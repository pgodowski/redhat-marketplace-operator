/*
Copyright 2023 IBM Co..

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	marketplacev1alpha1 "github.com/redhat-marketplace/redhat-marketplace-operator/datareporter/v2/api/v1alpha1"
)

// DataReporterConfigReconciler reconciles a DataReporterConfig object
type DataReporterConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// data-service
//+kubebuilder:rbac:urls=/dataservice.v1.fileserver.FileServer/*,verbs=create
// kube-rbac-proxy
//+kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
//+kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create

// controller CRDs
//+kubebuilder:rbac:groups=marketplace.redhat.com,resources=datareporterconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=marketplace.redhat.com,resources=datareporterconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=marketplace.redhat.com,resources=datareporterconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DataReporterConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *DataReporterConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DataReporterConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&marketplacev1alpha1.DataReporterConfig{}).
		Complete(r)
}
