// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package engine

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/internal/metrics"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/filter"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/mailbox"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/processors"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/stores"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/types"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/generated/clientset/versioned/typed/marketplace/v1beta1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/client"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/managers"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/metadata"
	"k8s.io/client-go/rest"
)

// Injectors from wire.go:

func NewEngine(ctx context.Context, namespaces types.Namespaces, scheme *runtime.Scheme, restConfig *rest.Config, clientOptions managers.ClientOptions, log logr.Logger, prometheusData *metrics.PrometheusData, statusFlushDuration processors.StatusFlushDuration) (*Engine, error) {
	namespaceWatcher := filter.ProvideNamespaceWatcher(log)
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	restMapper, err := managers.NewDynamicRESTMapper(restConfig)
	if err != nil {
		return nil, err
	}
	cache, err := managers.ProvideNewCache(restConfig, restMapper, scheme, clientOptions)
	if err != nil {
		return nil, err
	}
	clientClient, err := managers.ProvideCachedClient(restConfig, restMapper, scheme, cache, clientOptions)
	if err != nil {
		return nil, err
	}
	metadataInterface, err := metadata.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	metadataClient := client.NewMetadataClient(metadataInterface, restMapper)
	findOwnerHelper := client.NewFindOwnerHelper(ctx, metadataClient)
	cacheIsIndexed, err := managers.AddIndices(ctx, cache)
	if err != nil {
		return nil, err
	}
	cacheIsStarted, err := managers.StartCache(ctx, cache, log, cacheIsIndexed)
	if err != nil {
		return nil, err
	}
	meterDefinitionLookupFilterFactory := &filter.MeterDefinitionLookupFilterFactory{
		Client:    clientClient,
		Log:       log,
		FindOwner: findOwnerHelper,
		IsStarted: cacheIsStarted,
	}
	meterDefinitionDictionary := stores.NewMeterDefinitionDictionary(ctx, namespaces, log, meterDefinitionLookupFilterFactory)
	meterDefinitionStore := stores.NewMeterDefinitionStore(ctx, log, meterDefinitionDictionary, scheme)
	monitoringV1Client, err := v1.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	meterDefinitionStoreListWatchers := ProvideMeterDefinitionStoreListWatchers(clientset, meterDefinitionStore, monitoringV1Client)
	meterDefinitionStoreRunnable := ProvideMeterDefinitionStoreRunnable(namespaces, log, namespaceWatcher, meterDefinitionStoreListWatchers, meterDefinitionStore)
	marketplaceV1beta1Client, err := v1beta1.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	meterDefinitionDictionaryStoreRunnable := ProvideMeterDefinitionDictionaryStoreRunnable(clientset, namespaces, marketplaceV1beta1Client, meterDefinitionDictionary, log)
	mailboxMailbox := mailbox.ProvideMailbox(log)
	statusProcessor := processors.ProvideStatusProcessor(log, clientClient, mailboxMailbox, statusFlushDuration, cacheIsStarted)
	serviceAnnotatorProcessor := processors.ProvideServiceAnnotatorProcessor(log, clientClient, mailboxMailbox, cacheIsStarted)
	prometheusProcessor := processors.ProvidePrometheusProcessor(log, mailboxMailbox, scheme, prometheusData)
	prometheusMdefProcessor := processors.ProvidePrometheusMdefProcessor(log, mailboxMailbox, scheme, prometheusData, cacheIsStarted)
	meterDefinitionRemovalWatcher := processors.ProvideMeterDefinitionRemovalWatcher(meterDefinitionDictionary, meterDefinitionStore, mailboxMailbox, log, clientClient, namespaceWatcher)
	objectChannelProducer := mailbox.ProvideObjectChannelProducer(meterDefinitionStore, mailboxMailbox, log)
	meterDefinitionChannelProducer := mailbox.ProvideMeterDefinitionChannelProducer(meterDefinitionDictionary, mailboxMailbox, log)
	runnables := ProvideRunnables(meterDefinitionStoreRunnable, meterDefinitionDictionaryStoreRunnable, mailboxMailbox, statusProcessor, serviceAnnotatorProcessor, prometheusProcessor, prometheusMdefProcessor, meterDefinitionRemovalWatcher, objectChannelProducer, meterDefinitionChannelProducer, namespaceWatcher)
	engine := ProvideEngine(log, runnables)
	return engine, nil
}
