// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.ibm.com/symposium/redhat-marketplace-operator/pkg/controller"
	"github.ibm.com/symposium/redhat-marketplace-operator/pkg/managers"
)

// Injectors from wire.go:

func initializeMarketplaceController() *managers.ControllerMain {
	marketplaceController := controller.ProvideMarketplaceController()
	meterbaseController := controller.ProvideMeterbaseController()
	meterDefinitionController := controller.ProvideMeterDefinitionController()
	razeeDeployController := controller.ProvideRazeeDeployController()
	controllerFlagSet := controller.ProvideControllerFlagSet()
	opsSrcSchemeDefinition := managers.ProvideOpsSrcScheme()
	monitoringSchemeDefinition := managers.ProvideMonitoringScheme()
	controllerMain := makeMarketplaceController(marketplaceController, meterbaseController, meterDefinitionController, razeeDeployController, controllerFlagSet, opsSrcSchemeDefinition, monitoringSchemeDefinition)
	return controllerMain
}