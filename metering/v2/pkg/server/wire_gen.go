// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package server

import (
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/internal/metrics"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/engine"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/managers"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/reconcileutils"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Injectors from wire.go:

func NewServer(opts *Options) (*Service, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	restMapper, err := managers.NewDynamicRESTMapper(restConfig)
	if err != nil {
		return nil, err
	}
	scheme := provideScheme()
	clientOptions := getClientOptions()
	cache, err := managers.ProvideNewCache(restConfig, restMapper, scheme, clientOptions)
	if err != nil {
		return nil, err
	}
	client, err := managers.ProvideCachedClient(restConfig, restMapper, scheme, cache, clientOptions)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	options := ConvertOptions(opts)
	registry := provideRegistry()
	logger := _wireLoggerValue
	clientCommandRunner := reconcileutils.NewClientCommand(client, scheme, logger)
	context := provideContext()
	cacheIsIndexed, err := managers.AddIndices(context, cache)
	if err != nil {
		return nil, err
	}
	cacheIsStarted, err := managers.StartCache(context, cache, logger, cacheIsIndexed)
	if err != nil {
		return nil, err
	}
	namespaces := ProvideNamespaces(opts)
	prometheusData := metrics.ProvidePrometheusData()
	engineEngine, err := engine.NewEngine(context, namespaces, scheme, clientOptions, restConfig, logger, prometheusData)
	if err != nil {
		return nil, err
	}
	service := &Service{
		k8sclient:       client,
		k8sRestClient:   clientset,
		opts:            options,
		cache:           cache,
		metricsRegistry: registry,
		cc:              clientCommandRunner,
		indexed:         cacheIsIndexed,
		started:         cacheIsStarted,
		engine:          engineEngine,
		prometheusData:  prometheusData,
	}
	return service, nil
}

var (
	_wireLoggerValue = log
)
