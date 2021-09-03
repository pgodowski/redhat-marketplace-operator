/*
Copyright 2020 IBM Co..

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package marketplace

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	openshiftconfigv1 "github.com/openshift/api/config/v1"
	olmv1 "github.com/operator-framework/api/pkg/operators/v1"
	opsrcv1 "github.com/operator-framework/api/pkg/operators/v1"
	olmv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/onsi/gomega/types"
	osappsv1 "github.com/openshift/api/apps/v1"
	osimagev1 "github.com/openshift/api/image/v1"
	marketplaceredhatcomv1alpha1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1alpha1"
	marketplaceredhatcomv1beta1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/catalog"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/config"
	"k8s.io/apimachinery/pkg/api/errors"

	// "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/inject"loo
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/manifests"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var k8sManager ctrl.Manager
var k8sScheme *runtime.Scheme

const (
	imageStreamID string = "rhm-meterdefinition-file-server:v1"
	imageStreamTag string = "v1"
	listenerAddress string = "127.0.0.1:2100"
)


func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))
	os.Setenv("KUBEBUILDER_CONTROLPLANE_START_TIMEOUT", "2m")
	os.Setenv("POD_NAMESPACE","default")
	os.Setenv("IMAGE_STREAM_ID",imageStreamID)
	os.Setenv("IMAGE_STREAM_TAG",imageStreamTag)

	dcControllerMockServerAddr := fmt.Sprintf("%s%s","http://",listenerAddress)
	os.Setenv("CATALOG_URL", dcControllerMockServerAddr)

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:  []string{filepath.Join("..", "..", "config", "crd", "bases")},
		KubeAPIServerFlags: append(envtest.DefaultKubeAPIServerFlags, "--bind-address=127.0.0.1"),
	}

	var err error
	k8sScheme = provideScheme()
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = marketplaceredhatcomv1alpha1.AddToScheme(k8sScheme)
	Expect(err).NotTo(HaveOccurred())

	err = marketplaceredhatcomv1beta1.AddToScheme(k8sScheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme
	k8sClient, err = client.New(cfg, client.Options{Scheme: k8sScheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: k8sScheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&RemoteResourceS3Reconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("RemoteResourceS3"),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	operatorConfig, err := config.GetConfig()
	Expect(err).NotTo(HaveOccurred())

	factory := manifests.NewFactory(operatorConfig,k8sScheme)
	
	restConfig := k8sManager.GetConfig()

	clientset, err := kubernetes.NewForConfig(restConfig)
	Expect(err).NotTo(HaveOccurred())

	catalogClient,err := catalog.ProvideCatalogClient(k8sManager.GetClient(),operatorConfig,clientset)
	Expect(err).NotTo(HaveOccurred())

	catalogClient.UseInsecureClient()

	err = (&DeploymentConfigReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("DeploymentConfigReconciler"),
		Scheme: k8sManager.GetScheme(),
		cfg: operatorConfig,
		factory: factory,
		CatalogClient: catalogClient,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		fmt.Println(err)
		Expect(err).ToNot(HaveOccurred())
	}()
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	// dcControllerMockServer.Close()
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

func provideScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(openshiftconfigv1.AddToScheme(scheme))
	utilruntime.Must(olmv1.AddToScheme(scheme))
	utilruntime.Must(opsrcv1.AddToScheme(scheme))
	utilruntime.Must(olmv1alpha1.AddToScheme(scheme))
	utilruntime.Must(monitoringv1.AddToScheme(scheme))
	utilruntime.Must(marketplaceredhatcomv1alpha1.AddToScheme(scheme))
	utilruntime.Must(marketplaceredhatcomv1beta1.AddToScheme(scheme))
	utilruntime.Must(osappsv1.AddToScheme(scheme))
	utilruntime.Must(osimagev1.AddToScheme(scheme))
	return scheme
}

var SucceedOrAlreadyExist types.GomegaMatcher = SatisfyAny(
	Succeed(),
	WithTransform(errors.IsAlreadyExists, BeTrue()),
)

var IsNotFound types.GomegaMatcher = WithTransform(errors.IsNotFound, BeTrue())

