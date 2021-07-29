// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package inject

import (
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/config"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/managers"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/manifests"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/runnables"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/reconcileutils"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

// Injectors from wire.go:

func initializeInjectDependencies(cache2 cache.Cache, fields *managers.ControllerFields) (injectorDependencies, error) {
	logger := fields.Logger
	client := fields.Client
	scheme := fields.Scheme
	clientCommandRunner := reconcileutils.NewClientCommand(client, scheme, logger)
	restConfig := fields.Config
	restMapper, err := managers.NewDynamicRESTMapper(restConfig)
	if err != nil {
		return injectorDependencies{}, err
	}
	simpleClient, err := managers.ProvideSimpleClient(restConfig, restMapper, scheme)
	if err != nil {
		return injectorDependencies{}, err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return injectorDependencies{}, err
	}
	operatorConfig, err := config.ProvideInfrastructureAwareConfig(simpleClient, discoveryClient)
	if err != nil {
		return injectorDependencies{}, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return injectorDependencies{}, err
	}
	factory := manifests.NewFactory(operatorConfig, scheme)
	crdUpdater := &runnables.CRDUpdater{
		Logger:  logger,
		CC:      clientCommandRunner,
		Config:  operatorConfig,
		Rest:    restConfig,
		Client:  clientset,
		Factory: factory,
	}
	runnablesRunnables := runnables.ProvideRunnables(crdUpdater)
	clientCommandInjector := &ClientCommandInjector{
		Fields:        fields,
		CommandRunner: clientCommandRunner,
	}
	operatorConfigInjector := &OperatorConfigInjector{
		Config: operatorConfig,
	}
	patchInjector := &PatchInjector{}
	deployedNamespace := ProvideNamespace(operatorConfig)
	factoryInjector := &FactoryInjector{
		Fields:    fields,
		Config:    operatorConfig,
		Namespace: deployedNamespace,
		Scheme:    scheme,
		Factory:   factory,
	}
	kubeInterfaceInjector := &KubeInterfaceInjector{
		KubeInterface: clientset,
	}
	injectables := ProvideInjectables(clientCommandInjector, operatorConfigInjector, patchInjector, factoryInjector, kubeInterfaceInjector)
	injectInjectorDependencies := injectorDependencies{
		Runnables:   runnablesRunnables,
		Injectables: injectables,
	}
	return injectInjectorDependencies, nil
}

func initializeCommandRunner(fields *managers.ControllerFields) (reconcileutils.ClientCommandRunner, error) {
	client := fields.Client
	scheme := fields.Scheme
	logger := fields.Logger
	clientCommandRunner := reconcileutils.NewClientCommand(client, scheme, logger)
	return clientCommandRunner, nil
}
