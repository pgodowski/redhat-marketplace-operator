// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/common"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/status"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIKeyRef) DeepCopyInto(out *APIKeyRef) {
	*out = *in
	in.ValueFrom.DeepCopyInto(&out.ValueFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIKeyRef.
func (in *APIKeyRef) DeepCopy() *APIKeyRef {
	if in == nil {
		return nil
	}
	out := new(APIKeyRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccesKeyIDRef) DeepCopyInto(out *AccesKeyIDRef) {
	*out = *in
	in.ValueFrom.DeepCopyInto(&out.ValueFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccesKeyIDRef.
func (in *AccesKeyIDRef) DeepCopy() *AccesKeyIDRef {
	if in == nil {
		return nil
	}
	out := new(AccesKeyIDRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Auth) DeepCopyInto(out *Auth) {
	*out = *in
	if in.Hmac != nil {
		in, out := &in.Hmac, &out.Hmac
		*out = new(Hmac)
		(*in).DeepCopyInto(*out)
	}
	if in.Iam != nil {
		in, out := &in.Iam, &out.Iam
		*out = new(Iam)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Auth.
func (in *Auth) DeepCopy() *Auth {
	if in == nil {
		return nil
	}
	out := new(Auth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CSVNamespacedName) DeepCopyInto(out *CSVNamespacedName) {
	*out = *in
	if in.GroupVersionKind != nil {
		in, out := &in.GroupVersionKind, &out.GroupVersionKind
		*out = new(common.GroupVersionKind)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CSVNamespacedName.
func (in *CSVNamespacedName) DeepCopy() *CSVNamespacedName {
	if in == nil {
		return nil
	}
	out := new(CSVNamespacedName)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ErrorDetails) DeepCopyInto(out *ErrorDetails) {
	*out = *in
	if in.Details != nil {
		in, out := &in.Details, &out.Details
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ErrorDetails.
func (in *ErrorDetails) DeepCopy() *ErrorDetails {
	if in == nil {
		return nil
	}
	out := new(ErrorDetails)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Header) DeepCopyInto(out *Header) {
	{
		in := &in
		*out = make(Header, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Header.
func (in Header) DeepCopy() Header {
	if in == nil {
		return nil
	}
	out := new(Header)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Hmac) DeepCopyInto(out *Hmac) {
	*out = *in
	in.AccesKeyIDRef.DeepCopyInto(&out.AccesKeyIDRef)
	in.SecretAccessKeyRef.DeepCopyInto(&out.SecretAccessKeyRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Hmac.
func (in *Hmac) DeepCopy() *Hmac {
	if in == nil {
		return nil
	}
	out := new(Hmac)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Iam) DeepCopyInto(out *Iam) {
	*out = *in
	in.APIKeyRef.DeepCopyInto(&out.APIKeyRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Iam.
func (in *Iam) DeepCopy() *Iam {
	if in == nil {
		return nil
	}
	out := new(Iam)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Log) DeepCopyInto(out *Log) {
	{
		in := &in
		*out = make(Log, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Log.
func (in Log) DeepCopy() Log {
	if in == nil {
		return nil
	}
	out := new(Log)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MarketplaceConfig) DeepCopyInto(out *MarketplaceConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MarketplaceConfig.
func (in *MarketplaceConfig) DeepCopy() *MarketplaceConfig {
	if in == nil {
		return nil
	}
	out := new(MarketplaceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MarketplaceConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MarketplaceConfigList) DeepCopyInto(out *MarketplaceConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MarketplaceConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MarketplaceConfigList.
func (in *MarketplaceConfigList) DeepCopy() *MarketplaceConfigList {
	if in == nil {
		return nil
	}
	out := new(MarketplaceConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MarketplaceConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MarketplaceConfigSpec) DeepCopyInto(out *MarketplaceConfigSpec) {
	*out = *in
	if in.DeploySecretName != nil {
		in, out := &in.DeploySecretName, &out.DeploySecretName
		*out = new(string)
		**out = **in
	}
	if in.EnableMetering != nil {
		in, out := &in.EnableMetering, &out.EnableMetering
		*out = new(bool)
		**out = **in
	}
	if in.InstallIBMCatalogSource != nil {
		in, out := &in.InstallIBMCatalogSource, &out.InstallIBMCatalogSource
		*out = new(bool)
		**out = **in
	}
	if in.Features != nil {
		in, out := &in.Features, &out.Features
		*out = new(common.Features)
		(*in).DeepCopyInto(*out)
	}
	if in.NamespaceLabelSelector != nil {
		in, out := &in.NamespaceLabelSelector, &out.NamespaceLabelSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MarketplaceConfigSpec.
func (in *MarketplaceConfigSpec) DeepCopy() *MarketplaceConfigSpec {
	if in == nil {
		return nil
	}
	out := new(MarketplaceConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MarketplaceConfigStatus) DeepCopyInto(out *MarketplaceConfigStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(status.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.RazeeSubConditions != nil {
		in, out := &in.RazeeSubConditions, &out.RazeeSubConditions
		*out = make(status.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.MeterBaseSubConditions != nil {
		in, out := &in.MeterBaseSubConditions, &out.MeterBaseSubConditions
		*out = make(status.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MarketplaceConfigStatus.
func (in *MarketplaceConfigStatus) DeepCopy() *MarketplaceConfigStatus {
	if in == nil {
		return nil
	}
	out := new(MarketplaceConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterBase) DeepCopyInto(out *MeterBase) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterBase.
func (in *MeterBase) DeepCopy() *MeterBase {
	if in == nil {
		return nil
	}
	out := new(MeterBase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MeterBase) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterBaseList) DeepCopyInto(out *MeterBaseList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MeterBase, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterBaseList.
func (in *MeterBaseList) DeepCopy() *MeterBaseList {
	if in == nil {
		return nil
	}
	out := new(MeterBaseList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MeterBaseList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterBaseSpec) DeepCopyInto(out *MeterBaseSpec) {
	*out = *in
	if in.Prometheus != nil {
		in, out := &in.Prometheus, &out.Prometheus
		*out = new(PrometheusSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.AdditionalScrapeConfigs != nil {
		in, out := &in.AdditionalScrapeConfigs, &out.AdditionalScrapeConfigs
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterBaseSpec.
func (in *MeterBaseSpec) DeepCopy() *MeterBaseSpec {
	if in == nil {
		return nil
	}
	out := new(MeterBaseSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterBaseStatus) DeepCopyInto(out *MeterBaseStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(status.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.PrometheusStatus != nil {
		in, out := &in.PrometheusStatus, &out.PrometheusStatus
		*out = new(monitoringv1.PrometheusStatus)
		**out = **in
	}
	if in.MeterdefinitionCatalogServerStatus != nil {
		in, out := &in.MeterdefinitionCatalogServerStatus, &out.MeterdefinitionCatalogServerStatus
		*out = new(MeterdefinitionCatalogServerStatus)
		(*in).DeepCopyInto(*out)
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.UpdatedReplicas != nil {
		in, out := &in.UpdatedReplicas, &out.UpdatedReplicas
		*out = new(int32)
		**out = **in
	}
	if in.AvailableReplicas != nil {
		in, out := &in.AvailableReplicas, &out.AvailableReplicas
		*out = new(int32)
		**out = **in
	}
	if in.UnavailableReplicas != nil {
		in, out := &in.UnavailableReplicas, &out.UnavailableReplicas
		*out = new(int32)
		**out = **in
	}
	if in.Targets != nil {
		in, out := &in.Targets, &out.Targets
		*out = make([]common.Target, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterBaseStatus.
func (in *MeterBaseStatus) DeepCopy() *MeterBaseStatus {
	if in == nil {
		return nil
	}
	out := new(MeterBaseStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterDefinition) DeepCopyInto(out *MeterDefinition) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterDefinition.
func (in *MeterDefinition) DeepCopy() *MeterDefinition {
	if in == nil {
		return nil
	}
	out := new(MeterDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MeterDefinition) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterDefinitionList) DeepCopyInto(out *MeterDefinitionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MeterDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterDefinitionList.
func (in *MeterDefinitionList) DeepCopy() *MeterDefinitionList {
	if in == nil {
		return nil
	}
	out := new(MeterDefinitionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MeterDefinitionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterDefinitionSpec) DeepCopyInto(out *MeterDefinitionSpec) {
	*out = *in
	if in.InstalledBy != nil {
		in, out := &in.InstalledBy, &out.InstalledBy
		*out = new(common.NamespacedNameReference)
		(*in).DeepCopyInto(*out)
	}
	if in.VertexLabelSelector != nil {
		in, out := &in.VertexLabelSelector, &out.VertexLabelSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Workloads != nil {
		in, out := &in.Workloads, &out.Workloads
		*out = make([]Workload, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Version != nil {
		in, out := &in.Version, &out.Version
		*out = new(string)
		**out = **in
	}
	if in.ServiceMeterLabels != nil {
		in, out := &in.ServiceMeterLabels, &out.ServiceMeterLabels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.PodMeterLabels != nil {
		in, out := &in.PodMeterLabels, &out.PodMeterLabels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterDefinitionSpec.
func (in *MeterDefinitionSpec) DeepCopy() *MeterDefinitionSpec {
	if in == nil {
		return nil
	}
	out := new(MeterDefinitionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterDefinitionStatus) DeepCopyInto(out *MeterDefinitionStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(status.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.WorkloadResources != nil {
		in, out := &in.WorkloadResources, &out.WorkloadResources
		*out = make([]common.WorkloadResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Results != nil {
		in, out := &in.Results, &out.Results
		*out = make([]common.Result, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterDefinitionStatus.
func (in *MeterDefinitionStatus) DeepCopy() *MeterDefinitionStatus {
	if in == nil {
		return nil
	}
	out := new(MeterDefinitionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterLabelQuery) DeepCopyInto(out *MeterLabelQuery) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterLabelQuery.
func (in *MeterLabelQuery) DeepCopy() *MeterLabelQuery {
	if in == nil {
		return nil
	}
	out := new(MeterLabelQuery)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterReport) DeepCopyInto(out *MeterReport) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterReport.
func (in *MeterReport) DeepCopy() *MeterReport {
	if in == nil {
		return nil
	}
	out := new(MeterReport)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MeterReport) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterReportList) DeepCopyInto(out *MeterReportList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MeterReport, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterReportList.
func (in *MeterReportList) DeepCopy() *MeterReportList {
	if in == nil {
		return nil
	}
	out := new(MeterReportList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MeterReportList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterReportSpec) DeepCopyInto(out *MeterReportSpec) {
	*out = *in
	in.StartTime.DeepCopyInto(&out.StartTime)
	in.EndTime.DeepCopyInto(&out.EndTime)
	if in.PrometheusService != nil {
		in, out := &in.PrometheusService, &out.PrometheusService
		*out = new(common.ServiceReference)
		(*in).DeepCopyInto(*out)
	}
	if in.MeterDefinitions != nil {
		in, out := &in.MeterDefinitions, &out.MeterDefinitions
		*out = make([]MeterDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.MeterDefinitionReferences != nil {
		in, out := &in.MeterDefinitionReferences, &out.MeterDefinitionReferences
		*out = make([]v1beta1.MeterDefinitionReference, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ExtraArgs != nil {
		in, out := &in.ExtraArgs, &out.ExtraArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterReportSpec.
func (in *MeterReportSpec) DeepCopy() *MeterReportSpec {
	if in == nil {
		return nil
	}
	out := new(MeterReportSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterReportStatus) DeepCopyInto(out *MeterReportStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(status.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AssociatedJob != nil {
		in, out := &in.AssociatedJob, &out.AssociatedJob
		*out = new(common.JobReference)
		(*in).DeepCopyInto(*out)
	}
	if in.MetricUploadCount != nil {
		in, out := &in.MetricUploadCount, &out.MetricUploadCount
		*out = new(int)
		**out = **in
	}
	if in.UploadID != nil {
		in, out := &in.UploadID, &out.UploadID
		*out = new(types.UID)
		**out = **in
	}
	if in.Errors != nil {
		in, out := &in.Errors, &out.Errors
		*out = make([]ErrorDetails, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Warnings != nil {
		in, out := &in.Warnings, &out.Warnings
		*out = make([]ErrorDetails, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterReportStatus.
func (in *MeterReportStatus) DeepCopy() *MeterReportStatus {
	if in == nil {
		return nil
	}
	out := new(MeterReportStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeterdefinitionCatalogServerStatus) DeepCopyInto(out *MeterdefinitionCatalogServerStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(status.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterdefinitionCatalogServerStatus.
func (in *MeterdefinitionCatalogServerStatus) DeepCopy() *MeterdefinitionCatalogServerStatus {
	if in == nil {
		return nil
	}
	out := new(MeterdefinitionCatalogServerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusSpec) DeepCopyInto(out *PrometheusSpec) {
	*out = *in
	in.ResourceRequirements.DeepCopyInto(&out.ResourceRequirements)
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Storage.DeepCopyInto(&out.Storage)
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusSpec.
func (in *PrometheusSpec) DeepCopy() *PrometheusSpec {
	if in == nil {
		return nil
	}
	out := new(PrometheusSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RazeeConfigurationValues) DeepCopyInto(out *RazeeConfigurationValues) {
	*out = *in
	if in.IbmCosReaderKey != nil {
		in, out := &in.IbmCosReaderKey, &out.IbmCosReaderKey
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.RazeeDashOrgKey != nil {
		in, out := &in.RazeeDashOrgKey, &out.RazeeDashOrgKey
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.FileSourceURL != nil {
		in, out := &in.FileSourceURL, &out.FileSourceURL
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RazeeConfigurationValues.
func (in *RazeeConfigurationValues) DeepCopy() *RazeeConfigurationValues {
	if in == nil {
		return nil
	}
	out := new(RazeeConfigurationValues)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RazeeDeployment) DeepCopyInto(out *RazeeDeployment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RazeeDeployment.
func (in *RazeeDeployment) DeepCopy() *RazeeDeployment {
	if in == nil {
		return nil
	}
	out := new(RazeeDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RazeeDeployment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RazeeDeploymentList) DeepCopyInto(out *RazeeDeploymentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RazeeDeployment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RazeeDeploymentList.
func (in *RazeeDeploymentList) DeepCopy() *RazeeDeploymentList {
	if in == nil {
		return nil
	}
	out := new(RazeeDeploymentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RazeeDeploymentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RazeeDeploymentSpec) DeepCopyInto(out *RazeeDeploymentSpec) {
	*out = *in
	if in.DeploySecretName != nil {
		in, out := &in.DeploySecretName, &out.DeploySecretName
		*out = new(string)
		**out = **in
	}
	if in.TargetNamespace != nil {
		in, out := &in.TargetNamespace, &out.TargetNamespace
		*out = new(string)
		**out = **in
	}
	if in.DeployConfig != nil {
		in, out := &in.DeployConfig, &out.DeployConfig
		*out = new(RazeeConfigurationValues)
		(*in).DeepCopyInto(*out)
	}
	if in.ChildUrl != nil {
		in, out := &in.ChildUrl, &out.ChildUrl
		*out = new(string)
		**out = **in
	}
	if in.LegacyUninstallHasRun != nil {
		in, out := &in.LegacyUninstallHasRun, &out.LegacyUninstallHasRun
		*out = new(bool)
		**out = **in
	}
	if in.Features != nil {
		in, out := &in.Features, &out.Features
		*out = new(common.Features)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RazeeDeploymentSpec.
func (in *RazeeDeploymentSpec) DeepCopy() *RazeeDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(RazeeDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RazeeDeploymentStatus) DeepCopyInto(out *RazeeDeploymentStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(status.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.JobConditions != nil {
		in, out := &in.JobConditions, &out.JobConditions
		*out = new(batchv1.JobCondition)
		**out = **in
	}
	if in.JobState != nil {
		in, out := &in.JobState, &out.JobState
		*out = new(batchv1.JobStatus)
		(*in).DeepCopyInto(*out)
	}
	if in.MissingDeploySecretValues != nil {
		in, out := &in.MissingDeploySecretValues, &out.MissingDeploySecretValues
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RazeePrerequisitesCreated != nil {
		in, out := &in.RazeePrerequisitesCreated, &out.RazeePrerequisitesCreated
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.LocalSecretVarsPopulated != nil {
		in, out := &in.LocalSecretVarsPopulated, &out.LocalSecretVarsPopulated
		*out = new(bool)
		**out = **in
	}
	if in.RedHatMarketplaceSecretFound != nil {
		in, out := &in.RedHatMarketplaceSecretFound, &out.RedHatMarketplaceSecretFound
		*out = new(bool)
		**out = **in
	}
	if in.RazeeJobInstall != nil {
		in, out := &in.RazeeJobInstall, &out.RazeeJobInstall
		*out = new(RazeeJobInstallStruct)
		**out = **in
	}
	if in.NodesFromRazeeDeployments != nil {
		in, out := &in.NodesFromRazeeDeployments, &out.NodesFromRazeeDeployments
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RazeeDeploymentStatus.
func (in *RazeeDeploymentStatus) DeepCopy() *RazeeDeploymentStatus {
	if in == nil {
		return nil
	}
	out := new(RazeeDeploymentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RazeeJobInstallStruct) DeepCopyInto(out *RazeeJobInstallStruct) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RazeeJobInstallStruct.
func (in *RazeeJobInstallStruct) DeepCopy() *RazeeJobInstallStruct {
	if in == nil {
		return nil
	}
	out := new(RazeeJobInstallStruct)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RazeeLogs) DeepCopyInto(out *RazeeLogs) {
	*out = *in
	if in.Log != nil {
		in, out := &in.Log, &out.Log
		*out = make(Log, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RazeeLogs.
func (in *RazeeLogs) DeepCopy() *RazeeLogs {
	if in == nil {
		return nil
	}
	out := new(RazeeLogs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceS3) DeepCopyInto(out *RemoteResourceS3) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceS3.
func (in *RemoteResourceS3) DeepCopy() *RemoteResourceS3 {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceS3)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RemoteResourceS3) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceS3List) DeepCopyInto(out *RemoteResourceS3List) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RemoteResourceS3, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceS3List.
func (in *RemoteResourceS3List) DeepCopy() *RemoteResourceS3List {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceS3List)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RemoteResourceS3List) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceS3Spec) DeepCopyInto(out *RemoteResourceS3Spec) {
	*out = *in
	in.Auth.DeepCopyInto(&out.Auth)
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = make([]Request, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceS3Spec.
func (in *RemoteResourceS3Spec) DeepCopy() *RemoteResourceS3Spec {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceS3Spec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteResourceS3Status) DeepCopyInto(out *RemoteResourceS3Status) {
	*out = *in
	if in.Touched != nil {
		in, out := &in.Touched, &out.Touched
		*out = new(bool)
		**out = **in
	}
	in.RazeeLogs.DeepCopyInto(&out.RazeeLogs)
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(status.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteResourceS3Status.
func (in *RemoteResourceS3Status) DeepCopy() *RemoteResourceS3Status {
	if in == nil {
		return nil
	}
	out := new(RemoteResourceS3Status)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Request) DeepCopyInto(out *Request) {
	*out = *in
	in.Options.DeepCopyInto(&out.Options)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Request.
func (in *Request) DeepCopy() *Request {
	if in == nil {
		return nil
	}
	out := new(Request)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *S3Options) DeepCopyInto(out *S3Options) {
	*out = *in
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make(map[string]Header, len(*in))
		for key, val := range *in {
			var outVal map[string]string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make(Header, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new S3Options.
func (in *S3Options) DeepCopy() *S3Options {
	if in == nil {
		return nil
	}
	out := new(S3Options)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretAccessKeyRef) DeepCopyInto(out *SecretAccessKeyRef) {
	*out = *in
	in.ValueFrom.DeepCopyInto(&out.ValueFrom)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretAccessKeyRef.
func (in *SecretAccessKeyRef) DeepCopy() *SecretAccessKeyRef {
	if in == nil {
		return nil
	}
	out := new(SecretAccessKeyRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageSpec) DeepCopyInto(out *StorageSpec) {
	*out = *in
	if in.Class != nil {
		in, out := &in.Class, &out.Class
		*out = new(string)
		**out = **in
	}
	out.Size = in.Size.DeepCopy()
	if in.EmptyDir != nil {
		in, out := &in.EmptyDir, &out.EmptyDir
		*out = new(corev1.EmptyDirVolumeSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageSpec.
func (in *StorageSpec) DeepCopy() *StorageSpec {
	if in == nil {
		return nil
	}
	out := new(StorageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValueFrom) DeepCopyInto(out *ValueFrom) {
	*out = *in
	in.SecretKeyRef.DeepCopyInto(&out.SecretKeyRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValueFrom.
func (in *ValueFrom) DeepCopy() *ValueFrom {
	if in == nil {
		return nil
	}
	out := new(ValueFrom)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Workload) DeepCopyInto(out *Workload) {
	*out = *in
	if in.OwnerCRD != nil {
		in, out := &in.OwnerCRD, &out.OwnerCRD
		*out = new(common.GroupVersionKind)
		**out = **in
	}
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.AnnotationSelector != nil {
		in, out := &in.AnnotationSelector, &out.AnnotationSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.MetricLabels != nil {
		in, out := &in.MetricLabels, &out.MetricLabels
		*out = make([]MeterLabelQuery, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Workload.
func (in *Workload) DeepCopy() *Workload {
	if in == nil {
		return nil
	}
	out := new(Workload)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadStatus) DeepCopyInto(out *WorkloadStatus) {
	*out = *in
	in.LastReadTime.DeepCopyInto(&out.LastReadTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadStatus.
func (in *WorkloadStatus) DeepCopy() *WorkloadStatus {
	if in == nil {
		return nil
	}
	out := new(WorkloadStatus)
	in.DeepCopyInto(out)
	return out
}
