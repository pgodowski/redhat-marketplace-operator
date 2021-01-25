// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/redhat-marketplace/redhat-marketplace-operator/authchecker/v2/pkg/authchecker"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/client"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/managers"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Injectors from wire.go:

func InitializeAuthChecker(cfg authchecker.AuthCheckerConfig) (*authchecker.AuthChecker, error) {
	logger := _wireLoggerValue
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	dynamicInterface, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	restMapper, err := managers.NewDynamicRESTMapper(restConfig)
	if err != nil {
		return nil, err
	}
	dynamicClient := client.NewDynamicClient(dynamicInterface, restMapper)
	authChecker, err := authchecker.NewAuthChecker(logger, dynamicClient, cfg)
	if err != nil {
		return nil, err
	}
	return authChecker, nil
}

var (
	_wireLoggerValue = log
)