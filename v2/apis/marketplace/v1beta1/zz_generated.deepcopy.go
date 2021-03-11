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

package v1beta1

import (
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/common"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/status"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnnotationFilter) DeepCopyInto(out *AnnotationFilter) {
	*out = *in
	if in.AnnotationSelector != nil {
		in, out := &in.AnnotationSelector, &out.AnnotationSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnnotationFilter.
func (in *AnnotationFilter) DeepCopy() *AnnotationFilter {
	if in == nil {
		return nil
	}
	out := new(AnnotationFilter)
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
func (in *LabelFilter) DeepCopyInto(out *LabelFilter) {
	*out = *in
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LabelFilter.
func (in *LabelFilter) DeepCopy() *LabelFilter {
	if in == nil {
		return nil
	}
	out := new(LabelFilter)
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
	if in.ResourceFilters != nil {
		in, out := &in.ResourceFilters, &out.ResourceFilters
		*out = make([]ResourceFilter, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Meters != nil {
		in, out := &in.Meters, &out.Meters
		*out = make([]MeterWorkload, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InstalledBy != nil {
		in, out := &in.InstalledBy, &out.InstalledBy
		*out = new(common.NamespacedNameReference)
		(*in).DeepCopyInto(*out)
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
func (in *MeterWorkload) DeepCopyInto(out *MeterWorkload) {
	*out = *in
	if in.GroupBy != nil {
		in, out := &in.GroupBy, &out.GroupBy
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Without != nil {
		in, out := &in.Without, &out.Without
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Period != nil {
		in, out := &in.Period, &out.Period
		*out = new(v1.Duration)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeterWorkload.
func (in *MeterWorkload) DeepCopy() *MeterWorkload {
	if in == nil {
		return nil
	}
	out := new(MeterWorkload)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceFilter) DeepCopyInto(out *NamespaceFilter) {
	*out = *in
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceFilter.
func (in *NamespaceFilter) DeepCopy() *NamespaceFilter {
	if in == nil {
		return nil
	}
	out := new(NamespaceFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OwnerCRDFilter) DeepCopyInto(out *OwnerCRDFilter) {
	*out = *in
	out.GroupVersionKind = in.GroupVersionKind
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OwnerCRDFilter.
func (in *OwnerCRDFilter) DeepCopy() *OwnerCRDFilter {
	if in == nil {
		return nil
	}
	out := new(OwnerCRDFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceFilter) DeepCopyInto(out *ResourceFilter) {
	*out = *in
	if in.Namespace != nil {
		in, out := &in.Namespace, &out.Namespace
		*out = new(NamespaceFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.OwnerCRD != nil {
		in, out := &in.OwnerCRD, &out.OwnerCRD
		*out = new(OwnerCRDFilter)
		**out = **in
	}
	if in.Label != nil {
		in, out := &in.Label, &out.Label
		*out = new(LabelFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.Annotation != nil {
		in, out := &in.Annotation, &out.Annotation
		*out = new(AnnotationFilter)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceFilter.
func (in *ResourceFilter) DeepCopy() *ResourceFilter {
	if in == nil {
		return nil
	}
	out := new(ResourceFilter)
	in.DeepCopyInto(out)
	return out
}
