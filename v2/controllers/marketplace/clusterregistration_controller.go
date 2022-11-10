// Copyright 2021 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package marketplace

import (
	"context"
	"errors"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/gotidy/ptr"
	openshiftconfigv1 "github.com/openshift/api/config/v1"
	marketplacev1alpha1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1alpha1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/config"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/marketplace"
	mktypes "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/types"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/predicates"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// blank assignment to verify that ReconcileClusterRegistration implements reconcile.Reconciler
var _ reconcile.Reconciler = &ClusterRegistrationReconciler{}

// ClusterRegistrationReconciler reconciles a Registration object
type ClusterRegistrationReconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	Client client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger

	cfg *config.OperatorConfig
}

type SecretInfo struct {
	TypeOf     string
	Secret     *v1.Secret
	StatusKey  string
	MessageKey string
	SecretKey  string
	MissingMsg string
}

// +kubebuilder:rbac:groups="",namespace=system,resources=secrets,verbs=get;list;watch;create
// +kubebuilder:rbac:groups="",namespace=system,resources=secrets,resourceNames=redhat-marketplace-pull-secret;ibm-entitlement-key;rhm-operator-secret,verbs=update;patch
// +kubebuilder:rbac:groups="apps",namespace=system,resources=deployments,verbs=get;list;watch
// +kubebuilder:rbac:groups="apps",namespace=system,resources=deployments/finalizers,verbs=get;list;watch;update;patch,resourceNames=redhat-marketplace-controller-manager
// +kubebuilder:rbac:groups=marketplace.redhat.com,namespace=system,resources=marketplaceconfigs;marketplaceconfigs/finalizers;marketplaceconfigs/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="config.openshift.io",resources=clusterversions,verbs=get;list;watch
// +kubebuilder:rbac:groups="apps",namespace=system,resources=deployments;deployments/finalizers,verbs=get;list;watch
// +kubebuilder:rbac:groups="apps",namespace=system,resources=deployments/finalizers,verbs=get;list;watch;update;patch,resourceNames=redhat-marketplace-controller-manager

// Reconcile reads that state of the cluster for a ClusterRegistration object and makes changes based on the state read
// and what is in the ClusterRegistration.Spec
func (r *ClusterRegistrationReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	reqLogger := r.Log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ClusterRegistration")
	secretFetcher := utils.ProvideSecretFetcherBuilder(r.Client, context.TODO(), request.Namespace)
	si, err := secretFetcher.ReturnSecret()
	if err == utils.NoSecretsFound {
		// marketplaceconfig controller will provide status in this case
		reqLogger.Info("no redhat-marketplace-pull-secret or ibm-entitlement-key secret found, secret is required in a connected environment")
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("found secret", "secret", si.Name)

	annotations := si.Secret.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	jwtToken, err := secretFetcher.ParseAndValidate(si)
	if err != nil {
		reqLogger.Error(err, "error validating secret")
		if errors.Is(err, utils.TokenFieldMissingOrEmpty) {
			return r.updateSecretWithMessage(si, annotations, reqLogger)
		}

		return reconcile.Result{}, err
	}

	if jwtToken == "" {
		err := errors.New("jwt token is empty")
		reqLogger.Error(err, "couldn't find secret field")
		return reconcile.Result{}, err
	}

	tokenClaims, err := marketplace.GetJWTTokenClaim(jwtToken)
	if err != nil {
		reqLogger.Error(err, "Token is missing account id")
		annotations[si.StatusKey] = "error"
		annotations[si.MessageKey] = "Account id is not available in provided token, please generate token from RH Marketplace again"
		si.Secret.SetAnnotations(annotations)
		if err := r.Client.Update(context.TODO(), si.Secret); err != nil {
			reqLogger.Error(err, "Failed to patch secret with Endpoint status")
		}
		reqLogger.Info("Secret updated with account id missing error")
		return reconcile.Result{}, err
	}

	reqLogger.Info("Marketplace Token Claims set")

	marketplaceConfig := &marketplacev1alpha1.MarketplaceConfig{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: request.Namespace,
		Name:      utils.MARKETPLACECONFIG_NAME,
	}, marketplaceConfig)

	// Create a default MarketplaceConfig if a pull-secret is created, connected environment quickstart
	if k8serrors.IsNotFound(err) {
		marketplaceConfig.Name = utils.MARKETPLACECONFIG_NAME
		marketplaceConfig.Namespace = r.cfg.DeployedNamespace

		// IS_DISCONNECTED flag takes precedence, default IsDisconnected to false
		if r.cfg.IsDisconnected {
			marketplaceConfig.Spec.IsDisconnected = ptr.Bool(true)
		} else if marketplaceConfig.Spec.IsDisconnected == nil {
			marketplaceConfig.Spec.IsDisconnected = ptr.Bool(false)
		}

		// Set required ClusterID
		reqLogger.Info("finding clusterversion resource")
		clusterVersion := &openshiftconfigv1.ClusterVersion{}
		err = r.Client.Get(context.Background(), client.ObjectKey{
			Name: "version",
		}, clusterVersion)

		if err != nil {
			if !k8serrors.IsNotFound(err) && !meta.IsNoMatchError(err) {
				reqLogger.Error(err, "Failed to retrieve clusterversion resource")
				return reconcile.Result{}, err
			}
			clusterVersion = nil
		}

		var clusterID string
		if clusterVersion != nil {
			clusterID = string(clusterVersion.Spec.ClusterID)
			reqLogger.Info("Clusterversion object found with clusterID", "clusterID", clusterID)
		} else {
			clusterID = uuid.New().String()
			reqLogger.Info("Clusterversion object not found, generating clusterID", "clusterID", clusterID)
		}

		marketplaceConfig.Spec.ClusterUUID = string(clusterID)

		// Set required AccountID, reconcile later

		marketplaceConfig.Spec.RhmAccountID = ""

		// Set the controller deployment as the controller-ref, since it owns the finalizer
		dep := &appsv1.Deployment{}
		err := r.Client.Get(context.TODO(), types.NamespacedName{
			Name:      utils.RHM_CONTROLLER_DEPLOYMENT_NAME,
			Namespace: r.cfg.DeployedNamespace,
		}, dep)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err = controllerutil.SetControllerReference(dep, marketplaceConfig, r.Scheme); err != nil {
			return reconcile.Result{}, err
		}

		if err = r.Client.Create(context.TODO(), marketplaceConfig); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Set MarketplaceConfig as the Owner of the pull secret
	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		secretFetcher.GetPullSecret()
		return secretFetcher.AddOwnerRefToAll(marketplaceConfig, r.Scheme)
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	if !ptr.ToBool(marketplaceConfig.Spec.IsDisconnected) {
		//only check registration status, compare pull secret from COS if we are not in a disconnected environment
		mclient, err := marketplace.NewMarketplaceClientBuilder(r.cfg).
			NewMarketplaceClient(jwtToken, tokenClaims)

		if err != nil {
			reqLogger.Error(err, "failed to build marketplaceclient")
			return reconcile.Result{}, nil
		}

		reqLogger.Info("token found", "from secret", si.Name)

		//Calling POST endpoint to pull the secret definition
		newOptSecretObj, err := mclient.GetMarketplaceSecret()
		if err != nil {
			reqLogger.Error(err, "failed to get operator secret from marketplace")
			annotations[si.StatusKey] = "error"
			annotations[si.MessageKey] = err.Error()
			si.Secret.SetAnnotations(annotations)
			if err := r.Client.Update(context.TODO(), si.Secret); err != nil {
				reqLogger.Error(err, "Failed to patch secret with Endpoint status")
			}
			return reconcile.Result{}, err
		}

		// Secret comes back with namespace "redhat-marketplace-operator", set to request namespace
		newOptSecretObj.Namespace = request.Namespace

		//Fetch the Secret with name rhm-operator-secret
		secretKeyname := types.NamespacedName{
			Name:      newOptSecretObj.Name,
			Namespace: newOptSecretObj.Namespace,
		}

		reqLogger.Info("retrieving operator secret", "name", secretKeyname.Name, "namespace", secretKeyname.Namespace)

		err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			optSecret := &v1.Secret{}
			err = r.Client.Get(context.TODO(), secretKeyname, optSecret)
			if err != nil {
				if !k8serrors.IsNotFound(err) {
					reqLogger.Error(err, "error getting operator secret")
					return err
				}

				controllerutil.SetOwnerReference(si.Secret, newOptSecretObj, r.Scheme)
				reqLogger.Info("secret not found, creating")
				err = r.Client.Create(context.TODO(), newOptSecretObj)
				if err != nil {
					reqLogger.Error(err, "Failed to create operator object")
					return err
				}
			} else {
				reqLogger.Info("Comparing old and new rhm-operator-secret")

				if !reflect.DeepEqual(newOptSecretObj.Data, optSecret.Data) {
					reqLogger.Info("rhm-operator-secret are different copy")
					optSecret.Data = newOptSecretObj.Data
					controllerutil.SetOwnerReference(si.Secret, newOptSecretObj, r.Scheme)
					reqLogger.Info("updating rhm-operator-secret")
					return r.Client.Update(context.TODO(), optSecret)
				}
			}

			return nil
		})
		if err != nil {
			return reconcile.Result{}, err
		}

		// Update the AccountID & annotations field in MarketplaceConfig
		err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			if err = r.Client.Get(context.TODO(), types.NamespacedName{
				Namespace: request.Namespace,
				Name:      utils.MARKETPLACECONFIG_NAME,
			}, marketplaceConfig); err != nil {
				return err
			}

			updated := false
			if si.Name == utils.IBMEntitlementKeySecretName {
				if marketplaceConfig.Spec.RhmAccountID != string(newOptSecretObj.Data[utils.BUCKET_NAME_FIELD]) {
					marketplaceConfig.Spec.RhmAccountID = string(newOptSecretObj.Data[utils.BUCKET_NAME_FIELD])
					updated = true
				}
			}
			if si.Name == utils.RHMPullSecretName {
				if marketplaceConfig.Spec.RhmAccountID != tokenClaims.AccountID {
					marketplaceConfig.Spec.RhmAccountID = tokenClaims.AccountID
					updated = true
				}
			}

			// If the annotation is not set, or (unlikely) does not match tokenClaims.Env
			env, ok := marketplaceConfig.GetAnnotations()["marketplace.redhat.com/environment"]
			if ok {
				if env != tokenClaims.Env {
					ok = false
				}
			}
			if !ok {
				marketplaceConfig.SetAnnotations(map[string]string{"marketplace.redhat.com/environment": tokenClaims.Env})
				updated = true
			}

			if updated {
				reqLogger.Info("updating marketplaceconfig")
				return r.Client.Update(context.TODO(), marketplaceConfig)
			}

			return nil
		})
		if err != nil {
			return reconcile.Result{}, err
		}

		//Setting MarketplaceClientAccount
		marketplaceClientAccount := &marketplace.MarketplaceClientAccount{
			AccountId:   marketplaceConfig.Spec.RhmAccountID,
			ClusterUuid: marketplaceConfig.Spec.ClusterUUID,
		}

		// Marketplace config object found
		reqLogger.Info("Pulling MarketPlace config object status")
		registrationStatusOutput, _ := mclient.RegistrationStatus(marketplaceClientAccount)

		if registrationStatusOutput.RegistrationStatus == marketplace.RegistrationStatusInstalled {
			reqLogger.Info("MarketPlace config object is already registered for account")

			err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
				si, err := secretFetcher.ReturnSecret()
				if err != nil {
					return err
				}
				//Update secret with status
				if annotations[si.SecretKey] != "success" {
					reqLogger.Info("Updating secret with success status")
					annotations[si.StatusKey] = "success"
					annotations[si.MessageKey] = "rhm-operator-secret generated successfully"
					si.Secret.SetAnnotations(annotations)
					return r.Client.Update(context.TODO(), si.Secret)
				}

				return nil
			})
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	reqLogger.Info("ClusterRegistrationController reconcile finished")
	return reconcile.Result{}, nil
}

func (r *ClusterRegistrationReconciler) Inject(injector mktypes.Injectable) mktypes.SetupWithManager {
	injector.SetCustomFields(r)
	return r
}

func (m *ClusterRegistrationReconciler) InjectOperatorConfig(cfg *config.OperatorConfig) error {
	m.cfg = cfg
	return nil
}

func (r *ClusterRegistrationReconciler) updateSecretWithMessage(si *utils.SecretInfo, annotations map[string]string, reqLogger logr.Logger) (reconcile.Result, error) {
	reqLogger.Info("Missing token field in secret")
	annotations[si.StatusKey] = "error"
	annotations[si.MessageKey] = si.MissingMsg
	si.Secret.SetAnnotations(annotations)
	err := r.Client.Update(context.TODO(), si.Secret)
	if err != nil {
		reqLogger.Error(err, "Failed to patch secret with Endpoint status")
	}
	reqLogger.Info("Secret updated with status on failiure")
	return reconcile.Result{}, err
}

func (r *ClusterRegistrationReconciler) setControllerReference(controlled metav1.Object) error {
	// Set the controller deployment as the controller-ref, since it owns the finalizer
	dep := &appsv1.Deployment{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      utils.RHM_CONTROLLER_DEPLOYMENT_NAME,
		Namespace: r.cfg.DeployedNamespace,
	}, dep)
	if err != nil {
		return err
	}
	if err = controllerutil.SetControllerReference(dep, controlled, r.Scheme); err != nil {
		return err
	}
	return nil
}

func (r *ClusterRegistrationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	namespacePredicate := predicates.NamespacePredicate(r.cfg.DeployedNamespace)
	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(namespacePredicate).
		For(&v1.Secret{}, builder.WithPredicates(
			predicate.Funcs{
				CreateFunc: func(e event.CreateEvent) bool {
					secret, ok := e.Object.(*v1.Secret)
					if !ok {
						return false
					}
					secretName := secret.ObjectMeta.Name
					if _, ok := secret.Data[utils.RHMPullSecretKey]; ok && secretName == utils.RHMPullSecretName {
						return true
					}

					if _, ok := secret.Data[utils.IBMEntitlementDataKey]; ok && secretName == utils.IBMEntitlementKeySecretName {
						return true
					}

					return false
				},
				UpdateFunc: func(e event.UpdateEvent) bool {
					secret, ok := e.ObjectNew.(*v1.Secret)
					secretName := secret.ObjectMeta.Name
					if !ok {
						return false
					}

					secretOld, ok := e.ObjectOld.(*v1.Secret)
					if !ok {
						return false
					}

					if secretName == utils.RHMPullSecretName {
						if _, ok := secret.Data[utils.RHMPullSecretKey]; ok && e.ObjectOld != e.ObjectNew {
							return true
						}
					}

					if secretName == utils.IBMEntitlementKeySecretName {
						if _, ok := secret.Data[utils.IBMEntitlementDataKey]; ok && e.ObjectOld != e.ObjectNew {
							return true
						}
					}

					if secretName == utils.RHM_OPERATOR_SECRET_NAME {
						if !reflect.DeepEqual(secret.Data, secretOld.Data) {
							return true
						}
					}

					return false
				},
				DeleteFunc: func(e event.DeleteEvent) bool {
					return e.Object.GetName() == utils.RHM_OPERATOR_SECRET_NAME
				},
				GenericFunc: func(e event.GenericEvent) bool {
					secret, ok := e.Object.(*v1.Secret)
					if !ok {
						return false
					}
					secretName := secret.ObjectMeta.Name
					if _, ok := secret.Data[utils.RHMPullSecretKey]; ok && secretName == utils.RHMPullSecretName {
						return true
					}

					if _, ok := secret.Data[utils.IBMEntitlementDataKey]; ok && secretName == utils.IBMEntitlementKeySecretName {
						return true
					}

					if secretName == utils.RHM_OPERATOR_SECRET_NAME {
						return true
					}

					return false
				},
			},
		)).
		// This controller is responsible for keeping the AccountID & Env annotation updated on MarketplaceConfig (connected)
		// Since the values are dervived from the pull secret
		// It also must reconcile due a change from IsDisconnected to connected mode
		Watches(
			&source.Kind{Type: &marketplacev1alpha1.MarketplaceConfig{}},
			handler.EnqueueRequestsFromMapFunc(func(a client.Object) []reconcile.Request {
				return []reconcile.Request{
					{NamespacedName: types.NamespacedName{
						Name:      utils.RHMPullSecretName,
						Namespace: a.GetNamespace(),
					}},
				}
			}),
			builder.WithPredicates(predicate.Funcs{
				// Queue the reconciler in a connected environment
				// Handle the case where start in disconnected mode, and switch to connected
				CreateFunc: func(e event.CreateEvent) bool {
					marketplaceConfig, ok := e.Object.(*marketplacev1alpha1.MarketplaceConfig)
					if ok {
						return !ptr.ToBool(marketplaceConfig.Spec.IsDisconnected)
					}
					return false
				},
				UpdateFunc: func(e event.UpdateEvent) bool {

					marketplaceConfigNew, newOk := e.ObjectNew.(*marketplacev1alpha1.MarketplaceConfig)
					marketplaceConfigOld, oldOk := e.ObjectOld.(*marketplacev1alpha1.MarketplaceConfig)

					if newOk && oldOk {
						if ptr.ToBool(marketplaceConfigNew.Spec.IsDisconnected) == false && ptr.ToBool(marketplaceConfigOld.Spec.IsDisconnected) == true {
							return true
						}

						if ptr.ToBool(marketplaceConfigNew.Spec.IsDisconnected) {
							if marketplaceConfigNew.Spec.RhmAccountID != marketplaceConfigOld.Spec.RhmAccountID {
								return true
							}
							if !reflect.DeepEqual(marketplaceConfigNew.GetAnnotations(), marketplaceConfigOld.GetAnnotations()) {
								return true
							}
						}
					}

					return false
				},
				DeleteFunc:  func(event.DeleteEvent) bool { return false },
				GenericFunc: func(e event.GenericEvent) bool { return false },
			},
			)).
		Complete(r)
}
