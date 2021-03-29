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

package controller_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/common"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1alpha1"
	. "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/reconcileutils"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/status"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe("MeterReportController", func() {
	Context("MeterReport reconcile", func() {
		var (
			meterreport *v1alpha1.MeterReport
			meterdef    *v1alpha1.MeterDefinition
			start       time.Time
			end         time.Time
		)

		BeforeEach(func() {
			start, end = time.Now(), time.Now()
			start.Add(-5 * time.Minute)

			meterdef = &v1alpha1.MeterDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: Namespace,
				},
				Spec: v1alpha1.MeterDefinitionSpec{
					Group:              "testgroup",
					Kind:               "testkind",
					WorkloadVertexType: v1alpha1.WorkloadVertexOperatorGroup,
					Workloads: []v1alpha1.Workload{
						{
							Name:         "test",
							WorkloadType: v1alpha1.WorkloadTypePod,
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/name": "rhm-metric-state",
								},
							},
							MetricLabels: []v1alpha1.MeterLabelQuery{
								{
									Aggregation: "sum",
									Label:       "test",
									Query:       "kube_pod_info",
								},
							},
						},
					},
				},
			}

			meterreport = &v1alpha1.MeterReport{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test1",
					Namespace: "openshift-redhat-marketplace",
				},
				Spec: v1alpha1.MeterReportSpec{
					StartTime: metav1.NewTime(start),
					EndTime:   metav1.NewTime(end),
					PrometheusService: &common.ServiceReference{
						Namespace:  Namespace,
						Name:       "rhm-prometheus-meterbase",
						TargetPort: intstr.FromString("rbac"),
					},
					ExtraArgs: []string{
						"--uploadTarget", "noop",
						"--brokenarg",
					},
				},
			}

			testHarness.Delete(context.TODO(), meterreport)
			Expect(testHarness.Create(context.TODO(), meterdef)).Should(SucceedOrAlreadyExist)
		})

		AfterEach(func() {
			Expect(testHarness.Delete(context.TODO(), meterdef)).Should(Succeed())
			Expect(testHarness.Delete(context.TODO(), meterreport)).Should(Succeed())
		})

		It("should retry a job", func() {
			Expect(testHarness.Create(context.TODO(), meterreport)).Should(Succeed())
			job := &batchv1.Job{}

			Eventually(func() bool {
				result, _ := testHarness.Do(
					context.TODO(),
					GetAction(types.NamespacedName{Name: meterreport.Name, Namespace: Namespace}, job),
				)
				return result.Is(Continue)
			}, timeout, interval).Should(BeTrue())

			Expect("the job to fail and report on the condition")
			Eventually(func() map[string]interface{} {
				result, _ := testHarness.Do(
					context.TODO(),
					GetAction(types.NamespacedName{Name: meterreport.Name, Namespace: Namespace}, meterreport),
				)

				if !result.Is(Continue) {
					return map[string]interface{}{
						"result": result.Status,
					}
				}

				var reason status.ConditionReason
				for _, cond := range meterreport.Status.Conditions {
					if cond.Type == v1alpha1.ReportConditionTypeJobRunning {
						reason = cond.Reason
					}
				}

				return map[string]interface{}{
					"result": result.Status,
					"reason": reason,
				}
			}, timeout, interval).Should(
				MatchAllKeys(Keys{
					"result": Equal(Continue),
					"reason": Equal(v1alpha1.ReportConditionReasonJobErrored),
				}))

			job2 := &batchv1.Job{}

			// Expect a new job to be created
			Eventually(func() bool {
				result, _ := testHarness.Do(
					context.TODO(),
					GetAction(types.NamespacedName{Name: meterreport.Name, Namespace: Namespace}, job2),
				)

				if result.Is(Continue) {
					return job.UID != job2.UID
				}
				return false
			}, time.Second*300, interval).Should(BeTrue())

			Expect(job2.UID).To(Not(Equal(job.UID)))
		})
	})
})
