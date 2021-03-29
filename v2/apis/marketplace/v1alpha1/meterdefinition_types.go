// Copyright 2020 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"bytes"

	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/common"
	status "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// MeterDefinitionSpec defines the desired metering spec
// +k8s:openapi-gen=true
type MeterDefinitionSpec struct {
	// Group defines the operator group of the meter
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Group string `json:"meterGroup"`

	// Kind defines the primary CRD kind of the meter
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Kind string `json:"meterKind"`

	// InstalledBy is a reference to the CSV that install the meter
	// definition. This is used to determine an operator group.
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:hidden"
	// +optional
	InstalledBy *common.NamespacedNameReference `json:"installedBy,omitempty"`

	// WorkloadVertexType is the top most object of a workload. It allows
	// you to identify the upper bounds of your workloads.
	// +kubebuilder:validation:Enum=Namespace;OperatorGroup
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:select:Namespace,urn:alm:descriptor:com.tectonic.ui:select:OperatorGroup"
	WorkloadVertexType WorkloadVertex `json:"workloadVertexType,omitempty"`

	// VertexFilters are used when Namespace is selected. Can be omitted
	// if you select OperatorGroup
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:fieldDependency:workloadVertexType:Namespace"
	// +optional
	VertexLabelSelector *metav1.LabelSelector `json:"workloadVertexLabelSelectors,omitempty"`

	// Workloads identify the workloads to meter.
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +kubebuilder:validation:MinItems=1
	Workloads []Workload `json:"workloads,omitempty"`

	// Version defines the primary CRD version of the meter. This field is no longer used.
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:hidden"
	// +optional
	Version *string `json:"meterVersion,omitempty"`

	// ServiceMeterLabels name of the meterics you want to track. Use workloads instead.
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:hidden"
	// +optional
	ServiceMeterLabels []string `json:"serviceMeterLabels,omitempty"`

	// PodMeterLabels name of the prometheus metrics you want to track. User workloads instead.
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:hidden"
	// +optional
	PodMeterLabels []string `json:"podMeterLabels,omitempty"`
}

const (
	WorkloadVertexOperatorGroup WorkloadVertex = "OperatorGroup"
	WorkloadVertexNamespace                    = "Namespace"
)
const (
	WorkloadTypePod            WorkloadType = "Pod"
	WorkloadTypeService        WorkloadType = "Service"
	WorkloadTypeServiceMonitor WorkloadType = "ServiceMonitor"
	WorkloadTypePVC            WorkloadType = "PersistentVolumeClaim"
)

type WorkloadVertex string
type WorkloadType string
type CSVNamespacedName common.NamespacedNameReference

// Workload helps identify what to target for metering.
type Workload struct {
	// Name of the workload, must be unique in a meter definition.
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Name string `json:"name"`

	// WorkloadType identifies the type of workload to look for. This can be
	// pod or service right now.
	// +kubebuilder:validation:Enum=Pod;Service;PersistentVolumeClaim
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:select:Pod,urn:alm:descriptor:com.tectonic.ui:select:Service,urn:alm:descriptor:com.tectonic.ui:select:PersistentVolumeClaim"
	WorkloadType WorkloadType `json:"type"`

	// OwnerCRD is the name of the GVK to look for as the owner of all the
	// meterable assets. If omitted, the labels and annotations are used instead.
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	OwnerCRD *common.GroupVersionKind `json:"ownerCRD,omitempty"`

	// LabelSelector are used to filter to the correct workload.
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty"`

	// AnnotationSelector are used to filter to the correct workload.
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	AnnotationSelector *metav1.LabelSelector `json:"annotationSelector,omitempty"`

	// MetricLabels are the labels to collect
	// +required
	// +kubebuilder:validation:MinItems=1
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	MetricLabels []MeterLabelQuery `json:"metricLabels,omitempty"`
}

// WorkloadStatus provides quick status to check if
// workloads are working correctly
type WorkloadStatus struct {
	// Name of the workload, must be unique in a meter definition.
	Name string `json:"name"`

	CurrentMetricValue string `json:"currentValue"`

	LastReadTime metav1.Time `json:"startTime"`
}

// MeterLabelQuery helps define a meter label to build and search for
type MeterLabelQuery struct {
	// Label is the name of the meter
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Label string `json:"label"`

	// Query to use for the label
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Query string `json:"query,omitempty"`

	// Aggregation to use with the query
	// +kubebuilder:validation:Enum:=sum;min;max;avg
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:select:sum,urn:alm:descriptor:com.tectonic.ui:select:min,urn:alm:descriptor:com.tectonic.ui:select:max,urn:alm:descriptor:com.tectonic.ui:select:avg"
	Aggregation string `json:"aggregation,omitempty"`
}

// MeterDefinitionStatus defines the observed state of MeterDefinition
// +k8s:openapi-gen=true
// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
type MeterDefinitionStatus struct {

	// Conditions represent the latest available observations of an object's state
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes.conditions"
	// +optional
	Conditions status.Conditions `json:"conditions,omitempty"`

	// WorkloadResources is the list of resoruces discovered by
	// this meter definition
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	WorkloadResources []common.WorkloadResource `json:"workloadResource,omitempty"`

	// Results is a list of Results that get returned from a query to prometheus
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	Results []common.Result `json:"results,omitempty"`
}

// MeterDefinition defines the meter workloads used to enable pay for
// use billing.
// +kubebuilder:object:root=true
//
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=meterdefinitions,scope=Namespaced
// +operator-sdk:gen-csv:customresourcedefinitions.displayName="Meter Definitions"
// +genclient
type MeterDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MeterDefinitionSpec   `json:"spec,omitempty"`
	Status MeterDefinitionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MeterDefinitionList contains a list of MeterDefinition
type MeterDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MeterDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MeterDefinition{}, &MeterDefinitionList{})
}

func (meterdef *MeterDefinition) BuildMeterDefinitionFromString(meterdefString, name, namespace, nameLabel, namespaceLabel string) error {
	data := []byte(meterdefString)
	err := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 100).Decode(meterdef)
	if err != nil {
		return err
	}

	csvInfo := make(map[string]string)
	csvInfo[nameLabel] = name
	csvInfo[namespaceLabel] = namespace
	meterdef.SetAnnotations(csvInfo)

	meterdef.Namespace = namespace
	meterdef.Spec.InstalledBy = &common.NamespacedNameReference{
		Name:      name,
		Namespace: namespace,
	}

	return nil
}

func (meterdef *MeterDefinition) ToPrometheusLabels() []map[string]string {
	allMdefs := []map[string]string{}

	meterDefinitionUID := string(meterdef.UID)
	meterGroup := string(meterdef.Spec.Group)
	meterKind := string(meterdef.Spec.Kind)
	for _, workload := range meterdef.Spec.Workloads {
		for _, metricLabel := range workload.MetricLabels {
			labels := map[string]string{
				"meter_definition_uid": meterDefinitionUID,
				"workload_vertex_type": string(meterdef.Spec.WorkloadVertexType),
				"workload_name":        workload.Name,
				"workload_type":        string(workload.WorkloadType),
				"metric_label":         metricLabel.Label,
				"metric_aggregation":   metricLabel.Aggregation,
				"metric_query":         metricLabel.Query,
				"meter_group":          meterGroup,
				"meter_kind":           meterKind,
				"metric_period":        "1h",
			}

			allMdefs = append(allMdefs, labels)
		}
	}

	return allMdefs
}
