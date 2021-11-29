// Copyright 2020 IBM Corp.
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
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/util/retry"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	emperrors "emperror.dev/errors"
	"github.com/go-logr/logr"
	olmv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/generated/clientset/versioned/scheme"
	marketplacev1alpha1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1alpha1"
	marketplacev1beta1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	utils "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

//var log = logf.Log.WithName("controller_olm_clusterserviceversion_watcher")

const (
	watchTag         string = "razee/watch-resource"
	olmCopiedFromTag string = "olm.copiedFrom"
	olmNamespace     string = "olm.operatorNamespace"
	ignoreTag        string = "marketplace.redhat.com/ignore"
	ignoreTagValue   string = "2"
	meterDefStatus   string = "marketplace.redhat.com/meterDefinitionStatus"
	meterDefError    string = "marketplace.redhat.com/meterDefinitionError"
)

// blank assignment to verify that ReconcileClusterServiceVersion implements reconcile.Reconciler
var _ reconcile.Reconciler = &ClusterServiceVersionReconciler{}

// ClusterServiceVersionReconciler reconciles a ClusterServiceVersion object
type ClusterServiceVersionReconciler struct {
	// This Client, initialized using mgr.Client() above, is a split Client
	// that reads objects from the cache and writes to the apiserver
	Client client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

// +kubebuilder:rbac:groups="operators.coreos.com",resources=clusterserviceversions;subscriptions,verbs=get;list;watch
// +kubebuilder:rbac:groups="operators.coreos.com",resources=clusterserviceversions,verbs=update;patch
// +kubebuilder:rbac:groups=marketplace.redhat.com,resources=meterdefinitions;meterdefinitions/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=marketplace.redhat.com,resources=meterdefinitions;meterdefinitions/status,verbs=get;list;watch;create;update;patch;delete

// Reconcile reads that state of the cluster for a ClusterServiceVersion object and makes changes based on the state read
// and what is in the ClusterServiceVersion.Spec
func (r *ClusterServiceVersionReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	reqLogger := r.Log.WithValues("Request.Name", request.Name, "Request.Namespace", request.Namespace)
	reqLogger.Info("Reconciling ClusterServiceVersion")
	// Fetch the ClusterServiceVersion instance
	CSV := &olmv1alpha1.ClusterServiceVersion{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, CSV)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Error(err, "clusterserviceversion does not exist")
			return reconcile.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get clusterserviceversion")
		return reconcile.Result{}, err
	}

	annotations := CSV.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	// check if CSV is being deleted
	// if yes -> finalizer logic
	// if no -> do nothing
	if !CSV.ObjectMeta.DeletionTimestamp.IsZero() {
		//Run finalization logic for the CSV.
		result, isReconcile, err := r.finalizeCSV(CSV)
		if isReconcile {
			return result, err
		}
	}

	result, isRequeue, err := r.reconcileMeterDefAnnotation(CSV, annotations)

	// check if err is instance of json.parsing error
	// if yes -> add failiure annotation

	if isRequeue {
		return result, err
	}
	sub := &olmv1alpha1.SubscriptionList{}

	if err := r.Client.List(context.TODO(), sub, client.InNamespace(request.NamespacedName.Namespace)); err != nil {
		return reconcile.Result{}, err
	}

	hasMarketplaceSub := false
	if len(sub.Items) > 0 {
		reqLogger.V(4).Info("found Subscription in namespaces", "count", len(sub.Items))
		// add razee watch label to CSV if subscription has rhm/operator label
		for _, s := range sub.Items {
			if value, ok := s.GetLabels()[operatorTag]; ok {
				if value == "true" {
					if len(s.Status.InstalledCSV) == 0 {
						reqLogger.Info("Requeue clusterserviceversion to wait for subscription getting installedCSV updated")
						return reconcile.Result{RequeueAfter: time.Second * 5}, nil
					}

					if s.Status.InstalledCSV == request.NamespacedName.Name {
						reqLogger.Info("found Subscription with installed CSV")
						hasMarketplaceSub = true

						if v, ok := CSV.GetLabels()[watchTag]; !ok || v != "lite" {
							err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
								err := r.Client.Get(context.TODO(),
									types.NamespacedName{
										Name:      CSV.GetName(),
										Namespace: CSV.GetNamespace(),
									},
									CSV)

								if err != nil {
									return err
								}

								labels := CSV.GetLabels()

								if labels == nil {
									labels = make(map[string]string)
								}

								labels[watchTag] = "lite"
								CSV.SetLabels(labels)

								return r.Client.Update(context.TODO(), CSV)
							})

							if err != nil {
								reqLogger.Error(err, "Failed to patch clusterserviceversion with razee/watch-resource: lite label")
								return reconcile.Result{}, err
							}
							reqLogger.Info("Patched clusterserviceversion with razee/watch-resource: lite label")
						} else {
							reqLogger.Info("No patch needed on clusterserviceversion resource")
						}
					}
				}
			}
		}
	} else {
		reqLogger.Info("Did not find Subscription in namespaces")
	}

	if !hasMarketplaceSub {
		reqLogger.Info("Does not have marketplace sub, ignoring CSV for future")

		if v, ok := CSV.GetAnnotations()[ignoreTag]; !ok || v != ignoreTagValue {
			retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				err := r.Client.Get(context.TODO(),
					types.NamespacedName{
						Name:      CSV.GetName(),
						Namespace: CSV.GetNamespace(),
					},
					CSV)

				if err != nil {
					return err
				}

				annotations := CSV.GetAnnotations()

				if annotations == nil {
					annotations = make(map[string]string)
				}

				annotations[ignoreTag] = ignoreTagValue
				CSV.SetAnnotations(annotations)

				return r.Client.Update(context.TODO(), CSV)
			})

			if retryErr != nil {
				reqLogger.Error(retryErr, "Failed to patch clusterserviceversion ignore tag")
				return reconcile.Result{Requeue: true}, retryErr
			}
			reqLogger.V(4).Info("Patched clusterserviceversion with ignore tag")
		} else {
			reqLogger.V(4).Info("No patch needed on clusterserviceversion resource for ignore tag")
		}
	}

	reqLogger.Info("reconciliation complete")
	return reconcile.Result{RequeueAfter: time.Minute * 1}, nil
}

func (r *ClusterServiceVersionReconciler) finalizeCSV(CSV *olmv1alpha1.ClusterServiceVersion) (reconcile.Result, bool, error) {
	reqLogger := r.Log.WithValues("Request.Name", CSV.GetName(), "Request.Namespace", CSV.GetNamespace())

	reqLogger.Info("deleting csv")
	if err := r.deleteExternalResources(CSV); err != nil {
		reqLogger.Error(err, "unable to delete csv")
		return reconcile.Result{}, false, err
	}

	// Stop reconciliation as the item is being deleted
	return reconcile.Result{}, true, nil
}

// deleteExternalResources searches for the MeterDefinition created by the CSV, if it's found delete it
func (r *ClusterServiceVersionReconciler) deleteExternalResources(CSV *olmv1alpha1.ClusterServiceVersion) error {
	reqLogger := r.Log.WithValues("Request.Name", CSV.GetName(), "Request.Namespace", CSV.GetNamespace())
	reqLogger.Info("deleting csv")

	annotations := CSV.GetAnnotations()
	if annotations == nil {
		reqLogger.Info("No annotations for this CSV")
		return nil
	}

	meterDefinitionString, ok := annotations[utils.CSV_METERDEFINITION_ANNOTATION]
	if !ok {
		reqLogger.Info("No value for ", "key: ", utils.CSV_METERDEFINITION_ANNOTATION)
		return nil
	}

	var errAlpha, errBeta error
	meterDefinitionBeta := &marketplacev1beta1.MeterDefinition{}
	meterDefinitionAlpha := &marketplacev1alpha1.MeterDefinition{}

	errBeta = meterDefinitionBeta.BuildMeterDefinitionFromString(meterDefinitionString, CSV.GetName(), CSV.GetNamespace(), utils.CSV_ANNOTATION_NAME, utils.CSV_ANNOTATION_NAMESPACE)

	if errBeta != nil {
		errAlpha = meterDefinitionAlpha.BuildMeterDefinitionFromString(meterDefinitionString, CSV.GetName(), CSV.GetNamespace(), utils.CSV_ANNOTATION_NAME, utils.CSV_ANNOTATION_NAMESPACE)
	}

	switch {
	case errBeta == nil:
		err := r.Client.Delete(context.TODO(), meterDefinitionBeta, client.PropagationPolicy(metav1.DeletePropagationForeground))
		if err != nil && errors.IsNotFound(err) {
			return err
		}
	case errAlpha == nil:
		err := r.Client.Delete(context.TODO(), meterDefinitionAlpha, client.PropagationPolicy(metav1.DeletePropagationForeground))
		if err != nil && errors.IsNotFound(err) {
			return err
		}
	default:
		err := emperrors.Combine(errBeta, errAlpha)
		reqLogger.Error(err, "Could not build a local copy of the MeterDefinition")
		return err
	}

	reqLogger.Info("found and deleted MeterDefinition")
	return nil
}

// reconcileMeterDefAnnotation checks the Annotations for the rhm CSV
// If the CSV is new, we tag it and create a MeterDefinition
// If the CSV is old, we check if the actual MeterDefinition matches the CSV (json) MeterDefinition
func (r *ClusterServiceVersionReconciler) reconcileMeterDefAnnotation(CSV *olmv1alpha1.ClusterServiceVersion, annotations map[string]string) (reconcile.Result, bool, error) {
	var err error
	reqLogger := r.Log.WithValues("CSV.Name", CSV.Name, "CSV.Namespace", CSV.Namespace)

	// checks if it is possible to build MeterDefinition from annotations of CSV
	reqLogger.Info("retrieving MeterDefinition string from csv")
	meterDefinitionString, ok := annotations[utils.CSV_METERDEFINITION_ANNOTATION]
	if !ok || len(meterDefinitionString) == 0 {
		reqLogger.Info("No value for ", "key: ", utils.CSV_METERDEFINITION_ANNOTATION)
		delete(annotations, meterDefError)
		delete(annotations, meterDefStatus)
		return reconcile.Result{}, false, nil
	}

	if ns, ok := CSV.GetAnnotations()[olmNamespace]; ok && ns != CSV.GetNamespace() {
		reqLogger.Info("MeterDef is global and this CSV is not the head")
		return reconcile.Result{}, false, nil
	}

	// builds a meterdefinition from our string (from the annotation)
	reqLogger.Info("retrieval successful", "str", meterDefinitionString)

	var errAlpha, errBeta error
	meterDefinitionBeta := &marketplacev1beta1.MeterDefinition{}
	meterDefinitionAlpha := &marketplacev1alpha1.MeterDefinition{}

	unstructured := &unstructured.Unstructured{}
	meterDefinition := &marketplacev1beta1.MeterDefinition{}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	_, objectKind, err := decode([]byte(meterDefinitionString), nil, nil)

	err = yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(meterDefinitionString)), 100).Decode(unstructured)
	if err == nil {
		switch {
		case objectKind.Version == "v1beta1":
			reqLogger.Info("mdef is a v1beta1", "value", meterDefinitionBeta)
			errBeta = meterDefinitionBeta.BuildMeterDefinitionFromString(
				meterDefinitionString,
				CSV.GetName(), CSV.GetNamespace(),
				utils.CSV_ANNOTATION_NAME, utils.CSV_ANNOTATION_NAMESPACE)

			meterDefinition = meterDefinitionBeta
		case objectKind.Version == "v1alpha1":
			reqLogger.Info("mdef is an v1alpha1")
			errAlpha = meterDefinitionAlpha.BuildMeterDefinitionFromString(
				meterDefinitionString,
				CSV.GetName(), CSV.GetNamespace(),
				utils.CSV_ANNOTATION_NAME, utils.CSV_ANNOTATION_NAMESPACE)

			if errAlpha == nil {
				err = meterDefinitionAlpha.ConvertTo(meterDefinition)

				if err != nil {
					reqLogger.Error(err, "Failed to convert to v1beta1")
				}
			}
		default:
			reqLogger.Info("mdef is neither")
			err = emperrors.Combine(err, errBeta, errAlpha)
			reqLogger.Error(err, "Failed to read the json annotation as a meterdefinition")
		}
	}

	if err != nil {
		reqLogger.Error(err, "Could not build a local copy of the MeterDefinition")
		reqLogger.Info("Adding failiure annotation in csv file ")
		annotations[meterDefStatus] = "error"
		annotations[meterDefError] = err.Error()
		CSV.SetAnnotations(annotations)
		if err := r.Client.Update(context.TODO(), CSV); err != nil {
			reqLogger.Error(err, "Failed to patch clusterserviceversion with MeterDefinition status")
			return reconcile.Result{}, true, err
		}
		reqLogger.Info("Patched clusterserviceversion with MeterDefinition status")
		return reconcile.Result{}, true, err
	}
	reqLogger.Info("marketplacev1beta1.MeterDefinitionList >>>> ")

	// Case 1: The CSV is old: compare vs. expected MeterDefinition
	list := &marketplacev1beta1.MeterDefinitionList{}
	err = r.Client.List(context.TODO(), list, client.InNamespace(meterDefinition.GetNamespace()))

	if err != nil {
		reqLogger.Error(err, "Could not retrieve the existing MeterDefinition")
		return reconcile.Result{}, true, err
	}
	reqLogger.Info("marketplacev1beta1.MeterDefinitionList End --- ")
	var actualMeterDefinition *marketplacev1beta1.MeterDefinition

	// Find the meterdef, we're use the InstalledBy field
	for _, meterDef := range list.Items {
		if meterDef.Spec.InstalledBy != nil &&
			meterDef.Spec.InstalledBy.Namespace == CSV.Namespace &&
			meterDef.Spec.InstalledBy.Name == CSV.Name {
			actualMeterDefinition = &meterDef
			reqLogger.Info("Found meterDef", "meterDef", meterDef)
			break
		}
	}

	// Check if the name has changed
	if actualMeterDefinition != nil && actualMeterDefinition.Name != meterDefinition.Name {
		reqLogger.Info("Discovered name change", "name", actualMeterDefinition.Name, "newName", meterDefinition.Name)
		err := r.Client.Delete(context.TODO(), actualMeterDefinition)

		if err != nil {
			return reconcile.Result{}, true, err
		}

		actualMeterDefinition = nil
	}

	// If not nil, we update
	if actualMeterDefinition != nil {
		if !reflect.DeepEqual(meterDefinition.Spec, actualMeterDefinition.Spec) &&
			!reflect.DeepEqual(meterDefinition.ObjectMeta, actualMeterDefinition.ObjectMeta) {
			reqLogger.Info("The actual meterdefinition is different from the expected meterdefinition")

			patch, err := json.Marshal(meterDefinition)
			if err != nil {
				return reconcile.Result{}, true, err
			}
			err = r.Client.Patch(context.TODO(), meterDefinition, client.RawPatch(types.MergePatchType, patch))
			if err != nil {
				return reconcile.Result{Requeue: true}, true, err
			}
			reqLogger.Info("Patch to update MeterDefinition successful. Requeuing")
			return reconcile.Result{Requeue: true}, true, nil
		}
		reqLogger.Info("meter definition matches")
		return reconcile.Result{}, false, nil
	}

	// Case 2: The CSV is new: we must track it & we must create the Meter Definition
	gvk, err := apiutil.GVKForObject(CSV, r.Scheme)
	if err != nil {
		return reconcile.Result{}, true, err
	}

	ref := metav1.OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               CSV.GetName(),
		UID:                CSV.GetUID(),
		BlockOwnerDeletion: pointer.BoolPtr(false),
		Controller:         pointer.BoolPtr(false),
	}

	meterDefinition.ObjectMeta.OwnerReferences = append(meterDefinition.ObjectMeta.OwnerReferences, ref)

	if err != nil {
		reqLogger.Error(err, "Failed to create.", "obj", meterDefinition)
		return reconcile.Result{}, true, err
	}

	meterDefinition.ObjectMeta.Namespace = CSV.Namespace

	err = r.Client.Create(context.TODO(), meterDefinition)
	if err != nil {
		reqLogger.Error(err, "Could not create MeterDefinition", "mdef", meterDefinition)
		reqLogger.Info("Adding failiure annotation in csv file ")
		annotations[meterDefStatus] = "error"
		annotations[meterDefError] = err.Error()
		CSV.SetAnnotations(annotations)
		if err := r.Client.Update(context.TODO(), CSV); err != nil {
			reqLogger.Error(err, "Failed to patch clusterserviceversion with MeterDefinition status")
			return reconcile.Result{}, true, err
		}
		reqLogger.Info("Patched clusterserviceversion with MeterDefinition status")
		return reconcile.Result{}, true, err
	}

	//Add success message annotation to csv
	delete(annotations, meterDefError)
	annotations[meterDefStatus] = "success"
	CSV.SetAnnotations(annotations)
	if err := r.Client.Update(context.TODO(), CSV); err != nil {
		reqLogger.Error(err, "Failed to patch clusterserviceversion with MeterDefinition status")
		return reconcile.Result{}, true, err
	}
	reqLogger.Info("Patched clusterserviceversion with MeterDefinition status")

	return reconcile.Result{}, true, nil
}

func csvFilter(metaNew metav1.Object) int {
	ann := metaNew.GetAnnotations()

	//annotation values
	ignoreVal, hasIgnoreTag := ann[ignoreTag]
	_, hasCopiedFrom := ann[olmCopiedFromTag]
	_, hasMeterDefinition := ann[utils.CSV_METERDEFINITION_ANNOTATION]

	sameNamespace := ann[olmNamespace] == metaNew.GetNamespace()

	switch {
	case hasMeterDefinition && !hasCopiedFrom && sameNamespace:
		return 1
	case !hasMeterDefinition && (!hasIgnoreTag || ignoreVal != ignoreTagValue):
		return 2
	default:
	}

	return 0
}

func checkForUpdateToMdef(evt event.UpdateEvent) bool {
	oldMeterDefVal, oldOk := evt.ObjectOld.GetAnnotations()[utils.CSV_METERDEFINITION_ANNOTATION]
	newMeterDefVal, newOk := evt.ObjectNew.GetAnnotations()[utils.CSV_METERDEFINITION_ANNOTATION]
	return oldOk && newOk && oldMeterDefVal != newMeterDefVal
}

var clusterServiceVersionPredictates predicate.Funcs = predicate.Funcs{
	UpdateFunc: func(evt event.UpdateEvent) bool {
		return csvFilter(evt.ObjectNew) > 0 && checkForUpdateToMdef(evt)
	},
	DeleteFunc: func(evt event.DeleteEvent) bool {
		return true
	},
	CreateFunc: func(evt event.CreateEvent) bool {
		return csvFilter(evt.Object) > 0
	},
	GenericFunc: func(evt event.GenericEvent) bool {
		return false
	},
}

func (r *ClusterServiceVersionReconciler) SetupWithManager(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&olmv1alpha1.ClusterServiceVersion{}, builder.WithPredicates(clusterServiceVersionPredictates)).
		Watches(
			&source.Kind{Type: &marketplacev1beta1.MeterDefinition{}}, &handler.EnqueueRequestForOwner{
				IsController: false,
				OwnerType:    &olmv1alpha1.ClusterServiceVersion{},
			}).
		Complete(r)
}
