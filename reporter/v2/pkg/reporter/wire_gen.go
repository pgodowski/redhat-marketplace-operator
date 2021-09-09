// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package reporter

import (
	"context"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/managers"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/prometheus"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/reconcileutils"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Injectors from wire.go:

func NewTask(ctx context.Context, reportName ReportName, taskConfig *Config) (*Task, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	restMapper, err := managers.NewDynamicRESTMapper(restConfig)
	if err != nil {
		return nil, err
	}
	scheme := provideScheme()
	simpleClient, err := managers.ProvideSimpleClient(restConfig, restMapper, scheme)
	if err != nil {
		return nil, err
	}
	logrLogger := _wireLoggerValue
	clientCommandRunner := reconcileutils.NewClientCommand(simpleClient, scheme, logrLogger)
	uploader, err := ProvideUploader(ctx, clientCommandRunner, logrLogger, taskConfig)
	if err != nil {
		return nil, err
	}
	task := &Task{
		ReportName: reportName,
		CC:         clientCommandRunner,
		K8SClient:  simpleClient,
		Ctx:        ctx,
		Config:     taskConfig,
		K8SScheme:  scheme,
		Uploader:   uploader,
	}
	return task, nil
}

var (
	_wireLoggerValue = logger
)

func NewReporter(task *Task) (*MarketplaceReporter, error) {
	reporterConfig := task.Config
	contextContext := task.Ctx
	simpleClient := task.K8SClient
	scheme := task.K8SScheme
	logrLogger := _wireLogrLoggerValue
	clientCommandRunner := reconcileutils.NewClientCommand(simpleClient, scheme, logrLogger)
	reportName := task.ReportName
	meterReport, err := getMarketplaceReport(contextContext, clientCommandRunner, reportName)
	if err != nil {
		return nil, err
	}
	marketplaceConfig, err := getMarketplaceConfig(contextContext, clientCommandRunner)
	if err != nil {
		return nil, err
	}
	service, err := getPrometheusService(contextContext, meterReport, clientCommandRunner)
	if err != nil {
		return nil, err
	}
	prometheusAPISetup := providePrometheusSetup(reporterConfig, meterReport, service)
	prometheusAPI, err := prometheus.NewPrometheusAPIForReporter(prometheusAPISetup)
	if err != nil {
		return nil, err
	}
	v, err := getMeterDefinitionReferences(contextContext, meterReport, clientCommandRunner)
	if err != nil {
		return nil, err
	}
	marketplaceReporter, err := NewMarketplaceReporter(reporterConfig, meterReport, marketplaceConfig, prometheusAPI, v)
	if err != nil {
		return nil, err
	}
	return marketplaceReporter, nil
}

var (
	_wireLogrLoggerValue = logger
)

func NewReporterV2(task *Task) (*MarketplaceReporterV2, error) {
	reporterConfig := task.Config
	contextContext := task.Ctx
	simpleClient := task.K8SClient
	scheme := task.K8SScheme
	logrLogger := _wireLoggerValue2
	clientCommandRunner := reconcileutils.NewClientCommand(simpleClient, scheme, logrLogger)
	reportName := task.ReportName
	meterReport, err := getMarketplaceReport(contextContext, clientCommandRunner, reportName)
	if err != nil {
		return nil, err
	}
	marketplaceConfig, err := getMarketplaceConfig(contextContext, clientCommandRunner)
	if err != nil {
		return nil, err
	}
	service, err := getPrometheusService(contextContext, meterReport, clientCommandRunner)
	if err != nil {
		return nil, err
	}
	prometheusAPISetup := providePrometheusSetup(reporterConfig, meterReport, service)
	prometheusAPI, err := prometheus.NewPrometheusAPIForReporter(prometheusAPISetup)
	if err != nil {
		return nil, err
	}
	v, err := getMeterDefinitionReferences(contextContext, meterReport, clientCommandRunner)
	if err != nil {
		return nil, err
	}
	marketplaceReporterV2, err := NewMarketplaceReporterV2(reporterConfig, meterReport, marketplaceConfig, prometheusAPI, v)
	if err != nil {
		return nil, err
	}
	return marketplaceReporterV2, nil
}

var (
	_wireLoggerValue2 = logger
)

func NewUploadTask(ctx context.Context, config2 *Config) (*UploadTask, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	restMapper, err := managers.NewDynamicRESTMapper(restConfig)
	if err != nil {
		return nil, err
	}
	scheme := provideScheme()
	simpleClient, err := managers.ProvideSimpleClient(restConfig, restMapper, scheme)
	if err != nil {
		return nil, err
	}
	logrLogger := _wireLoggerValue3
	clientCommandRunner := reconcileutils.NewClientCommand(simpleClient, scheme, logrLogger)
	downloader, err := ProvideDownloader(ctx, clientCommandRunner, logrLogger, config2)
	if err != nil {
		return nil, err
	}
	uploader, err := ProvideUploader(ctx, clientCommandRunner, logrLogger, config2)
	if err != nil {
		return nil, err
	}
	admin, err := ProvideAdmin(ctx, clientCommandRunner, logrLogger, config2)
	if err != nil {
		return nil, err
	}
	uploadTask := &UploadTask{
		CC:         clientCommandRunner,
		K8SClient:  simpleClient,
		Ctx:        ctx,
		Config:     config2,
		K8SScheme:  scheme,
		Downloader: downloader,
		Uploader:   uploader,
		Admin:      admin,
	}
	return uploadTask, nil
}

var (
	_wireLoggerValue3 = logger
)
