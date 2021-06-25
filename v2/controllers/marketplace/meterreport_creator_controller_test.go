package marketplace

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/common"
	marketplacev1alpha1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1alpha1"
	marketplacev1beta1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/reconcileutils"
	status "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("MeterbaseController", func() {
	Describe("check date functions", func() {
		var (
			ctrl *MeterReportCreatorReconciler
		)

		meterDef1 := &marketplacev1beta1.MeterDefinition{}
		meterDef2 := &marketplacev1beta1.MeterDefinition{}
		meterDef3 := &marketplacev1beta1.MeterDefinition{}
		catNameA := "labelA"
		mA := make(map[string]string)
		mA["marketplace.redhat.com/category"] = catNameA
		catNameB := "labelB"
		mB := make(map[string]string)
		mB["marketplace.redhat.com/category"] = catNameB

		meterDef1.SetLabels(mA)
		meterDef2.SetLabels(mB)
		meterDef3.SetLabels(mA)
		meterDefinitionList := []marketplacev1beta1.MeterDefinition{*meterDef1, *meterDef2, *meterDef3}

		endDate := time.Now().UTC()
		endDate = endDate.AddDate(0, 0, 0)
		minDate := endDate.AddDate(0, 0, 0)

		BeforeEach(func() {
			ctrl = &MeterReportCreatorReconciler{}
		})

		It("reports should calculate the correct dates to create", func() {

			exp := ctrl.generateExpectedDates(endDate, time.UTC, -30, minDate)
			Expect(exp).To(HaveLen(1))

			minDate = endDate.AddDate(0, 0, -2)

			exp = ctrl.generateExpectedDates(endDate, time.UTC, -30, minDate)
			Expect(exp).To(HaveLen(3))
		})

		It("should return category names, excluding duplicated", func() {
			Expect(meterDefinitionList).To(HaveLen((3)))
			categoryList := getCategoriesFromMeterDefinitions(meterDefinitionList)
			Expect(categoryList).To(HaveLen(2))
			Expect(categoryList).To(ContainElements([]string{catNameA, catNameB}))
		})

		It("should put category name into newly created meter report name", func() {
			endDate := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)
			minDate := endDate.AddDate(0, 0, 0)
			exp := ctrl.generateExpectedDates(endDate, time.UTC, -30, minDate)
			nameFromString := ctrl.newMeterReportNameFromString(catNameA, exp[0])
			Expect(nameFromString).To(Equal("meter-report-labelA-2021-06-01"))
		})

		It("should retrieve date properly for names with and without category in report name", func() {
			endDate := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)
			// old report name (without category)
			foundTime, _ := ctrl.retrieveCreatedDate("meter-report-2021-06-01")
			Expect(foundTime).To(Equal(endDate))
			// new report name (with category)
			foundTime2, _ := ctrl.retrieveCreatedDate("meter-report-label-2021-06-01")
			Expect(foundTime2).To(Equal(endDate))
			// if wrong format return actuall local time and error
			foundTime3, err := ctrl.retrieveCreatedDate("meter-2021-06-01")
			Expect(foundTime3.Unix()).To(Equal(time.Now().Unix()))
			Expect(err).NotTo(BeNil())
		})

		It("should return only non-duplicated categories", func() {
			categoryList := getCategoriesFromMeterDefinitions(meterDefinitionList)
			Expect(categoryList).To(HaveLen(2))
			Expect(categoryList).To(ContainElements([]string{catNameA, catNameB}))
		})
	})

	FDescribe("check reconciller", func() {
		var (
			name      = utils.METERBASE_NAME
			namespace = "openshift-redhat-marketplace"
			created   *marketplacev1alpha1.MeterBase
		)

		key := types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		}

		BeforeEach(func() {
			reportCreatorReconciler := &MeterReportCreatorReconciler{
				Log:    ctrl.Log.WithName("controllers").WithName("MeterReportCreator"),
				Client: k8sClient,
				Scheme: k8sScheme,
				CC:     reconcileutils.NewClientCommand(k8sManager.GetClient(), k8sScheme, ctrl.Log),
				cfg:    operatorCfg,
			}
			err := reportCreatorReconciler.SetupWithManager(k8sManager, doneChan)
			Expect(err).ToNot(HaveOccurred())

			k8sClient.Create(context.TODO(), &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
				Name: "openshift-redhat-marketplace",
			}})

			created = &marketplacev1alpha1.MeterBase{
				ObjectMeta: metav1.ObjectMeta{Name: key.Name, Namespace: key.Namespace},
				Spec: marketplacev1alpha1.MeterBaseSpec{
					Enabled: true,
				},
			}

			Expect(k8sClient.Create(context.TODO(), created)).Should(Succeed())

			created.Status = marketplacev1alpha1.MeterBaseStatus{
				Conditions: status.Conditions{
					status.Condition{
						Type:               marketplacev1alpha1.ConditionInstalling,
						Status:             corev1.ConditionFalse,
						Reason:             marketplacev1alpha1.ReasonMeterBaseFinishInstall,
						Message:            "finished install",
						LastTransitionTime: metav1.Now(),
					},
				},
			}

			Expect(k8sClient.Status().Update(context.TODO(), created)).Should(Succeed())

			spec := marketplacev1beta1.MeterDefinitionSpec{
				Group: "app.partner.metering.com",
				Kind:  "App",
				ResourceFilters: []marketplacev1beta1.ResourceFilter{
					{
						Namespace: &marketplacev1beta1.NamespaceFilter{UseOperatorGroup: true},
						OwnerCRD: &marketplacev1beta1.OwnerCRDFilter{
							GroupVersionKind: common.GroupVersionKind{
								APIVersion: "apps.partner.metering.com/v1",
								Kind:       "App",
							},
						},
						WorkloadType: marketplacev1beta1.WorkloadTypePod,
					},
				},
				Meters: []marketplacev1beta1.MeterWorkload{
					{
						Aggregation:  "sum",
						Query:        "simple_query",
						Metric:       "rpc_durations_seconds",
						Label:        "{{ .Label.meter_query }}",
						WorkloadType: marketplacev1beta1.WorkloadTypePod,
						Description:  "{{ .Label.meter_domain | lower }} description",
					},
				},
			}

			meterDef1 := &marketplacev1beta1.MeterDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mdef1",
					Namespace: "openshift-redhat-marketplace",
				},
				Spec: spec,
			}
			meterDef2 := &marketplacev1beta1.MeterDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mdef2",
					Namespace: "openshift-redhat-marketplace",
				},
				Spec: spec,
			}
			meterDef3 := &marketplacev1beta1.MeterDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mdef3",
					Namespace: "openshift-redhat-marketplace",
				},
				Spec: spec,
			}

			catNameA := "labelA"
			mA := make(map[string]string)
			mA["marketplace.redhat.com/category"] = catNameA
			catNameB := "labelB"
			mB := make(map[string]string)
			mB["marketplace.redhat.com/category"] = catNameB

			meterDef1.SetLabels(mA)
			meterDef2.SetLabels(mB)
			meterDef3.SetLabels(mA)
			meterDefinitionList := []marketplacev1beta1.MeterDefinition{*meterDef1, *meterDef2, *meterDef3}

			for i := range meterDefinitionList {
				Expect(k8sClient.Create(context.TODO(), &meterDefinitionList[i])).Should(Succeed())
			}
		})

		AfterEach(func() {
			// Add any teardown steps that needs to be executed after each test
			Expect(k8sClient.Delete(context.TODO(), created)).Should(Succeed())
		})

		It("should run meterbase reconciler", func() {
			k8sClient.Get(context.TODO(), key, created)
			By("Expecting status c")
			Eventually(func() *status.Condition {
				f := &marketplacev1alpha1.MeterBase{}
				k8sClient.Get(context.TODO(), key, f)
				return f.Status.Conditions.GetCondition(marketplacev1alpha1.ConditionInstalling)
			}, timeout, interval).Should(PointTo(MatchFields(IgnoreExtras, Fields{
				"Status": Equal(corev1.ConditionFalse),
				"Reason": Equal(marketplacev1alpha1.ReasonMeterBaseFinishInstall),
			})))

			Eventually(func() []marketplacev1alpha1.MeterReport {
				meterReportList := &marketplacev1alpha1.MeterReportList{}
				err := k8sClient.List(context.TODO(), meterReportList, client.InNamespace("openshift-redhat-marketplace"))
				Expect(err).To(Not(HaveOccurred()))
				return meterReportList.Items
			}, timeout, interval).Should(Not(BeEmpty()))
		})
	})
})
