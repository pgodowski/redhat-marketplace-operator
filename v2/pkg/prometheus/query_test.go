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

package prometheus

import (
	"time"

	"github.com/prometheus/common/model"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	. "github.com/redhat-marketplace/redhat-marketplace-operator/v2/tests/mock/mock_query"
)

var _ = Describe("Query", func() {

	var (
		v1api         v1.API
		prometheusAPI PrometheusAPI
		start, _      = time.Parse(time.RFC3339, "2020-04-19T13:00:00Z")
		end, _        = time.Parse(time.RFC3339, "2020-04-19T16:00:00Z")

		testQuery *PromQuery
	)

	BeforeEach(func() {
		testQuery = NewPromQuery(&PromQueryArgs{
			Metric: "rpc_durations_seconds_count",
			Query:  `foo{bar="true"}`,
			Type:   v1beta1.WorkloadTypePod,
			Start:  start,
			End:    end,
			Step:   time.Minute * 60,
		})

		v1api = GetTestAPI(MockResponseRoundTripper("../../../reporter/v2/test/mockresponses/prometheus-query-range.json", []v1beta1.MeterDefinition{}))
		prometheusAPI = PrometheusAPI{
			v1api,
		}
	})

	It("should query a range", func() {
		result, warnings, err := prometheusAPI.ReportQuery(testQuery)

		Expect(err).To(Succeed())
		Expect(warnings).To(BeEmpty(), "warnings should be empty")
		Expect(model.ValMatrix).To(Equal(result.Type()), "value type matrix expected")

		matrixResult, ok := result.(model.Matrix)

		Expect(ok).To(BeTrue(), "result is not a matrix")
		Expect(len(matrixResult)).To(Equal(2))
	})

	It("should build a query", func() {
		q1 := NewPromQuery(&PromQueryArgs{
			Metric: "foo",
			Query:  "kube_persistentvolumeclaim_resource_requests_storage_bytes",
			MeterDef: types.NamespacedName{
				Name:      "foo",
				Namespace: "foons",
			},
			AggregateFunc: "sum",
			Type:          v1beta1.WorkloadTypePVC,
		})

		expected := `sum by (namespace,persistentvolumeclaim) (avg(meterdef_persistentvolumeclaim_info{meter_def_name="foo",meter_def_namespace="foons",phase="Bound"}) without(cluster_ip,container,endpoint,instance,job,pod,pod_ip,pod_uid,service) * on(namespace,persistentvolumeclaim) group_right kube_persistentvolumeclaim_resource_requests_storage_bytes) * on(namespace,persistentvolumeclaim) group_right group(kube_persistentvolumeclaim_resource_requests_storage_bytes) without(cluster_ip,container,endpoint,instance,job)`
		q, err := q1.Print()
		Expect(err).To(Succeed())
		Expect(q).To(Equal(expected), "failed to create query for pvc")
	})

	It("should build a query with a groupby and without", func() {
		q1 := NewPromQuery(&PromQueryArgs{
			Metric: "foo",
			Query:  "kube_persistentvolumeclaim_resource_requests_storage_bytes",
			MeterDef: types.NamespacedName{
				Name:      "foo",
				Namespace: "foons",
			},
			AggregateFunc: "sum",
			Type:          v1beta1.WorkloadTypePVC,
			GroupBy:       []string{"foo"},
			Without:       []string{"bar"},
		})

		expected := `sum by (foo) (avg(meterdef_persistentvolumeclaim_info{meter_def_name="foo",meter_def_namespace="foons",phase="Bound"}) without(cluster_ip,container,endpoint,instance,job,pod,pod_ip,pod_uid,service) * on(persistentvolumeclaim,namespace) group_right kube_persistentvolumeclaim_resource_requests_storage_bytes) * on(foo) group_right group(kube_persistentvolumeclaim_resource_requests_storage_bytes) without(bar,cluster_ip,container,endpoint,instance,job)`
		q, err := q1.Print()
		Expect(err).To(Succeed())
		Expect(q).To(Equal(expected), "failed to create query for pvc")
	})

	It("should handle groupby clauses", func() {
		q1 := NewPromQuery(&PromQueryArgs{
			Metric: "foo",
			Query:  "kube_persistentvolumeclaim_resource_requests_storage_bytes",
			MeterDef: types.NamespacedName{
				Name:      "foo",
				Namespace: "foons",
			},
			AggregateFunc: "sum",
			Type:          v1beta1.WorkloadTypePVC,
			GroupBy:       []string{"persistentvolumeclaim"},
		})

		expected := `sum by (persistentvolumeclaim) (avg(meterdef_persistentvolumeclaim_info{meter_def_name="foo",meter_def_namespace="foons",phase="Bound"}) without(cluster_ip,container,endpoint,instance,job,pod,pod_ip,pod_uid,service) * on(persistentvolumeclaim,namespace) group_right kube_persistentvolumeclaim_resource_requests_storage_bytes) * on(persistentvolumeclaim) group_right group(kube_persistentvolumeclaim_resource_requests_storage_bytes) without(cluster_ip,cluster_ip,container,container,endpoint,endpoint,instance,instance,job,job,namespace,pod,pod_ip,pod_uid,service)`
		q, err := q1.Print()
		Expect(err).To(Succeed())
		Expect(q).To(Equal(expected), "failed to create query for pvc")
	})
})
