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
	"crypto/tls"
	"crypto/x509"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	emperror "emperror.dev/errors"
	"github.com/go-logr/logr"
	olmv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	common "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/common"
	marketplacev1beta1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/config"
	prom "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/prometheus"
	mktypes "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/types"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils"
	. "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/reconcileutils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// CSVInfo to hold information about CSVs
type CSVInfo struct {
	CsvName      string
	CsvVersion   string
	CsvNamespace string
}

//go:embed meterdef-templates
var meterdefTemplates embed.FS

// blank assignment to verify that ReconcileClusterServiceVersion implements reconcile.Reconciler
var _ reconcile.Reconciler = &MeterdefinitionInstallReconciler{}

// MeterdefinitionInstallReconciler reconciles a ClusterServiceVersion object
type MeterdefinitionInstallReconciler struct {
	// This Client, initialized using mgr.Client() above, is a split Client
	// that reads objects from the cache and writes to the apiserver
	Client        client.Client
	Scheme        *runtime.Scheme
	Log           logr.Logger
	cfg           *config.OperatorConfig
	kubeInterface kubernetes.Interface
}

// +kubebuilder:rbac:groups="operators.coreos.com",resources=clusterserviceversions;subscriptions,verbs=get;list;watch
// +kubebuilder:rbac:groups="operators.coreos.com",resources=clusterserviceversions,verbs=update;patch
// +kubebuilder:rbac:groups=marketplace.redhat.com,resources=meterdefinitions;meterdefinitions/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=marketplace.redhat.com,resources=meterdefinitions;meterdefinitions/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="operators.coreos.com",resources=clusterserviceversions;subscriptions,verbs=get;list;watch
// +kubebuilder:rbac:urls=/list-for-version/*,verbs=get;
// +kubebuilder:rbac:groups="authentication.k8s.io",resources=tokenreviews,verbs=create;get
// +kubebuilder:rbac:groups="authorization.k8s.io",resources=subjectaccessreviews,verbs=create;get

// Reconcile reads that state of the cluster for a ClusterServiceVersion object and creates corresponding meter definitions if found
func (r *MeterdefinitionInstallReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := r.Log.WithValues("Request.Name", request.Name, "Request.Namespace", request.Namespace)
	reqLogger.Info("Reconciling ClusterServiceVersion")

	// Fetch the ClusterServiceVersion instance
	CSV := &olmv1alpha1.ClusterServiceVersion{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, CSV)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, check the meterdef store if there is an existing InstallMapping,delete, and return empty result
			// If no install mapping found return empty result
			_csvName := strings.Split(request.Name, ".")[0]
			delimiter := "."
			rightOfDelimiter := strings.Join(strings.Split(request.Name, delimiter)[1:], delimiter)
			_version := strings.Split(rightOfDelimiter, "v")[1]

			reqLogger.Info("clusterserviceversion does not exist, checking install map for package", "name", _csvName, "version", _version)
			result := deleteInstallMapping(_csvName, _version, r.cfg.DeployedNamespace, request, r.Client, reqLogger)
			if !result.Is(Continue) {
				if result.Is(Error) {
					reqLogger.Error(result.GetError(), "Failed deleting install mapping from meter definition store")
				}
				return result.Return()
			}

		}

		reqLogger.Error(err, "Failed to get clusterserviceversion")
		return reconcile.Result{}, err
	}

	csvName := strings.Split(CSV.Name, ".")[0]
	reqLogger.Info("csv name", "name", csvName)
	csvVersion := CSV.Spec.Version.Version.String()
	reqLogger.Info("csv version", "version", csvVersion)

	// New CSV install
	sub := &olmv1alpha1.SubscriptionList{}
	if err := r.Client.List(context.TODO(), sub, client.InNamespace(request.NamespacedName.Namespace)); err != nil {
		return reconcile.Result{}, err
	}

	if len(sub.Items) > 0 {
		reqLogger.Info("found Subscription in namespaces", "count", len(sub.Items))

		// apply meter definition if subscription has rhm/operator label
		for _, subscription := range sub.Items {
			if value, ok := subscription.GetLabels()[operatorTag]; ok {

				if value == "true" {
					if len(subscription.Status.InstalledCSV) == 0 {
						reqLogger.Info("Requeue clusterserviceversion to wait for subscription getting")
						return reconcile.Result{RequeueAfter: time.Second * 5}, nil
					}

					installedOperatorName, ok := subscription.GetAnnotations()[installedOperatorNameTag]
					if ok {
						if subscription.Status.InstalledCSV != request.NamespacedName.Name && csvName == installedOperatorName {
							reqLogger.Info("subscription installed csv", "installed csv", subscription.Status.InstalledCSV)
							return reconcile.Result{RequeueAfter: time.Second * 5}, nil
						}
					}

					if subscription.Status.InstalledCSV == request.NamespacedName.Name {
						reqLogger.Info("found Subscription with installed CSV")

						_, selectedMeterDefinitions, result := ListMeterdefintionsFromFileServer(csvName, csvVersion, CSV.Namespace, r.Client, r.kubeInterface, r.cfg.DeployedNamespace, reqLogger)
						if !result.Is(Continue) {

							if result.Is(Error) {
								reqLogger.Error(result.GetError(), "Failed retrieving meterdefinitions from file server", "CSV", csvName)
							}

							return result.Return()
						}

						globalMeterdefinitions, result := r.fetchGlobalMeterdefs(CSV, reqLogger)
						if !result.Is(Continue) {

							if result.Is(Error) {
								reqLogger.Error(result.GetError(), "Failed retrieving global meterdefinitions", "CSV", csvName)
							}

							result.Return()
						}

						allMeterDefinitions := append(globalMeterdefinitions, selectedMeterDefinitions...)

						gvk, err := apiutil.GVKForObject(CSV, r.Scheme)
						if err != nil {
							return reconcile.Result{}, err
						}

						// create CSV specific and global meter definitions
						reqLogger.Info("creating meterdefinitions", "CSV", csvName)
						for _, meterDefItem := range allMeterDefinitions {
							reqLogger.Info("checking for existing meterdefinition", "meterdef", meterDefItem.Name, "CSV", csvName)

							// Check if the meterdef is on the cluster already
							meterdef := &marketplacev1beta1.MeterDefinition{}
							err = r.Client.Get(context.TODO(), types.NamespacedName{Name: meterDefItem.Name, Namespace: request.Namespace}, meterdef)
							if err != nil {
								if errors.IsNotFound(err) {
									reqLogger.Info("meterdefinition not found, creating", "meterdef name", meterDefItem.Name, "CSV", CSV.Name)

									// TODO: might have to change this string matching logic later based on global meter def names
									// or we can use an annotation to maintain this info in global meter defs
									isGlobalMeterDef := strings.Contains(meterDefItem.Name, "global")
									result = r.createMeterdefAndInstallMapping(csvName, csvVersion, meterDefItem, CSV, gvk, !isGlobalMeterDef, request, reqLogger)
									if !result.Is(Continue) {

										if result.Is(Error) {
											reqLogger.Error(result.GetError(), "Failed while creating meterdefinition", "meterdef name", meterDefItem.Name, "CSV", CSV.Name)
										}
										return result.Return()
									}

									reqLogger.Info("created meterdefinition", "meterdef name", meterDefItem.Name, "CSV", CSV.Name)
									return reconcile.Result{Requeue: true}, nil
								}

								reqLogger.Error(err, "Failed to get meterdefinition", "meterdef name", meterDefItem.Name, "CSV", CSV.Name)
								return reconcile.Result{}, err
							}
						}
					}
				}
			}
		}
	} else {
		reqLogger.Info("Subscriptions not found in the namespace")
		return reconcile.Result{RequeueAfter: time.Minute * 1}, nil
	}

	reqLogger.Info("reconcilation complete")
	return reconcile.Result{}, nil
}

func reconcileCSV(metaNew metav1.Object) bool {
	ann := metaNew.GetAnnotations()

	ignoreVal, hasIgnoreTag := ann[ignoreTag]

	// we need to pick up the csv
	if !hasIgnoreTag || ignoreVal != ignoreTagValue {
		return true
	}

	//ignore
	return false
}

func ListMeterdefintionsFromFileServer(csvName string, version string, namespace string, client client.Client, kubeInterface kubernetes.Interface, deployedNamespace string, reqLogger logr.Logger) ([]string, []marketplacev1beta1.MeterDefinition, *ExecResult) {
	reqLogger.Info("retrieving meterdefinitions", "csvName", csvName, "csvVersion", version)

	_client, err := createCatalogServerClient(client, deployedNamespace, kubeInterface, reqLogger)
	if err != nil {
		return nil, nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	// returns all the meterdefinitions for associated with a particular CSV version
	url := fmt.Sprintf("https://rhm-meterdefinition-file-server.openshift-redhat-marketplace.svc.cluster.local:8200/list-for-version/%s/%s", csvName, version)
	response, err := _client.Get(url)
	if err != nil {
		reqLogger.Error(err, "Error on GET to Catalog Server")
		if err == io.EOF {
			reqLogger.Error(err, "Meterdefintion not found")
			return nil, nil, &ExecResult{
				ReconcileResult: reconcile.Result{},
				Err:             emperror.New("empty response"),
			}
		}

		reqLogger.Error(err, "Error querying file server")
		return nil, nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	meterDefNames := []string{}
	mdefSlice := []marketplacev1beta1.MeterDefinition{}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		reqLogger.Error(err, "error reading body")
		return nil, nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	if strings.Contains(string(data), "no meterdefinitions found") {
		reqLogger.Info(string(data), "CSV Name", csvName, "CSV Version", version, "CSV Namesapce", namespace)
		return meterDefNames, mdefSlice, &ExecResult{
			Status: ActionResultStatus(Continue),
		}
	}

	reqLogger.Info("response data", "data", string(data))

	meterDefsData := strings.Replace(string(data), "<<NAMESPACE-PLACEHOLDER>>", namespace, -1)
	err = yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(meterDefsData)), 100).Decode(&mdefSlice)
	if err != nil {
		reqLogger.Error(err, "error decoding response from ListMeterdefinitionsFromFileServer()")
		return nil, nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	for _, meterDefItem := range mdefSlice {
		meterDefNames = append(meterDefNames, meterDefItem.ObjectMeta.Name)
	}

	reqLogger.Info("meterdefintions returned from file server", csvName, meterDefNames)

	return meterDefNames, mdefSlice, &ExecResult{
		Status: ActionResultStatus(Continue),
	}
}

func (r *MeterdefinitionInstallReconciler) createMeterdefAndInstallMapping(csvName string, csvVersion string, meterDefinition marketplacev1beta1.MeterDefinition, csv *olmv1alpha1.ClusterServiceVersion, groupVersionKind schema.GroupVersionKind, createInstallMapping bool, request reconcile.Request, reqLogger logr.InfoLogger) *ExecResult {

	// create owner ref object
	ref := metav1.OwnerReference{
		APIVersion:         groupVersionKind.GroupVersion().String(),
		Kind:               groupVersionKind.Kind,
		Name:               csv.GetName(),
		UID:                csv.GetUID(),
		BlockOwnerDeletion: pointer.BoolPtr(false),
		Controller:         pointer.BoolPtr(false),
	}

	meterDefinition.ObjectMeta.SetOwnerReferences([]metav1.OwnerReference{ref})
	meterDefinition.ObjectMeta.Namespace = csv.Namespace

	meterDefName := meterDefinition.Name
	err := r.Client.Create(context.TODO(), &meterDefinition)
	if err != nil {
		reqLogger.Error(err, "Could not create meterdefinition", "mdef", meterDefName, "CSV", csv.Name)
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}
	reqLogger.Info("Created meterdefinition", "mdef", meterDefName, "CSV", csv.Name)

	if createInstallMapping {
		// Fetch the mdefStoreCM instance
		mdefStoreCM := &corev1.ConfigMap{}
		err = r.Client.Get(context.TODO(), types.NamespacedName{Name: utils.METERDEF_INSTALL_MAP_NAME, Namespace: r.cfg.DeployedNamespace}, mdefStoreCM)
		if err != nil {
			if errors.IsNotFound(err) {
				return &ExecResult{
					ReconcileResult: reconcile.Result{},
					Err:             err,
				}

			}

			reqLogger.Error(err, "Failed to get MeterdefintionConfigMap", "mdef", meterDefName, "CSV", csv.Name)
			return &ExecResult{
				ReconcileResult: reconcile.Result{},
				Err:             err,
			}
		}

		mdefStore := mdefStoreCM.Data[utils.MeterDefinitionStoreKey]
		updatedStore, err := createOrUpdateInstallMapping(csvName, csvVersion, meterDefName, mdefStore, request, reqLogger)
		if err != nil {
			return &ExecResult{
				ReconcileResult: reconcile.Result{},
				Err:             err,
			}
		}

		mdefStoreCM.Data[utils.MeterDefinitionStoreKey] = updatedStore

		err = r.Client.Update(context.TODO(), mdefStoreCM)
		if err != nil {
			return &ExecResult{
				ReconcileResult: reconcile.Result{},
				Err:             err,
			}
		}
		reqLogger.Info("Created Install Mapping", "mdef", meterDefName, "CSV", csv.Name)
	}
	return &ExecResult{
		Status: ActionResultStatus(Continue),
	}
}

func (r *MeterdefinitionInstallReconciler) fetchGlobalMeterdefs(csv *olmv1alpha1.ClusterServiceVersion, reqLogger logr.Logger) ([]marketplacev1beta1.MeterDefinition, *ExecResult) {

	csvInfo := CSVInfo{csv.Name, csv.Spec.Version.String(), csv.Namespace}
	globalMeterdefs := []marketplacev1beta1.MeterDefinition{}

	mdefTemplates := r.cfg.ControllerValues.MeterDefTemplates
	for _, _template := range mdefTemplates {
		reqLogger.Info("Started processing global meter definition", "Meterdef Tempalte", _template, "CSV", csv.Name)

		mdef := &marketplacev1beta1.MeterDefinition{}

		templatePath := fmt.Sprintf("meterdef-templates/%s", _template)
		parsedMeterdef, err := template.ParseFS(meterdefTemplates, templatePath)
		if err != nil {
			reqLogger.Error(err, "Error occured while parsing a global meterdefinition", "CSV", csv.Name, "Meterdef YAML", _template)
			return nil, &ExecResult{
				ReconcileResult: reconcile.Result{},
				Err:             err,
			}
		}

		var bufferWriter bytes.Buffer
		err = parsedMeterdef.Execute(&bufferWriter, csvInfo)
		if err != nil {
			reqLogger.Error(err, "Error occured while parsing a global meterdefinition", "CSV", csv.Name, "Meterdef YAML", _template)
			return nil, &ExecResult{
				ReconcileResult: reconcile.Result{},
				Err:             err,
			}
		}

		err = yaml.NewYAMLOrJSONDecoder(bytes.NewReader(bufferWriter.Bytes()), 100).Decode(mdef)
		if err != nil {
			reqLogger.Error(err, "Error occured while decoding a global meterdefinition", "CSV", csv.Name, "Meterdef YAML", _template)
			return nil, &ExecResult{
				ReconcileResult: reconcile.Result{},
				Err:             err,
			}
		}

		globalMeterdefs = append(globalMeterdefs, *mdef)
		reqLogger.Info("Finished processing global meter definition", "Meterdef Tempalte", _template, "CSV", csv.Name)
	}

	reqLogger.Info("Parsed and added all global meter definitions", "Total Global MeterDefs", len(globalMeterdefs), "CSV", csv.Name)
	return globalMeterdefs, &ExecResult{
		Status: ActionResultStatus(Continue),
	}
}

// creates or updates install mapping
func createOrUpdateInstallMapping(csvName string, csvVersion string, installedMeterdef string, cmMdefStore string, request reconcile.Request, reqLogger logr.Logger) (string, error) {

	meterdefStore := &common.MeterdefinitionStore{}

	err := json.Unmarshal([]byte(cmMdefStore), meterdefStore)
	if err != nil {
		reqLogger.Error(err, "error unmarshaling meterdefinition store")
		return "", err
	}

	// if the install mapping already exists update it
	var installMappingFound bool
	for i, item := range meterdefStore.InstallMappings {
		// update a certain package's meterdefinition install list
		if item.CsvName == csvName && item.CsvVersion == csvVersion && item.Namespace == request.Namespace {
			installMappingFound = true
			reqLogger.Info("found existing meterdefinition mapping", "mapping", item)
			// if there are missing meterdefinitions from the list add them
			if !utils.Contains(item.InstalledMeterdefinitions, installedMeterdef) {
				reqLogger.Info("updating installed meterdef list", "updates", installedMeterdef)
				meterdefStore.InstallMappings[i].InstalledMeterdefinitions = append(meterdefStore.InstallMappings[i].InstalledMeterdefinitions, installedMeterdef)
			}
		}
	}

	// if the install mapping isn't found add it
	if !installMappingFound {
		reqLogger.Info("no meterdef mapping found, adding")
		newInstallMapping := common.InstallMapping{
			Namespace:                 request.Namespace,
			CsvName:                   csvName,
			CsvVersion:                csvVersion,
			InstalledMeterdefinitions: []string{installedMeterdef},
		}

		meterdefStore.InstallMappings = append(meterdefStore.InstallMappings, newInstallMapping)
	}

	out, err := json.Marshal(meterdefStore)
	if err != nil {
		return "", err
	}

	meterdefStoreJSON := string(out)
	return meterdefStoreJSON, nil
}

// deletes install mapping for (csvName + namespace) combination
func deleteInstallMapping(csvName string, csvVersion string, deployedNamespace string, request reconcile.Request, client client.Client, reqLogger logr.Logger) *ExecResult {

	// Fetch the mdefStoreCM instance
	mdefStoreCM := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: utils.METERDEF_INSTALL_MAP_NAME, Namespace: deployedNamespace}, mdefStoreCM)
	if err != nil {
		if errors.IsNotFound(err) {
			return &ExecResult{
				ReconcileResult: reconcile.Result{},
				Err:             err,
			}

		}

		reqLogger.Error(err, "Failed to get MeterdefintionConfigMap")
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	mdefStore := mdefStoreCM.Data[utils.MeterDefinitionStoreKey]

	meterdefStore := &common.MeterdefinitionStore{}

	err = json.Unmarshal([]byte(mdefStore), meterdefStore)
	if err != nil {
		reqLogger.Error(err, "error unmarshaling meterdefinition store")
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	var installMappingFound bool
	for i, item := range meterdefStore.InstallMappings {
		// remove the package's InstallMapping. Need to search by csvName here that's what is returned in the request
		if item.CsvName == csvName && item.CsvVersion == csvVersion && item.Namespace == request.Namespace {
			reqLogger.Info("deleting install mapping for", "csv name", csvName, "csv version", csvVersion, "install mapping", item)
			installMappingFound = true
			meterdefStore.InstallMappings = append(meterdefStore.InstallMappings[:i], meterdefStore.InstallMappings[i+1:]...)
			break
		}
	}

	if !installMappingFound {
		reqLogger.Info("deleteInstallMapping: no InstallMap found for", "csv name", csvName, "csv version", csvVersion)
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             nil,
		}
	}

	out, err := json.Marshal(meterdefStore)
	if err != nil {
		reqLogger.Error(err, "error marshaling meterdefinition store")
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	meterdefStoreJSON := string(out)
	mdefStoreCM.Data[utils.MeterDefinitionStoreKey] = meterdefStoreJSON

	err = client.Update(context.TODO(), mdefStoreCM)
	if err != nil {
		return &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	reqLogger.Info("deleted install mapping", "CSVName", csvName, "Namespace", request.Namespace)
	// return &ExecResult{
	// 	Status: ActionResultStatus(Continue),
	// }
	return &ExecResult{
		ReconcileResult: reconcile.Result{},
		Err:             nil,
	}
}

// adds install mapping for (packageName + csvName + csvVersion + namespace) combination
// func addInstallMapping(csvName string, csvVersion, deployedNamespace string, meterDefNames []string, request reconcile.Request, client client.Client, reqLogger logr.Logger) *ExecResult {

// 	// Fetch the mdefStoreCM instance
// 	mdefStoreCM := &corev1.ConfigMap{}
// 	err := client.Get(context.TODO(), types.NamespacedName{Name: utils.METERDEF_INSTALL_MAP_NAME, Namespace: deployedNamespace}, mdefStoreCM)
// 	if err != nil {
// 		if errors.IsNotFound(err) {
// 			return &ExecResult{
// 				ReconcileResult: reconcile.Result{},
// 				Err:             err,
// 			}

// 		}

// 		reqLogger.Error(err, "Failed to get MeterdefintionConfigMap")
// 		return &ExecResult{
// 			ReconcileResult: reconcile.Result{},
// 			Err:             err,
// 		}
// 	}

// 	mdefStore := mdefStoreCM.Data[utils.MeterDefinitionStoreKey]

// 	meterdefStore := &common.MeterdefinitionStore{}

// 	err = json.Unmarshal([]byte(mdefStore), meterdefStore)
// 	if err != nil {
// 		reqLogger.Error(err, "error unmarshaling meterdefinition store")
// 		return &ExecResult{
// 			ReconcileResult: reconcile.Result{},
// 			Err:             err,
// 		}
// 	}

// 	newInstallMapping := common.InstallMapping{
// 		Namespace:                 request.Namespace,
// 		CsvName:                   csvName,
// 		CsvVersion:                csvVersion,
// 		InstalledMeterdefinitions: meterDefNames,
// 	}

// 	meterdefStore.InstallMappings = append(meterdefStore.InstallMappings, newInstallMapping)

// 	out, err := json.Marshal(meterdefStore)
// 	if err != nil {
// 		reqLogger.Error(err, "error marshaling meterdefinition store")
// 		return &ExecResult{
// 			ReconcileResult: reconcile.Result{},
// 			Err:             err,
// 		}
// 	}

// 	meterdefStoreJSON := string(out)
// 	mdefStoreCM.Data[utils.MeterDefinitionStoreKey] = meterdefStoreJSON

// 	err = client.Update(context.TODO(), mdefStoreCM)
// 	if err != nil {
// 		return &ExecResult{
// 			ReconcileResult: reconcile.Result{},
// 			Err:             err,
// 		}
// 	}

// 	reqLogger.Info("added install mapping", "CsvName", csvName, "Namespace", request.Namespace)
// 	return &ExecResult{
// 		ReconcileResult: reconcile.Result{},
// 		Err:             nil,
// 	}
// }

func getCatalogServerService(deployedNamespace string, client client.Client, reqLogger logr.InfoLogger) (*corev1.Service, error) {
	service := &corev1.Service{}

	err := client.Get(context.TODO(), types.NamespacedName{Namespace: deployedNamespace, Name: utils.CATALOG_SERVER_SERVICE_NAME}, service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func getCertFromConfigMap(client client.Client, deployedNamespace string, reqLogger logr.Logger) ([]byte, error) {
	cm := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Namespace: deployedNamespace, Name: "serving-certs-ca-bundle"}, cm)
	if err != nil {
		return nil, err
	}

	reqLogger.Info("extracting cert from config map")

	out, ok := cm.Data["service-ca.crt"]

	if !ok {
		err = emperror.New("Error retrieving cert from config map")
		return nil, err
	}

	cert := []byte(out)
	return cert, nil

}

func createCatalogServerClient(client client.Client, deployedNamespace string, kubeInterface kubernetes.Interface, reqLogger logr.Logger) (*http.Client, error) {
	service, err := getCatalogServerService(deployedNamespace, client, reqLogger)
	if err != nil {
		return nil, &ExecResult{
			ReconcileResult: reconcile.Result{},
			Err:             err,
		}
	}

	cert, err := getCertFromConfigMap(client, deployedNamespace, reqLogger)
	if err != nil {
		return nil, err
	}

	saClient := prom.NewServiceAccountClient(deployedNamespace, kubeInterface)
	authToken, err := saClient.NewServiceAccountToken(utils.OPERATOR_SERVICE_ACCOUNT, utils.FileServerAudience, 3600, reqLogger)
	if err != nil {
		return nil, err
	}

	if service != nil && len(cert) != 0 && authToken != "" {
		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}

		ok := caCertPool.AppendCertsFromPEM(cert)
		if !ok {
			err = emperror.New("failed to append cert to cert pool")
			reqLogger.Error(err, "cert pool error")
			return nil, err
		}

		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}

		var transport http.RoundTripper = &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
		}

		transport = WithBearerAuth(transport, authToken)

		catalogServerClient := http.Client{
			Transport: transport,
			Timeout:   1 * time.Second,
		}

		reqLogger.Info("Catalog Server client created successfully")
		return &catalogServerClient, nil
	}

	return nil, &ExecResult{
		ReconcileResult: reconcile.Result{},
		Err:             emperror.New("catalog server client prerequisites not ready"),
	}

}

//TODO: these funcs are duplicated here, in marketplace_client.go, and prometheus_client.go,may want to create a library
func WithBearerAuth(rt http.RoundTripper, token string) http.RoundTripper {
	addHead := WithHeader(rt)
	addHead.Header.Set("Authorization", "Bearer "+token)
	return addHead
}

type withHeader struct {
	http.Header
	rt http.RoundTripper
}

func WithHeader(rt http.RoundTripper) withHeader {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return withHeader{Header: make(http.Header), rt: rt}
}

func (h withHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}

func checkForCSVVersionChanges(e event.UpdateEvent) bool {
	oldCSV, ok := e.ObjectOld.(*olmv1alpha1.ClusterServiceVersion)
	if !ok {
		return false
	}

	newCSV, ok := e.ObjectNew.(*olmv1alpha1.ClusterServiceVersion)
	if !ok {
		return false
	}

	return oldCSV.Spec.Version.String() != newCSV.Spec.Version.String()
}

var rhmCSVControllerPredicates predicate.Funcs = predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		if !reconcileCSV(e.MetaNew) || !reconcileCSV(e.MetaOld) {
			return false
		}
		return checkForCSVVersionChanges(e)
	},

	DeleteFunc: func(e event.DeleteEvent) bool {
		return reconcileCSV(e.Meta)
	},

	CreateFunc: func(e event.CreateEvent) bool {
		return reconcileCSV(e.Meta)

	},

	GenericFunc: func(e event.GenericEvent) bool {
		return reconcileCSV(e.Meta)
	},
}

func (r *MeterdefinitionInstallReconciler) Inject(injector mktypes.Injectable) mktypes.SetupWithManager {
	injector.SetCustomFields(r)
	return r
}

func (m *MeterdefinitionInstallReconciler) InjectOperatorConfig(cfg *config.OperatorConfig) error {
	m.cfg = cfg
	return nil
}

func (r *MeterdefinitionInstallReconciler) InjectKubeInterface(k kubernetes.Interface) error {
	r.kubeInterface = k
	return nil
}

func (r *MeterdefinitionInstallReconciler) SetupWithManager(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&olmv1alpha1.ClusterServiceVersion{}, builder.WithPredicates(rhmCSVControllerPredicates)).
		Complete(r)
}
