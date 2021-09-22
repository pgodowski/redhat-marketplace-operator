package marketplace

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/blang/semver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/onsi/gomega/ghttp"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/gotidy/ptr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	osappsv1 "github.com/openshift/api/apps/v1"
	osimagev1 "github.com/openshift/api/image/v1"
	"github.com/operator-framework/api/pkg/lib/version"
	olmv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/common"
	marketplacev1alpha1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1alpha1"
	marketplacev1beta1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/catalog"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = FDescribe("DeploymentConfig Controller Test", func() {

	var (
		namespace = "default"

		/* rhm csv */
		csvName = "test-csv-1.v0.0.1"
		csvVersion        = "0.0.1"
		subName           = "test-csv-1-sub"
		packageName       = "test-csv-1-rhmp"
		catalogSourceName = "redhat-marketplace"

		/* non-rhm csv */
		nonRhmCsvName           = "non-rhm-csv.v0.0.1"
		nonRhmSubNmae           = "non-rhm-sub"
		nonRhmPackageName       = "non-rhm-package"
		nonRhmCatalogSourceName = "non-rhm-catalog-source"

		/* system meterdefs */
		systemMeterDef1Name  = csvName + "-" + "pod-count"
		systemMeterDef2Name = csvName + "-" + "cpu-usage"

		listMeterDefsForCsvPath          = "/" + catalog.ListForVersionEndpoint + "/" + csvName + "/" + csvVersion + "/" + namespace
		indexLabelsPath                  = "/" + catalog.GetMeterdefinitionIndexLabelEndpoint + "/" + csvName
		systemMeterDefIndexLabelsPath    = "/" + catalog.GetSystemMeterDefIndexLabelEndpoint + "/" + csvName
		healthEndpoint 				     = "/" + catalog.HealthEndpoint
		communityMeterDefIndexLabelsBody []byte
		systemMeterDefIndexLabelsBody    []byte
		dcControllerMockServer           *ghttp.Server
	)

	idFn := func(element interface{}) string {
		return fmt.Sprintf("%v", element)
	}

	meterDef1Key := types.NamespacedName{
		Name:      "meterdef-1",
		Namespace: namespace,
	}

	meterDef2Key := types.NamespacedName{
		Name:      "meterdef-2",
		Namespace: namespace,
	}

	systemMeterDef1Key := types.NamespacedName{
		Name:      systemMeterDef1Name,
		Namespace: namespace,
	}

	systemMeterDef2Key := types.NamespacedName{
		Name:      systemMeterDef2Name,
		Namespace: namespace,
	}

	rhmCsvKey  := types.NamespacedName{
		Name:      csvName,
		Namespace: namespace,
	}

	meterBaseKey := types.NamespacedName{
		Name: utils.METERBASE_NAME, 
		Namespace: namespace,
	}

	systemMeterDef1 := marketplacev1beta1.MeterDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MeterDefinition",
			APIVersion: "marketplace.redhat.com/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      systemMeterDef1Name,
			Namespace: "default",
			Annotations: map[string]string{
				"versionRange": "<=0.0.1",
			},
			Labels: map[string]string{
				"marketplace.redhat.com/installedOperatorNameTag": csvName,
				"marketplace.redhat.com/isSystemMeterDefinition":  "1",
			},
		},
		Spec: marketplacev1beta1.MeterDefinitionSpec{
			Group: "apps.partner.metering.com",
			Kind:  "App",
			ResourceFilters: []marketplacev1beta1.ResourceFilter{
				{
					WorkloadType: marketplacev1beta1.WorkloadTypeService,
					OwnerCRD: &marketplacev1beta1.OwnerCRDFilter{
						GroupVersionKind: common.GroupVersionKind{
							APIVersion: "test_package_1.com/v2",
							Kind:       "test_package_1Cluster",
						},
					},
					Namespace: &marketplacev1beta1.NamespaceFilter{
						UseOperatorGroup: true,
					},
				},
			},
			Meters: []marketplacev1beta1.MeterWorkload{
				{
					Aggregation: "sum",
					GroupBy:     []string{"namespace"},
					Period: &metav1.Duration{
						Duration: time.Duration(time.Hour * 1),
					},
					Query:        "kube_service_labels{}",
					Metric:       "test_package_1_cluster_count",
					WorkloadType: marketplacev1beta1.WorkloadTypeService,
					Without:      []string{"label_test_package_1_cluster", "label_app", "label_operator_test_package_1_com_version"},
				},
			},
		},
	}

	systemMeterDef2 := marketplacev1beta1.MeterDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MeterDefinition",
			APIVersion: "marketplace.redhat.com/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      systemMeterDef2Name,
			Namespace: "default",
			Annotations: map[string]string{
				"versionRange": "<=0.0.1",
			},
			Labels: map[string]string{
				"marketplace.redhat.com/installedOperatorNameTag": csvName,
				"marketplace.redhat.com/isSystemMeterDefinition":  "1",
			},
		},
		Spec: marketplacev1beta1.MeterDefinitionSpec{
			Group: "apps.partner.metering.com",
			Kind:  "App",
			ResourceFilters: []marketplacev1beta1.ResourceFilter{
				{
					WorkloadType: marketplacev1beta1.WorkloadTypeService,
					OwnerCRD: &marketplacev1beta1.OwnerCRDFilter{
						GroupVersionKind: common.GroupVersionKind{
							APIVersion: "test_package_1.com/v2",
							Kind:       "test_package_1Cluster",
						},
					},
					Namespace: &marketplacev1beta1.NamespaceFilter{
						UseOperatorGroup: true,
					},
				},
			},
			Meters: []marketplacev1beta1.MeterWorkload{
				{
					Aggregation: "sum",
					GroupBy:     []string{"namespace"},
					Period: &metav1.Duration{
						Duration: time.Duration(time.Hour * 1),
					},
					Query:        "kube_service_labels{}",
					Metric:       "test_package_1_cluster_count",
					WorkloadType: marketplacev1beta1.WorkloadTypeService,
					Without:      []string{"label_test_package_1_cluster", "label_app", "label_operator_test_package_1_com_version"},
				},
			},
		},
	}

	meterDef1 := marketplacev1beta1.MeterDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MeterDefinition",
			APIVersion: "marketplace.redhat.com/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      meterDef1Key.Name,
			Namespace: namespace,
			Annotations: map[string]string{
				"versionRange": "<=0.0.1",
			},
			Labels: map[string]string{
				"marketplace.redhat.com/installedOperatorNameTag":  csvName,
				"marketplace.redhat.com/isCommunityMeterdefintion": "1",
			},
		},
		Spec: marketplacev1beta1.MeterDefinitionSpec{
			Group: "apps.partner.metering.com",
			Kind:  "App",
			ResourceFilters: []marketplacev1beta1.ResourceFilter{
				{
					WorkloadType: marketplacev1beta1.WorkloadTypeService,
					OwnerCRD: &marketplacev1beta1.OwnerCRDFilter{
						GroupVersionKind: common.GroupVersionKind{
							APIVersion: "test_package_1.com/v2",
							Kind:       "test_package_1Cluster",
						},
					},
					Namespace: &marketplacev1beta1.NamespaceFilter{
						UseOperatorGroup: true,
					},
				},
			},
			Meters: []marketplacev1beta1.MeterWorkload{
				{
					Aggregation: "sum",
					GroupBy:     []string{"namespace"},
					Period: &metav1.Duration{
						Duration: time.Duration(time.Hour * 1),
					},
					Query:        "kube_service_labels{}",
					Metric:       "test_package_1_cluster_count",
					WorkloadType: marketplacev1beta1.WorkloadTypeService,
					Without:      []string{"label_test_package_1_cluster", "label_app", "label_operator_test_package_1_com_version"},
				},
			},
		},
	}

	meterDef2 := marketplacev1beta1.MeterDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MeterDefinition",
			APIVersion: "marketplace.redhat.com/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      meterDef2Key.Name,
			Namespace: "default",
			Annotations: map[string]string{
				"versionRange": "<=0.0.1",
			},
			Labels: map[string]string{
				"marketplace.redhat.com/installedOperatorNameTag":  csvName,
				"marketplace.redhat.com/isCommunityMeterdefintion": "1",
			},
		},
		Spec: marketplacev1beta1.MeterDefinitionSpec{
			Group: "apps.partner.metering.com",
			Kind:  "App",
			ResourceFilters: []marketplacev1beta1.ResourceFilter{
				{
					WorkloadType: marketplacev1beta1.WorkloadTypeService,
					OwnerCRD: &marketplacev1beta1.OwnerCRDFilter{
						GroupVersionKind: common.GroupVersionKind{
							APIVersion: "test_package_1.com/v2",
							Kind:       "test_package_1Cluster",
						},
					},
					Namespace: &marketplacev1beta1.NamespaceFilter{
						UseOperatorGroup: true,
					},
				},
			},
			Meters: []marketplacev1beta1.MeterWorkload{
				{
					Aggregation: "sum",
					GroupBy:     []string{"namespace"},
					Period: &metav1.Duration{
						Duration: time.Duration(time.Hour * 1),
					},
					Query:        "kube_service_labels{}",
					Metric:       "test_package_1_cluster_count",
					WorkloadType: marketplacev1beta1.WorkloadTypeService,
					Without:      []string{"label_test_package_1_cluster", "label_app", "label_operator_test_package_1_com_version"},
				},
			},
		},
	}

	subSectionMeterBase := &marketplacev1alpha1.MeterBase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.METERBASE_NAME,
			Namespace: namespace,
		},
		Spec: marketplacev1alpha1.MeterBaseSpec{
			Enabled: false,
			Prometheus: &marketplacev1alpha1.PrometheusSpec{
				Storage: marketplacev1alpha1.StorageSpec{
					Size: resource.MustParse("30Gi"),
				},
				Replicas: ptr.Int32(2),
			},
			MeterdefinitionCatalogServerConfig: &common.MeterDefinitionCatalogServerConfig{
				SyncCommunityMeterDefinitions:      true,
				SyncSystemMeterDefinitions:         true,
				DeployMeterDefinitionCatalogServer: true,
			},
		},
	}

	catalogSource := &olmv1alpha1.CatalogSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      catalogSourceName,
			Namespace: namespace,
		},
		Spec: olmv1alpha1.CatalogSourceSpec{
			SourceType: olmv1alpha1.SourceType(olmv1alpha1.SourceTypeGrpc),
			Image:      "quay.io/mxpaspa/memcached-ansible-index:1.0.1",
		},
	}

	csvOnCluster := olmv1alpha1.ClusterServiceVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      csvName,
			Namespace: namespace,
			Annotations: map[string]string{
				"operatorframework.io/properties": fmt.Sprintf(`{"properties":[{"type":"olm.gvk","value":{"group":"app.joget.com","kind":"JogetDX","version":"v1alpha1"}},{"type":"olm.package","value":{"packageName":"%v","version":"0.0.1"}}]}`, packageName),
			},
		},
		Spec: olmv1alpha1.ClusterServiceVersionSpec{
			InstallStrategy: olmv1alpha1.NamedInstallStrategy{
				StrategySpec: olmv1alpha1.StrategyDetailsDeployment{
					DeploymentSpecs: []olmv1alpha1.StrategyDeploymentSpec{},
				},
			},
			Version: version.OperatorVersion{
				Version: semver.Version{
					Major: 0,
					Minor: 0,
					Patch: 1,
				},
			},
		},
		Status: olmv1alpha1.ClusterServiceVersionStatus{},
	}

	nonRhmCsv := olmv1alpha1.ClusterServiceVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nonRhmCsvName,
			Namespace: namespace,
			Annotations: map[string]string{
				"operatorframework.io/properties": fmt.Sprintf(`{"properties":[{"type":"olm.gvk","value":{"group":"app.joget.com","kind":"JogetDX","version":"v1alpha1"}},{"type":"olm.package","value":{"packageName":"%v","version":"0.0.1"}}]}`, nonRhmPackageName),
			},
		},
		Spec: olmv1alpha1.ClusterServiceVersionSpec{
			InstallStrategy: olmv1alpha1.NamedInstallStrategy{
				StrategySpec: olmv1alpha1.StrategyDetailsDeployment{
					DeploymentSpecs: []olmv1alpha1.StrategyDeploymentSpec{},
				},
			},
			Version: version.OperatorVersion{
				Version: semver.Version{
					Major: 0,
					Minor: 0,
					Patch: 1,
				},
			},
		},
		Status: olmv1alpha1.ClusterServiceVersionStatus{},
	}

	subs := []olmv1alpha1.Subscription{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subName,
				Namespace: namespace,
				Labels: map[string]string{
					operatorTag: "true",
				},
			},
			Spec: &olmv1alpha1.SubscriptionSpec{
				CatalogSource:          catalogSourceName,
				CatalogSourceNamespace: namespace,
				Package:                packageName,
			},
			Status: olmv1alpha1.SubscriptionStatus{
				InstalledCSV: csvName,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      nonRhmSubNmae,
				Namespace: namespace,
				Labels: map[string]string{
					operatorTag: "true",
				},
			},
			Spec: &olmv1alpha1.SubscriptionSpec{
				CatalogSource:          nonRhmCatalogSourceName,
				CatalogSourceNamespace: namespace,
				Package:                nonRhmPackageName,
			},
			Status: olmv1alpha1.SubscriptionStatus{
				InstalledCSV: nonRhmCsvName,
			},
		},
	}

	Status200 := 200
	systemMeterdefsPath := "/" + catalog.GetSystemMeterdefinitionTemplatesEndpoint

	BeforeEach(func() {
		customListener, err := net.Listen("tcp", listenerAddress)
		Expect(err).ToNot(HaveOccurred())

		dcControllerMockServer = ghttp.NewUnstartedServer()
		dcControllerMockServer.HTTPTestServer.Listener.Close()
		dcControllerMockServer.HTTPTestServer.Listener = customListener
		dcControllerMockServer.SetAllowUnhandledRequests(true)
		dcControllerMockServer.Start()

		dc := &osappsv1.DeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      utils.DeploymentConfigName,
				Namespace: "default",
			},
			Spec: osappsv1.DeploymentConfigSpec{
				Triggers: osappsv1.DeploymentTriggerPolicies{
					{
						Type: osappsv1.DeploymentTriggerOnConfigChange,
						ImageChangeParams: &osappsv1.DeploymentTriggerImageChangeParams{
							Automatic:      true,
							ContainerNames: []string{"rhm-meterdefinition-file-server"},
							From: corev1.ObjectReference{
								Kind: "ImageStreamTag",
								Name: "rhm-meterdefinition-file-server:v1",
							},
						},
					},
				},
			},
			Status: osappsv1.DeploymentConfigStatus{
				LatestVersion: 1,
				Conditions: []osappsv1.DeploymentCondition{
					{
						Type:               osappsv1.DeploymentConditionType(osappsv1.DeploymentAvailable),
						Reason:             "NewReplicationControllerAvailable",
						Status:             corev1.ConditionTrue,
						LastTransitionTime: metav1.Now(),
						LastUpdateTime:     metav1.Now(),
					},
				},
			},
		}

		is := &osimagev1.ImageStream{
			ObjectMeta: metav1.ObjectMeta{
				Name:      utils.DeploymentConfigName,
				Namespace: "default",
			},
			Spec: osimagev1.ImageStreamSpec{
				LookupPolicy: osimagev1.ImageLookupPolicy{
					Local: false,
				},
				Tags: []osimagev1.TagReference{
					{
						Annotations: map[string]string{
							"openshift.io/imported-from": "quay.io/mxpaspa/rhm-meterdefinition-file-server:return-204-1.0.0",
						},
						From: &corev1.ObjectReference{
							Name: "quay.io/mxpaspa/rhm-meterdefinition-file-server:return-204-1.0.0",
							Kind: "DockerImage",
						},
						ImportPolicy: osimagev1.TagImportPolicy{
							Insecure:  true,
							Scheduled: true,
						},
						Name: "v1",
						ReferencePolicy: osimagev1.TagReferencePolicy{
							Type: osimagev1.SourceTagReferencePolicy,
						},
						Generation: ptr.Int64(1),
					},
				},
			},
		}

		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      utils.DeploymentConfigName,
				Namespace: namespace,
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:       "foo",
						Port:       int32(8180),
						TargetPort: intstr.FromString("foo"),
					},
				},
			},
		}

		Expect(k8sClient.Create(context.TODO(), dc)).Should(SucceedOrAlreadyExist, "create test deploymentconfig")
		Expect(k8sClient.Create(context.TODO(), is)).Should(SucceedOrAlreadyExist, "create test image stream")
		Expect(k8sClient.Create(context.TODO(), service)).Should(SucceedOrAlreadyExist, "create file server service")
		Expect(k8sClient.Create(context.TODO(), catalogSource.DeepCopy())).Should(Succeed(), "create catalog source")

		communityMeterDefIndexLabelsBody = []byte(fmt.Sprintf(`{
				"marketplace.redhat.com/installedOperatorNameTag": "%v",
				"marketplace.redhat.com/isCommunityMeterdefintion": "1"
			}`, csvName))

		dcControllerMockServer.RouteToHandler(
			"GET", indexLabelsPath, ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", indexLabelsPath),
				ghttp.RespondWithPtr(&Status200, &communityMeterDefIndexLabelsBody),
			))

		systemMeterDefIndexLabelsBody = []byte(fmt.Sprintf(`{
				"marketplace.redhat.com/installedOperatorNameTag": "%v",
				"marketplace.redhat.com/isSystemMeterDefinition": "1"
			}`, csvName))

		dcControllerMockServer.RouteToHandler(
			"GET", systemMeterDefIndexLabelsPath, ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", systemMeterDefIndexLabelsPath),
				ghttp.RespondWithPtr(&Status200, &systemMeterDefIndexLabelsBody),
			))
		
		healthBody := []byte(`status ok`)

		dcControllerMockServer.RouteToHandler(
			"GET", healthEndpoint, ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", healthEndpoint),
				ghttp.RespondWith(http.StatusOK,healthBody),
			))
	})

	AfterEach(func() {
		dcControllerMockServer.Close()

		_meterDef1 := &marketplacev1beta1.MeterDefinition{}
		k8sClient.Get(context.TODO(), meterDef1Key, _meterDef1)
		k8sClient.Delete(context.TODO(), _meterDef1)

		_meterDef2 := &marketplacev1beta1.MeterDefinition{}
		k8sClient.Get(context.TODO(), meterDef2Key, _meterDef2)
		k8sClient.Delete(context.TODO(), _meterDef2)

		_systemMeterDef1 := &marketplacev1beta1.MeterDefinition{}
		k8sClient.Get(context.TODO(), systemMeterDef1Key, _systemMeterDef1)
		k8sClient.Delete(context.TODO(), _systemMeterDef1)

		_systemMeterDef2 := &marketplacev1beta1.MeterDefinition{}
		k8sClient.Get(context.TODO(), systemMeterDef2Key, _systemMeterDef2)
		k8sClient.Delete(context.TODO(), _systemMeterDef2)

		_csv := &olmv1alpha1.ClusterServiceVersion{}
		k8sClient.Get(context.TODO(), rhmCsvKey, _csv)
		k8sClient.Delete(context.TODO(), _csv)

		_meterBase := &marketplacev1alpha1.MeterBase{}
		Expect(k8sClient.Get(context.TODO(), meterBaseKey, _meterBase)).Should(Succeed(), "get meterbase")
		k8sClient.Delete(context.TODO(), _meterBase)

		_catalogSource := &olmv1alpha1.CatalogSource{}
		Expect(k8sClient.Get(context.TODO(), types.NamespacedName{Name: catalogSourceName, Namespace: namespace}, _catalogSource)).Should(Succeed(), "get catalogsource")
		k8sClient.Delete(context.TODO(), _catalogSource)
	})

	Context("Create", func() {
		BeforeEach(func() {
			listSubs = func(k8sclient client.Client, csv *olmv1alpha1.ClusterServiceVersion) ([]olmv1alpha1.Subscription, error) {
				return subs, nil
			}

			Expect(k8sClient.Create(context.TODO(), csvOnCluster.DeepCopy())).Should(Succeed(), "create csv on cluster")
			Expect(k8sClient.Create(context.TODO(), subSectionMeterBase.DeepCopy())).Should(Succeed(), "create sub-section meterbase")

			returnedCommunityMeterdefGoSlice := []marketplacev1beta1.MeterDefinition{*meterDef1.DeepCopy()}
			communityMeterDefsBody, err := json.Marshal(returnedCommunityMeterdefGoSlice)
			if err != nil {
				log.Fatal(err)
			}

			dcControllerMockServer.RouteToHandler(
				"GET", listMeterDefsForCsvPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", listMeterDefsForCsvPath),
					ghttp.RespondWithPtr(&Status200, &communityMeterDefsBody),
				))

			returnedSystemMeterDefSlice := []marketplacev1beta1.MeterDefinition{*systemMeterDef1.DeepCopy()}
			systemMeterDefBody, err := json.Marshal(returnedSystemMeterDefSlice)
			if err != nil {
				log.Fatal(err)
			}

			dcControllerMockServer.RouteToHandler(
				"POST", systemMeterdefsPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", systemMeterdefsPath),
					ghttp.RespondWithPtr(&Status200, &systemMeterDefBody),
				))
		})

		It("Should create community defs if listed in the catalog", func() {
			Eventually(func() string {
				found := &marketplacev1beta1.MeterDefinition{}
				k8sClient.Get(context.TODO(), meterDef1Key, found)
				return found.Name
			}, timeout, interval).Should(Equal(meterDef1Key.Name))
		})

		It("Should create a system meterdef", func() {
			Eventually(func() string {
				foundSystemMeterDef := &marketplacev1beta1.MeterDefinition{}
				k8sClient.Get(context.TODO(), systemMeterDef1Key, foundSystemMeterDef)
				return foundSystemMeterDef.Name
			}, timeout, interval).Should(Equal(systemMeterDef1.Name))
		})
	})

	Context("Update", func() {
		BeforeEach(func() {
			listSubs = func(k8sclient client.Client, csv *olmv1alpha1.ClusterServiceVersion) ([]olmv1alpha1.Subscription, error) {
				return subs, nil
			}

			Expect(k8sClient.Create(context.TODO(), csvOnCluster.DeepCopy())).Should(Succeed(), "create csv on cluster")

			existingMeterDef := meterDef1.DeepCopy()
			Expect(k8sClient.Create(context.TODO(), existingMeterDef)).Should(Succeed(), "create existing meterdef")
			Expect(k8sClient.Create(context.TODO(), subSectionMeterBase.DeepCopy())).Should(Succeed(), "create sub-section meterbase")

			updatedMeterDef := meterDef1.DeepCopy()
			updatedMeterDef.Spec.Meters[0].Name = "updated"
			updatedCommunityMeterdefGoSlice := []marketplacev1beta1.MeterDefinition{*updatedMeterDef}
			communityMeterDefsBody, err := json.Marshal(updatedCommunityMeterdefGoSlice)
			if err != nil {
				log.Fatal(err)
			}

			dcControllerMockServer.RouteToHandler(
				"GET", listMeterDefsForCsvPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", listMeterDefsForCsvPath),
					ghttp.RespondWithPtr(&Status200, &communityMeterDefsBody),
				))

			existingSystemMeterDef := systemMeterDef1.DeepCopy()
			Expect(k8sClient.Create(context.TODO(), existingSystemMeterDef)).Should(Succeed(), "create existing system meterdef")

			updatedSystemMeterDef := systemMeterDef1.DeepCopy()
			updatedSystemMeterDef.Spec.Meters[0].Name = "updated system meterdef"
			updatedSystemMeterDefSlice := []marketplacev1beta1.MeterDefinition{*updatedSystemMeterDef}
			systemMeterDefBody, err := json.Marshal(updatedSystemMeterDefSlice)
			if err != nil {
				log.Fatal(err)
			}

			//overwriting RouteToHandler in top-level BeforeEach()
			dcControllerMockServer.RouteToHandler(
				"POST", systemMeterdefsPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", systemMeterdefsPath),
					ghttp.RespondWithPtr(&Status200, &systemMeterDefBody),
				))
		})

		It("community meterdefintion on cluster should be updated", func() {
			Eventually(func() string {
				found := &marketplacev1beta1.MeterDefinition{}
				k8sClient.Get(context.TODO(), meterDef1Key, found)

				if found.Spec.Meters != nil {
					return found.Spec.Meters[0].Name
				}

				return ""
			}, timeout, interval).Should(Equal("updated"))
		})

		It("system meterdefintion on cluster should be updated", func() {
			Eventually(func() string {
				foundSystemMeterDef := &marketplacev1beta1.MeterDefinition{}
				k8sClient.Get(context.TODO(), systemMeterDef1Key, foundSystemMeterDef)

				if foundSystemMeterDef.Spec.Meters != nil {
					return foundSystemMeterDef.Spec.Meters[0].Name
				}

				return ""
			}, timeout, interval).Should(Equal("updated system meterdef"))
		})
	})

	Context("Delete", func() {
		BeforeEach(func() {
			listSubs = func(k8sclient client.Client, csv *olmv1alpha1.ClusterServiceVersion) ([]olmv1alpha1.Subscription, error) {
				return subs, nil
			}

			Expect(k8sClient.Create(context.TODO(), csvOnCluster.DeepCopy())).Should(Succeed(), "create csv on cluster")

			_meterDef1 := meterDef1.DeepCopy()
			_meterDef2 := meterDef2.DeepCopy()
			_systemMeterDef1 := systemMeterDef1.DeepCopy()
			_systemMeterDef2 := systemMeterDef2.DeepCopy()

			existingMeterdefSlice := []marketplacev1beta1.MeterDefinition{*_meterDef1, *_meterDef2, *_systemMeterDef1, *_systemMeterDef2}
			for _, existingMeterdef := range existingMeterdefSlice {
				Expect(k8sClient.Create(context.TODO(), &existingMeterdef)).Should(Succeed(), "create existing meterdefs")
			}

			Expect(k8sClient.Create(context.TODO(), subSectionMeterBase.DeepCopy())).Should(Succeed(), "create sub-section meterbase")

			latestCommunityMeterdefsFromCatalog := []marketplacev1beta1.MeterDefinition{*_meterDef1}
			communityMeterDefsBody, err := json.Marshal(latestCommunityMeterdefsFromCatalog)
			if err != nil {
				log.Fatal(err)
			}

			dcControllerMockServer.RouteToHandler(
				"GET", listMeterDefsForCsvPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", listMeterDefsForCsvPath),
					ghttp.RespondWithPtr(&Status200, &communityMeterDefsBody),
				))

			returnedSystemMeterDefSlice := []marketplacev1beta1.MeterDefinition{*systemMeterDef1.DeepCopy()}
			systemMeterDefBody, err := json.Marshal(returnedSystemMeterDefSlice)
			if err != nil {
				log.Fatal(err)
			}

			dcControllerMockServer.RouteToHandler(
				"POST", systemMeterdefsPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", systemMeterdefsPath),
					ghttp.RespondWithPtr(&Status200, &systemMeterDefBody),
				))
		})

		It("meterdef-2 should be deleted off the cluster if removed from the catalog", func() {
			Eventually(func() []string {

				labelsMap := map[string]string{}
				err := json.Unmarshal(communityMeterDefIndexLabelsBody, &labelsMap)
				if err != nil {
					log.Fatal(err)
				}

				listOts := []client.ListOption{
					client.MatchingLabels(labelsMap),
				}

				mdefList := &marketplacev1beta1.MeterDefinitionList{}
				k8sClient.List(context.TODO(), mdefList, listOts...)

				var mdefNames []string
				for _, mdef := range mdefList.Items {
					mdefNames = append(mdefNames, mdef.Name)
				}

				return mdefNames
			}, timeout, interval).Should(And(
				HaveLen(1),
				MatchAllElements(idFn, Elements{
					"meterdef-1": Equal("meterdef-1"),
				}),
			))
		})

		It("systemmeterdef should be deleted off the cluster if removed from the catalog", func() {
			Eventually(func() []string {

				labelsMap := map[string]string{}
				err := json.Unmarshal(systemMeterDefIndexLabelsBody, &labelsMap)
				if err != nil {
					log.Fatal(err)
				}

				listOts := []client.ListOption{
					client.MatchingLabels(labelsMap),
				}

				mdefList := &marketplacev1beta1.MeterDefinitionList{}
				k8sClient.List(context.TODO(), mdefList, listOts...)

				var mdefNames []string
				for _, mdef := range mdefList.Items {
					mdefNames = append(mdefNames, mdef.Name)
				}

				return mdefNames
			}, timeout, interval).Should(And(
				HaveLen(1),
				MatchAllElements(idFn, Elements{
					systemMeterDef1Name: Equal(systemMeterDef1Name),
				}),
			))
		})
	})

	Context("Remove Catalog directory", func() {
		BeforeEach(func() {
			listSubs = func(k8sclient client.Client, csv *olmv1alpha1.ClusterServiceVersion) ([]olmv1alpha1.Subscription, error) {
				return subs, nil
			}

			Expect(k8sClient.Create(context.TODO(), csvOnCluster.DeepCopy())).Should(Succeed(), "create csv on cluster")

			_meterDef1 := meterDef1.DeepCopy()
			_meterDef2 := meterDef2.DeepCopy()

			existingMeterdefSlice := []marketplacev1beta1.MeterDefinition{*_meterDef1, *_meterDef2}

			for _, existingMeterdef := range existingMeterdefSlice {
				Expect(k8sClient.Create(context.TODO(), &existingMeterdef)).Should(SucceedOrAlreadyExist, "create existing meterdefs")
			}

			Expect(k8sClient.Create(context.TODO(), subSectionMeterBase.DeepCopy())).Should(Succeed(), "create sub-section meterbase")

			notFoundBody := []byte(`no meterdefs found`)
			dcControllerMockServer.RouteToHandler(
				"GET", listMeterDefsForCsvPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", listMeterDefsForCsvPath),
					ghttp.RespondWith(http.StatusNoContent, notFoundBody),
				))

			returnedSystemMeterDefSlice := []marketplacev1beta1.MeterDefinition{*systemMeterDef1.DeepCopy()}
			systemMeterDefBody, err := json.Marshal(returnedSystemMeterDefSlice)
			if err != nil {
				log.Fatal(err)
			}

			dcControllerMockServer.RouteToHandler(
				"POST", systemMeterdefsPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", systemMeterdefsPath),
					ghttp.RespondWithPtr(&Status200, &systemMeterDefBody),
				))
		})

		It("should delete all community meterdefinitions for an rhm csv with no catalog listing", func() {
			Eventually(func() []marketplacev1beta1.MeterDefinition {

				labelsMap := map[string]string{}
				err := json.Unmarshal(communityMeterDefIndexLabelsBody, &labelsMap)
				if err != nil {
					log.Fatal(err)
				}
				listOts := []client.ListOption{
					client.MatchingLabels(labelsMap),
				}

				mdefList := &marketplacev1beta1.MeterDefinitionList{}
				k8sClient.List(context.TODO(), mdefList, listOts...)

				return mdefList.Items
			}, timeout, interval).Should(HaveLen(0))
		})
	})

	Context("non-rhm resources", func() {
		BeforeEach(func() {
			listSubs = func(k8sclient client.Client, csv *olmv1alpha1.ClusterServiceVersion) ([]olmv1alpha1.Subscription, error) {
				return subs, nil
			}

			Expect(k8sClient.Create(context.TODO(), nonRhmCsv.DeepCopy())).Should(Succeed(), "create non-rhm-csv")
			Expect(k8sClient.Create(context.TODO(), subSectionMeterBase.DeepCopy())).Should(Succeed(), "create sub-section meterbase")

		})

		It("it should not create meterdefs for non-rhm resources", func() {
			Eventually(func() []marketplacev1beta1.MeterDefinition {
				mdefList := &marketplacev1beta1.MeterDefinitionList{}
				k8sClient.List(context.TODO(), mdefList)

				return mdefList.Items
			}, timeout, interval).Should(HaveLen(0))
		})
	})

	Context("feature flags", func() {
		BeforeEach(func() {
			listSubs = func(k8sclient client.Client, csv *olmv1alpha1.ClusterServiceVersion) ([]olmv1alpha1.Subscription, error) {
				return subs, nil
			}

			Expect(k8sClient.Create(context.TODO(), csvOnCluster.DeepCopy())).Should(Succeed(), "create csv on cluster")
			Expect(k8sClient.Create(context.TODO(), subSectionMeterBase.DeepCopy())).Should(Succeed(), "create sub-section meterbase")

			returnedCommunityMeterdefGoSlice := []marketplacev1beta1.MeterDefinition{*meterDef1.DeepCopy()}
			communityMeterDefsBody, err := json.Marshal(returnedCommunityMeterdefGoSlice)
			if err != nil {
				log.Fatal(err)
			}

			dcControllerMockServer.RouteToHandler(
				"GET", listMeterDefsForCsvPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", listMeterDefsForCsvPath),
					ghttp.RespondWithPtr(&Status200, &communityMeterDefsBody),
				))

			returnedSystemMeterDefSlice := []marketplacev1beta1.MeterDefinition{*systemMeterDef1.DeepCopy()}
			systemMeterDefBody, err := json.Marshal(returnedSystemMeterDefSlice)
			if err != nil {
				log.Fatal(err)
			}

			dcControllerMockServer.RouteToHandler(
				"POST", systemMeterdefsPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", systemMeterdefsPath),
					ghttp.RespondWithPtr(&Status200, &systemMeterDefBody),
				))
		})

		It("all system meterdefinitions should be deleted if SyncSystemMeterDefinitions is disabled", func() {
			Eventually(func() []marketplacev1beta1.MeterDefinition {
				labelsMap := map[string]string{}
				err := json.Unmarshal(communityMeterDefIndexLabelsBody, &labelsMap)
				if err != nil {
					log.Fatal(err)
				}

				listOts := []client.ListOption{
					client.MatchingLabels(labelsMap),
				}

				mdefList := &marketplacev1beta1.MeterDefinitionList{}
				k8sClient.List(context.TODO(), mdefList, listOts...)

				return mdefList.Items
			}, timeout, interval).Should(HaveLen(1), "community meterdefs should get created")

			Eventually(func() []marketplacev1beta1.MeterDefinition {
				labelsMap := map[string]string{}
				err := json.Unmarshal(systemMeterDefIndexLabelsBody, &labelsMap)
				if err != nil {
					log.Fatal(err)
				}

				listOts := []client.ListOption{
					client.MatchingLabels(labelsMap),
				}

				mdefList := &marketplacev1beta1.MeterDefinitionList{}
				k8sClient.List(context.TODO(), mdefList, listOts...)

				return mdefList.Items
			}, timeout, interval).Should(HaveLen(1), "system meterdefs should get created")

			foundMeterBase := &marketplacev1alpha1.MeterBase{}
			Expect(k8sClient.Get(context.TODO(), meterBaseKey, foundMeterBase))

			foundMeterBase.Spec.MeterdefinitionCatalogServerConfig.SyncSystemMeterDefinitions = false
			foundMeterBase.Spec.MeterdefinitionCatalogServerConfig.SyncCommunityMeterDefinitions = false
			Expect(k8sClient.Update(context.TODO(), foundMeterBase)).Should(Succeed(), "disable SyncSystemMeterDefinitions and SyncCommunityMeterDefinitions feature flags")

			Eventually(func() []marketplacev1beta1.MeterDefinition {
				labelsMap := map[string]string{}
				err := json.Unmarshal(communityMeterDefIndexLabelsBody, &labelsMap)
				if err != nil {
					log.Fatal(err)
				}

				listOts := []client.ListOption{
					client.MatchingLabels(labelsMap),
				}

				mdefList := &marketplacev1beta1.MeterDefinitionList{}
				k8sClient.List(context.TODO(), mdefList, listOts...)

				return mdefList.Items
			}, timeout, interval).Should(HaveLen(0), "community meterdefs should be deleted")

			Eventually(func() []marketplacev1beta1.MeterDefinition {
				labelsMap := map[string]string{}
				err := json.Unmarshal(systemMeterDefIndexLabelsBody, &labelsMap)
				if err != nil {
					log.Fatal(err)
				}

				listOts := []client.ListOption{
					client.MatchingLabels(labelsMap),
				}

				mdefList := &marketplacev1beta1.MeterDefinitionList{}
				k8sClient.List(context.TODO(), mdefList, listOts...)

				return mdefList.Items
			}, timeout, interval).Should(HaveLen(0), "system meterdefs should be deleted")
		})
	})

	Context("DeployMeterDefinitionCatalogServer feature flag", func() {
		BeforeEach(func() {
			_subSectionMeterBase := subSectionMeterBase.DeepCopy()
			_subSectionMeterBase.Spec.MeterdefinitionCatalogServerConfig.DeployMeterDefinitionCatalogServer = false
			Expect(k8sClient.Create(context.TODO(), _subSectionMeterBase)).Should(Succeed(), "create sub-section meterbase")
		})

		It("all file server resources should be deleted if DeployMeterDefinitionCatalogServer is disabled", func() {
			Eventually(func() bool {
				var dcNotFound bool
				var isNotFound bool
				var serviceIsNotFound bool

				dc := &osappsv1.DeploymentConfig{}
				err := k8sClient.Get(context.TODO(), types.NamespacedName{Name: utils.DeploymentConfigName, Namespace: namespace}, dc)
				if k8serrors.IsNotFound(err) {
					dcNotFound = true
				}

				is := &osimagev1.ImageStreamImage{}
				err = k8sClient.Get(context.TODO(), types.NamespacedName{Name: utils.DeploymentConfigName, Namespace: namespace}, is)
				if k8serrors.IsNotFound(err) {
					isNotFound = true
				}

				service := &corev1.Service{}
				err = k8sClient.Get(context.TODO(), types.NamespacedName{Name: utils.DeploymentConfigName, Namespace: namespace}, service)
				if k8serrors.IsNotFound(err) {
					serviceIsNotFound = true
				}

				return dcNotFound && isNotFound && serviceIsNotFound
			}, timeout, interval).Should(BeTrue())
		})
	})
})
