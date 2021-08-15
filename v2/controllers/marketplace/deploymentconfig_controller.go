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
	// "bytes"

	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-logr/logr"
	olmv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"

	osappsv1 "github.com/openshift/api/apps/v1"

	marketplacev1beta1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/config"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/manifests"
	mktypes "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/types"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/patch"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/predicates"
	. "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/reconcileutils"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"

	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/catalog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// blank assignment to verify that DeploymentConfigReconciler implements reconcile.Reconciler
var _ reconcile.Reconciler = &DeploymentConfigReconciler{}

// var GlobalMeterdefStoreDB = &MeterdefStoreDB{}
// DeploymentConfigReconciler reconciles the DataService of a MeterBase object
type DeploymentConfigReconciler struct {
	// This Client, initialized using mgr.Client() above, is a split Client
	// that reads objects from the cache and writes to the apiserver
	Client        client.Client
	Scheme        *runtime.Scheme
	Log           logr.Logger
	CC            ClientCommandRunner
	kubeInterface kubernetes.Interface
	cfg           *config.OperatorConfig
	factory       *manifests.Factory
	patcher       patch.Patcher
	catalogClient *catalog.CatalogClient
}

func (r *DeploymentConfigReconciler) Inject(injector mktypes.Injectable) mktypes.SetupWithManager {
	injector.SetCustomFields(r)
	return r
}

func (r *DeploymentConfigReconciler) InjectOperatorConfig(cfg *config.OperatorConfig) error {
	r.cfg = cfg
	return nil
}

func (r *DeploymentConfigReconciler) InjectCommandRunner(ccp ClientCommandRunner) error {
	r.Log.Info("command runner")
	r.CC = ccp
	return nil
}

func (r *DeploymentConfigReconciler) InjectCatalogClient(catalogClient *catalog.CatalogClient) error {
	r.Log.Info("catalog client")
	r.catalogClient = catalogClient
	return nil
}

func (r *DeploymentConfigReconciler) InjectPatch(p patch.Patcher) error {
	r.patcher = p
	return nil
}

func (r *DeploymentConfigReconciler) InjectFactory(f *manifests.Factory) error {
	r.factory = f
	return nil
}

func (r *DeploymentConfigReconciler) InjectKubeInterface(k kubernetes.Interface) error {
	r.kubeInterface = k
	return nil
}

// adds a new Controller to mgr with r as the reconcile.Reconciler
func (r *DeploymentConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {

	nsPred := predicates.NamespacePredicate(r.cfg.DeployedNamespace)

	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(nsPred).
		For(&osappsv1.DeploymentConfig{}, builder.WithPredicates(
			predicate.Funcs{
				CreateFunc: func(e event.CreateEvent) bool {
					return e.Meta.GetName() == utils.DEPLOYMENT_CONFIG_NAME

				},
				UpdateFunc: func(e event.UpdateEvent) bool {
					return e.MetaNew.GetName() == utils.DEPLOYMENT_CONFIG_NAME
				},
				DeleteFunc: func(e event.DeleteEvent) bool {
					return e.Meta.GetName() == utils.DEPLOYMENT_CONFIG_NAME

				},
				GenericFunc: func(e event.GenericEvent) bool {
					return e.Meta.GetName() == utils.DEPLOYMENT_CONFIG_NAME
				},
			},
		)).
		Complete(r)

}

// +kubebuilder:rbac:groups=apps.openshift.io,resources=deploymentconfigs,verbs=get;list;watch
// +kubebuilder:rbac:urls=/list-for-version/*,verbs=get;
// +kubebuilder:rbac:urls=/get-system-meterdefs/*,verbs=get;
// +kubebuilder:rbac:urls=/meterdef-index-label/*,verbs=get;
// +kubebuilder:rbac:groups="authentication.k8s.io",resources=tokenreviews,verbs=create;get
// +kubebuilder:rbac:groups="authorization.k8s.io",resources=subjectaccessreviews,verbs=create;get

// Reconcile reads that state of the cluster for a MeterdefConfigmap object and makes changes based on the state read
// and what is in the MeterdefConfigmap.Spec
func (r *DeploymentConfigReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := r.Log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	err := r.catalogClient.SetTransport(r.Client, r.cfg, r.kubeInterface, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	dc := &osappsv1.DeploymentConfig{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: utils.DEPLOYMENT_CONFIG_NAME, Namespace: request.Namespace}, dc)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Error(err, "deploymentconfig does not exist")
			return reconcile.Result{}, nil
		}

		reqLogger.Error(err, "Failed to get deploymentconfig")
		return reconcile.Result{}, err
	}

	// catch the deploymentconfig as it's rolling out a new deployment and requeue until finished
	for _, c := range dc.Status.Conditions {
		if c.Type == osappsv1.DeploymentProgressing {
			if c.Reason != "NewReplicationControllerAvailable" || c.Status != corev1.ConditionTrue || dc.Status.LatestVersion == dc.Status.ObservedGeneration {
				reqLogger.Info("deploymentconfig has not finished rollout, requeueing")
				return reconcile.Result{RequeueAfter: time.Minute * 2}, err
			}
		}
	}

	reqLogger.Info("deploymentconfig is in ready state")
	latestVersion := dc.Status.LatestVersion

	result := r.pruneDeployPods(latestVersion, request, reqLogger)
	if !result.Is(Continue) {

		if result.Is(Error) {
			reqLogger.Error(result.GetError(), "Failed during pruning operation")
		}

		return result.Return()
	}

	result = r.sync(request, reqLogger)
	if !result.Is(Continue) {

		if result.Is(Error) {
			reqLogger.Error(result.GetError(), "Failed during sync operation")
		}

		return result.Return()
	}

	reqLogger.Info("finished reconciling")
	return reconcile.Result{RequeueAfter: time.Minute * 1}, nil
}

func (r *DeploymentConfigReconciler) sync(request reconcile.Request, reqLogger logr.Logger) *ExecResult {
	csvList := &olmv1alpha1.ClusterServiceVersionList{}

	err := r.Client.List(context.TODO(), csvList)
	if err != nil {
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	for _, csv := range csvList.Items {

		/*
			csv.Name = memcached-operator.v0.0.1
			splitName = memcached-operator
		*/
		splitName := strings.Split(csv.Name, ".")[0]
		csvVersion := csv.Spec.Version.Version.String()
		namespace := request.Namespace

		/*
			pings the file server for a map of labels we use to index meterdefintions that originated from the file server
			these labels also get added to a meterdefinition by the file server	before it returns
			split-name gets dynamically added on the call to get labels - could probably just use the split name in the controller however
			{
				"marketplace.redhat.com/installedOperatorNameTag": "<split-name>",
				"marketplace.redhat.com/isCommunityMeterdefintion": "true"
			}

		*/
		indexLabels, result := r.catalogClient.GetMeterdefIndexLabels(reqLogger, splitName)
		if !result.Is(Continue) {
			return result
		}

		catalogResponse, result := r.catalogClient.ListMeterdefintionsFromFileServer(splitName, csvVersion, namespace, reqLogger)
		if !result.Is(Continue) {
			return result
		}

		/*
			csv is on the cluster but doesn't have a csv dir or doesn't have mdefs in it's catalog listing
			delete all meterdefs for that csv
			once all meterdefs are deleted, skip to the next csv - no need to keep processing
		*/
		if catalogResponse.CatalogStatus.CatlogStatusType == catalog.CsvDoesNotHaveCatalogDirStatus || catalogResponse.CatalogStatus.CatlogStatusType == catalog.CsvHasNoMeterdefinitionsStatus {
			result = r.deleteAllMeterdefsForCsv(indexLabels, reqLogger)
			if !result.Is(Continue) {
				return result
			}

			//skip to the next csv in csvList
			continue
		}

		/*
			if the csv has a listing in the meterdefinition catalog and the directory contains meterdefinitions, run a sync on those meterdefinitions
			i.e. catalogResponse.CatalogStatus.CatlogStatusType == catalog.CsvHasMeterdefinitions
		*/

		// fetch system meter definitions and append
		systemMeterdefsResponse, result := r.catalogClient.GetSystemMeterdefs(csv, reqLogger)
		if !result.Is(Continue) {
			return result
		}

		allMeterDefinitions := catalogResponse.MdefSlice
		allMeterDefinitions = append(allMeterDefinitions, systemMeterdefsResponse.MdefSlice...)

		installedMeterdefs, result := listAllCommunityMeterdefsOnCluster(r.Client, indexLabels)
		if !result.Is(Continue) {
			return result
		}

		/*
			delete a meterdef if there is a meterdef that originated from the catalog installed on the cluster, but that meterdef isn't in the latest file server image
		*/
		result = r.deleteOnDiff(installedMeterdefs.Items, allMeterDefinitions, reqLogger)
		if !result.Is(Continue) {
			return result
		}

		for _, meterdef := range allMeterDefinitions {

			installedMdef := &marketplacev1beta1.MeterDefinition{}
			reqLogger.Info("finding meterdefintion", meterdef.Name, meterdef.Namespace)
			err = r.Client.Get(context.TODO(), types.NamespacedName{Name: meterdef.Name, Namespace: meterdef.Namespace}, installedMdef)
			if err != nil && errors.IsNotFound(err) {
				/*
					create a meterdef for a csv if the csv has a meterdefinition listed in the catalog and that meterdef is not on the cluster
				*/
				result = r.createMeterdef(meterdef, &csv, reqLogger)
				if !result.Is(Continue) {
					return result
				}
			} else if err != nil {
				return &ExecResult{
					ReconcileResult: reconcile.Result{},
					Err:             err,
				}
			} else {
				/*
					update a meterdef for a csv if a meterdef from the catalog is also on the cluster
					&& the meterdef from the catalog contains an update to .Spec or .Annotations
					//TODO: what fields should we check a diff for ?
				*/
				result := r.updateMeterdef(installedMdef, meterdef, reqLogger)
				if !result.Is(Continue) {
					return result
				}
			}
		}
	}

	return &ExecResult{
		Status: ActionResultStatus(Continue),
	}
}

func (r *DeploymentConfigReconciler) updateMeterdef(installedMdef *marketplacev1beta1.MeterDefinition, catalogMdef marketplacev1beta1.MeterDefinition, reqLogger logr.Logger) *ExecResult {
	updatedMeterdefinition := installedMdef.DeepCopy()
	updatedMeterdefinition.Spec = catalogMdef.Spec
	updatedMeterdefinition.ObjectMeta.Annotations = catalogMdef.ObjectMeta.Annotations

	if !reflect.DeepEqual(updatedMeterdefinition, installedMdef) {
		reqLogger.Info("meterdefintion is out of sync with latest meterdef catalog", "Name", installedMdef.Name)
		err := r.Client.Update(context.TODO(), updatedMeterdefinition)
		if err != nil {
			reqLogger.Error(err, "Failed updating meter definition", "Name", updatedMeterdefinition.Name, "Namespace", updatedMeterdefinition.Namespace)
			return &ExecResult{
				ReconcileResult: reconcile.Result{},
				Err:             err,
			}
		}
		reqLogger.Info("Updated meterdefintion", "Name", updatedMeterdefinition.Name, "Namespace", updatedMeterdefinition.Namespace)

		return &ExecResult{
			ReconcileResult: reconcile.Result{Requeue: true},
			Err:             nil,
		}
	}

	return &ExecResult{
		Status: ActionResultStatus(Continue),
	}
}

func (r *DeploymentConfigReconciler) createMeterdef(meterDefinition marketplacev1beta1.MeterDefinition, csv *olmv1alpha1.ClusterServiceVersion, reqLogger logr.InfoLogger) *ExecResult {

	gvk, err := apiutil.GVKForObject(csv, r.Scheme)
	if err != nil {
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	// create owner ref object
	ref := metav1.OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               csv.GetName(),
		UID:                csv.GetUID(),
		BlockOwnerDeletion: pointer.BoolPtr(false),
		Controller:         pointer.BoolPtr(false),
	}

	meterDefinition.ObjectMeta.SetOwnerReferences([]metav1.OwnerReference{ref})
	meterDefinition.ObjectMeta.Namespace = csv.Namespace

	meterDefName := meterDefinition.Name
	err = r.Client.Create(context.TODO(), &meterDefinition)
	if err != nil {
		reqLogger.Error(err, "Could not create meterdefinition", "mdef", meterDefName, "CSV", csv.Name)
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}
	reqLogger.Info("Created meterdefinition", "mdef", meterDefName, "CSV", csv.Name)

	return &ExecResult{
		ReconcileResult: reconcile.Result{Requeue: true},
		Err:             nil,
	}
}

func (r *DeploymentConfigReconciler) deleteOnDiff(installedMeterdefs []marketplacev1beta1.MeterDefinition, allMeterDefinitions []marketplacev1beta1.MeterDefinition, reqLogger logr.Logger) *ExecResult {

	for _, installedMeterdef := range installedMeterdefs {
		found := false
		for _, meterdefFromCatalog := range allMeterDefinitions {

			if installedMeterdef.Name == meterdefFromCatalog.Name {
				reqLogger.Info("MATCH", installedMeterdef.Name, meterdefFromCatalog.Name)
				found = true
				break
			}
		}

		if !found {
			reqLogger.Info("meterdef has been selected for deletion", "meterdef", installedMeterdef.Name)
			result := r.deleteMeterDef(installedMeterdef.Name, installedMeterdef.Namespace, reqLogger)
			if !result.Is(Continue) {
				return result
			}
		}
	}

	return &ExecResult{
		Status: ActionResultStatus(Continue),
	}
}

func listAllCommunityMeterdefsOnCluster(runtimeClient client.Client, indexLabels map[string]string) (*marketplacev1beta1.MeterDefinitionList, *ExecResult) {
	installedMeterdefList := &marketplacev1beta1.MeterDefinitionList{}

	// look for meterdefs that originated from the meterdefinition catalog
	listOpts := []client.ListOption{
		client.MatchingLabels(indexLabels),
	}

	err := runtimeClient.List(context.TODO(), installedMeterdefList, listOpts...)
	if err != nil {
		return nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	return installedMeterdefList, &ExecResult{
		Status: ActionResultStatus(Continue),
	}
}

func (r *DeploymentConfigReconciler) deleteAllMeterdefsForCsv(indexLabels map[string]string, reqLogger logr.Logger) *ExecResult {
	installedMeterdefList := &marketplacev1beta1.MeterDefinitionList{}

	// look for meterdefs that are from the meterdefinition catalog
	listOpts := []client.ListOption{
		client.MatchingLabels(indexLabels),
	}

	err := r.Client.List(context.TODO(), installedMeterdefList, listOpts...)
	if err != nil {
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	for _, mdef := range installedMeterdefList.Items {
		result := r.deleteMeterDef(mdef.Name, mdef.Namespace, reqLogger)
		if !result.Is(Continue) {
			return result
		}
	}

	return &ExecResult{
		Status: ActionResultStatus(Continue),
	}
}

func (r *DeploymentConfigReconciler) deleteMeterDef(mdefName string, namespace string, reqLogger logr.Logger) *ExecResult {

	installedMeterDefn := &marketplacev1beta1.MeterDefinition{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: mdefName, Namespace: namespace}, installedMeterDefn)
	if err != nil && !errors.IsNotFound((err)) {
		reqLogger.Error(err, "could not get meter definition", "Name", mdefName)
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	// remove owner ref from meter definition before deleting
	installedMeterDefn.ObjectMeta.OwnerReferences = []metav1.OwnerReference{}
	err = r.Client.Update(context.TODO(), installedMeterDefn)
	if err != nil {
		reqLogger.Error(err, "Failed updating owner reference on meter definition", "Name", mdefName, "Namespace", namespace)
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}
	reqLogger.Info("Removed owner reference from meterdefintion", "Name", mdefName, "Namespace", namespace)

	//TODO: requeue here ?

	reqLogger.Info("Deleteing MeterDefinition")
	err = r.Client.Delete(context.TODO(), installedMeterDefn)
	if err != nil && !errors.IsNotFound(err) {
		reqLogger.Error(err, "could not delete MeterDefinition", "Name", mdefName)
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}
	reqLogger.Info("Deleted meterdefintion", "Name", mdefName, "Namespace", namespace)

	//TODO: requeue here ?

	return &ExecResult{
		Status: ActionResultStatus(Continue),
	}
}

/*
	TODO: test whether revisionHistoryLimit will handle this
*/
func (r *DeploymentConfigReconciler) pruneDeployPods(latestVersion int64, request reconcile.Request, reqLogger logr.Logger) *ExecResult {
	reqLogger.Info("pruning old deploy pods")

	latestPodName := fmt.Sprintf("rhm-meterdefinition-file-server-%d", latestVersion)
	reqLogger.Info("Prune", "latest version", latestVersion)
	reqLogger.Info("Prune", "latest pod name", latestPodName)

	dcPodList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(request.Namespace),
		client.HasLabels{"openshift.io/deployer-pod-for.name"},
	}

	err := r.Client.List(context.TODO(), dcPodList, listOpts...)
	if err != nil {
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	for _, pod := range dcPodList.Items {
		reqLogger.Info("Prune", "deploy pod", pod.Name)
		podLabelValue := pod.GetLabels()["openshift.io/deployer-pod-for.name"]
		if podLabelValue != latestPodName {

			err := r.Client.Delete(context.TODO(), &pod)
			if err != nil {
				return &ExecResult{
					ReconcileResult: reconcile.Result{},
					Err:             err,
				}
			}

			reqLogger.Info("Successfully pruned deploy pod", "pod name", pod.Name)
		}
	}

	return &ExecResult{
		Status: ActionResultStatus(Continue),
	}
}
