// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package engine

import (
	"context"
	"github.com/go-logr/logr"
	v1_2 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	"github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/internal/metrics"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/dictionary"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/mailbox"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/meterdefinition"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/processors"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/processorsenders"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/razee"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/types"
	v1alpha1_2 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/generated/clientset/versioned/typed/marketplace/v1alpha1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/generated/clientset/versioned/typed/marketplace/v1beta1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/client"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/managers"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Injectors from wire.go:

func NewEngine(ctx context.Context, namespaces types.Namespaces, scheme *runtime.Scheme, clientOptions managers.ClientOptions, k8sRestConfig *rest.Config, log logr.Logger, prometheusData *metrics.PrometheusData) (*Engine, error) {
	clientset, err := kubernetes.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	dynamicInterface, err := dynamic.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	restMapper, err := managers.NewDynamicRESTMapper(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	dynamicClient := client.NewDynamicClient(dynamicInterface, restMapper)
	findOwnerHelper := client.NewFindOwnerHelper(dynamicClient)
	monitoringV1Client, err := v1.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	marketplaceV1beta1Client, err := v1beta1.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	cache, err := managers.ProvideNewCache(k8sRestConfig, restMapper, scheme, clientOptions)
	if err != nil {
		return nil, err
	}
	clientClient, err := managers.ProvideCachedClient(k8sRestConfig, restMapper, scheme, cache, clientOptions)
	if err != nil {
		return nil, err
	}
	cacheIsIndexed, err := managers.AddIndices(ctx, cache)
	if err != nil {
		return nil, err
	}
	cacheIsStarted, err := managers.StartCache(ctx, cache, log, cacheIsIndexed)
	if err != nil {
		return nil, err
	}
	meterDefinitionList, err := dictionary.ProvideMeterDefinitionList(cacheIsStarted, clientClient)
	if err != nil {
		return nil, err
	}
	meterDefinitionsSeenStore := dictionary.NewMeterDefinitionsSeenStore()
	meterDefinitionDictionary := dictionary.NewMeterDefinitionDictionary(ctx, clientset, clientClient, findOwnerHelper, namespaces, log, meterDefinitionList, meterDefinitionsSeenStore)
	meterDefinitionStore := meterdefinition.NewMeterDefinitionStore(ctx, log, clientset, findOwnerHelper, monitoringV1Client, marketplaceV1beta1Client, meterDefinitionDictionary, scheme)
	meterDefinitionStoreRunnable := ProvideMeterDefinitionStoreRunnable(clientset, namespaces, meterDefinitionStore, monitoringV1Client, log)
	meterDefinitionDictionaryStoreRunnable := ProvideMeterDefinitionDictionaryStoreRunnable(clientset, namespaces, marketplaceV1beta1Client, meterDefinitionDictionary, log)
	meterDefinitionSeenStoreRunnable := ProvideMeterDefinitionSeenStoreRunnable(clientset, namespaces, marketplaceV1beta1Client, meterDefinitionsSeenStore, log)
	mailboxMailbox := mailbox.ProvideMailbox(log)
	statusProcessor := processors.ProvideStatusProcessor(log, clientClient, mailboxMailbox, scheme)
	serviceAnnotatorProcessor := processors.ProvideServiceAnnotatorProcessor(log, clientClient, mailboxMailbox)
	prometheusProcessor := processors.ProvidePrometheusProcessor(log, clientClient, mailboxMailbox, scheme, prometheusData)
	prometheusMdefProcessor := processors.ProvidePrometheusMdefProcessor(log, clientClient, mailboxMailbox, scheme, prometheusData)
	meterDefinitionRemovalWatcher := processors.ProvideMeterDefinitionRemovalWatcher(meterDefinitionDictionary, meterDefinitionStore, mailboxMailbox, log)
	objectChannelProducer := mailbox.ProvideObjectChannelProducer(meterDefinitionStore, mailboxMailbox, log)
	meterDefinitionChannelProducer := mailbox.ProvideMeterDefinitionChannelProducer(meterDefinitionDictionary, mailboxMailbox, log)
	runnables := ProvideRunnables(meterDefinitionStoreRunnable, meterDefinitionDictionaryStoreRunnable, meterDefinitionSeenStoreRunnable, mailboxMailbox, statusProcessor, serviceAnnotatorProcessor, prometheusProcessor, prometheusMdefProcessor, meterDefinitionRemovalWatcher, objectChannelProducer, meterDefinitionChannelProducer, meterDefinitionDictionary)
	engine := ProvideEngine(meterDefinitionStore, namespaces, log, clientset, monitoringV1Client, meterDefinitionDictionary, marketplaceV1beta1Client, runnables, prometheusData)
	return engine, nil
}

func NewRazeeEngine(ctx context.Context, namespaces types.Namespaces, scheme *runtime.Scheme, clientOptions managers.ClientOptions, k8sRestConfig *rest.Config, log logr.Logger) (*RazeeEngine, error) {
	restMapper, err := managers.NewDynamicRESTMapper(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	cache, err := managers.ProvideNewCache(k8sRestConfig, restMapper, scheme, clientOptions)
	if err != nil {
		return nil, err
	}
	cacheIsIndexed := managers.CacheIsIndexed{}
	cacheIsStarted, err := managers.StartCache(ctx, cache, log, cacheIsIndexed)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	dynamicInterface, err := dynamic.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	dynamicClient := client.NewDynamicClient(dynamicInterface, restMapper)
	findOwnerHelper := client.NewFindOwnerHelper(dynamicClient)
	razeeStoreGroup := razee.NewRazeeStore(ctx, log, clientset, findOwnerHelper, scheme)
	razeeStore := razeeStoreGroup.Store
	operatorsV1alpha1Client, err := v1alpha1.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	configV1Client, err := v1_2.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	marketplaceV1alpha1Client, err := v1alpha1_2.NewForConfig(k8sRestConfig)
	if err != nil {
		return nil, err
	}
	razeeStores := razeeStoreGroup.Stores
	razeeStoreRunnable := ProvideRazeeStoreRunnable(clientset, operatorsV1alpha1Client, configV1Client, marketplaceV1alpha1Client, namespaces, razeeStores, log)
	mailboxMailbox := mailbox.ProvideMailbox(log)
	clientClient, err := managers.ProvideCachedClient(k8sRestConfig, restMapper, scheme, cache, clientOptions)
	if err != nil {
		return nil, err
	}
	razeeProcessorSender := processorsenders.ProvideRazeeProcessorSender(log, clientClient, mailboxMailbox, scheme)
	razeeChannelProducer := mailbox.ProvideRazeeChannelProducer(razeeStore, mailboxMailbox, log)
	runnables := ProvideRazeeRunnables(razeeStoreRunnable, mailboxMailbox, razeeProcessorSender, razeeChannelProducer)
	razeeEngine := ProvideRazeeEngine(cacheIsStarted, razeeStore, namespaces, log, clientset, runnables)
	return razeeEngine, nil
}
