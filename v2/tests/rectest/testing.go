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

package rectest

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/tests/mock/mock_client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcileTester interface {
	Fail()
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Failed() bool
	Parallel()
	Skip(args ...interface{})
	Skipf(format string, args ...interface{})
	SkipNow()
	Skipped() bool
}

type OneOf struct {
	Object     client.Object
	ObjectList client.ObjectList
}

// - interfaces -
type ReconcilerTestValidationFunc func(*ReconcilerTest, ReconcileTester, client.Object)
type ReconcilerTestListValidationFunc func(*ReconcilerTest, ReconcileTester, client.ObjectList)
type ReconcilerSetupFunc func(*ReconcilerTest) error

type TestCaseStep interface {
	GetStepName() string
	Test(t ReconcileTester, reconcilerTest *ReconcilerTest)
}

// - end of interfaces -

// ReconcilerTest is the major test driver, create one of these for each test
type ReconcilerTest struct {
	runtimeObjs []client.Object
	SetupFunc   ReconcilerSetupFunc
	Reconciler  reconcile.Reconciler
	Client      client.Client
}

func (r *ReconcilerTest) SetReconciler(re reconcile.Reconciler) {
	r.Reconciler = re
}

func (r *ReconcilerTest) SetClient(c client.Client) {
	r.Client = c
}

func (r *ReconcilerTest) GetReconciler() reconcile.Reconciler {
	return r.Reconciler
}

func (r *ReconcilerTest) GetClient() client.Client {
	return r.Client
}

func (r *ReconcilerTest) GetGetObjects() []client.Object {
	return r.runtimeObjs
}

type ReconcileResult struct {
	reconcile.Result
	Err error
}

var testSetupLock sync.Mutex

func NewReconcilerTestSimple(
	reconciler reconcile.Reconciler,
	client client.Client,
) *ReconcilerTest {
	return &ReconcilerTest{
		Reconciler: reconciler,
		Client:     client,
	}
}

// NewReconcilerTest creates a new reconciler test with a setup func
// using the provided runtime objects to creat on the client.
func NewReconcilerTest(setup ReconcilerSetupFunc, predefinedObjs ...client.Object) *ReconcilerTest {
	testSetupLock.Lock()
	defer testSetupLock.Unlock()
	myObjs := []client.Object{}

	for _, obj := range predefinedObjs {
		myObjs = append(myObjs, obj)
	}

	return &ReconcilerTest{
		runtimeObjs: myObjs,
		SetupFunc:   setup,
	}
}

func Ignore(_ *ReconcilerTest, _ ReconcileTester, _ client.Object)         {}
func IgnoreList(_ *ReconcilerTest, _ ReconcileTester, _ client.ObjectList) {}

type ControllerReconcileStep struct {
	stepOptions
	reconcileStepOptions
	*testLine
}

func ReconcileStep(
	stepOptions []StepOption,
	options ...ReconcileStepOption,
) *ControllerReconcileStep {
	stepOpts, _ := newStepOptions(stepOptions...)
	opts, _ := newReconcileStepOptions(options...)

	return &ControllerReconcileStep{
		testLine:             NewTestLine("reconcileStep failure", 3),
		reconcileStepOptions: opts,
		stepOptions:          stepOpts,
	}
}

func (tc *ControllerReconcileStep) GetStepName() string {
	if tc.StepName == "" {
		return "ReconcileStep"
	}
	return tc.StepName
}

func (tc *ControllerReconcileStep) Test(t ReconcileTester, r *ReconcilerTest) {
	// Reconcile again so Reconcile() checks for the OperatorSource

	if tc.UntilDone {
		tc.Max = 1000
	}

	if tc.Max == 0 {
		tc.Max = len(tc.ExpectedResults)
	}

	for i := 0; i < tc.Max; i++ {
		exit := false

		indx := i

		ctx := context.Background()
		res, err := r.Reconciler.Reconcile(ctx, tc.Request)
		result := ReconcileResult{res, err}

		expectedResult := AnyResult

		if indx < len(tc.ExpectedResults) {
			expectedResult = tc.ExpectedResults[indx]
		}

		if expectedResult != AnyResult {
			assert.Equalf(t, expectedResult, result,
				"%+v", tc.TestLineError(fmt.Errorf("incorrect expected result")))
		} else {
			// stop if done or if there was an error
			if result == DoneResult {
				if len(tc.ExpectedResults) != 0 && indx >= len(tc.ExpectedResults) && !tc.UntilDone {
					assert.Equalf(t, len(tc.ExpectedResults)-1, indx,
						"%+v", tc.TestLineError(fmt.Errorf("expected reconcile count did not match")))
				}
				t.Logf("reconcile completed in %v turns", indx+1)
				exit = true
			}

			if !tc.IgnoreError {
				t.Logf("ignore error")
				if err != nil {
					assert.Equalf(t, DoneResult, result,
						"%+v", tc.TestLineError(err))
					exit = true
				}
			}

			if indx == tc.Max-1 {
				exit = true
			}
		}

		if exit {
			break
		}
	}
}

type ClientGetStep struct {
	*testLine
	stepOptions
	getStepOptions
}

func GetStep(
	stepOptions []StepOption,
	options ...GetStepOption,
) *ClientGetStep {
	stepOpts, _ := newStepOptions(stepOptions...)
	getOpts, _ := newGetStepOptions(options...)
	return &ClientGetStep{
		testLine:       NewTestLine("failed client get step", 3),
		stepOptions:    stepOpts,
		getStepOptions: getOpts,
	}
}

func (tc *ClientGetStep) GetStepName() string {
	if tc.StepName == "" {
		return "GetStep"
	}
	return tc.StepName
}

func (tc *ClientGetStep) Test(t ReconcileTester, r *ReconcilerTest) {
	// Reconcile again so Reconcile() checks for the OperatorSource
	err := r.GetClient().Get(
		context.TODO(),
		types.NamespacedName{
			Name:      tc.NamespacedName.Name,
			Namespace: tc.NamespacedName.Namespace,
		},
		tc.Obj,
	)

	require.NoErrorf(t, err, "get (%T): (%v); err=%+v", tc.Obj, err, tc.TestLineError(err))
	tc.CheckResult(r, t, tc.Obj)
}

type ClientListStep struct {
	*testLine
	stepOptions
	listStepOptions
}

func ListStep(
	stepOptions []StepOption,
	options ...ListStepOption,
) *ClientListStep {
	stepOpts, _ := newStepOptions(stepOptions...)
	listOpts, _ := newListStepOptions(options...)
	return &ClientListStep{
		testLine:        NewTestLine("failed client list step", 3),
		stepOptions:     stepOpts,
		listStepOptions: listOpts,
	}
}

func (tc *ClientListStep) GetStepName() string {
	if tc.StepName == "" {
		return "ListStep"
	}
	return tc.StepName
}

func (tc *ClientListStep) Test(t ReconcileTester, r *ReconcilerTest) {
	// Reconcile again so Reconcile() checks for the OperatorSource
	err := r.GetClient().List(context.TODO(),
		tc.Obj,
		tc.Filter...,
	)

	if err != nil {
		assert.FailNowf(t, "error encountered",
			"%+v", tc.TestLineError(errors.Errorf("get (%T): (%v)", tc.Obj, err)))
	}

	tc.CheckResult(r, t, tc.Obj)
}

var testAllMutex sync.Mutex

func (r *ReconcilerTest) TestAll(t ReconcileTester, testCases ...TestCaseStep) {
	if r.SetupFunc != nil {
		testAllMutex.Lock()
		err := r.SetupFunc(r)
		testAllMutex.Unlock()

		if err != nil {
			t.Fatalf("failed to setup test %v", err)
		}
	}

	for i, testData := range testCases {
		testName := fmt.Sprintf("%v %v", testData.GetStepName(), i+1)

		if testName == "" {
			testName = fmt.Sprintf("Step %v", i+1)
		}

		rectest := r
		testData := testData

		testData.Test(t, rectest)
	}
}

func getNamespacedName(action, version, name, namespace string) string {
	return fmt.Sprintf("%s-%s-%s.%s", action, version, name, namespace)
}

func ClientErrorStub(ctrl *gomock.Controller, clientImpl client.Client, mockErr error) client.Client {
	mock := mock_client.NewMockClient(ctrl)
	statusWriter := mock_client.NewMockStatusWriter(ctrl)
	called := make(map[string]bool)

	mock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
			name := getNamespacedName("create", obj.GetObjectKind().GroupVersionKind().String(), obj.GetName(), obj.GetNamespace())

			if _, ok := called[name]; !ok {
				called[name] = true
				return mockErr
			}

			return clientImpl.Create(ctx, obj, opts...)
		}).AnyTimes()

	mock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
			name := getNamespacedName("update", obj.GetObjectKind().GroupVersionKind().String(), obj.GetName(), obj.GetNamespace())

			if _, ok := called[name]; !ok {
				called[name] = true
				return mockErr
			}

			return clientImpl.Update(ctx, obj, opts...)
		}).AnyTimes()

	mock.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
			name := getNamespacedName("delete", obj.GetObjectKind().GroupVersionKind().String(), obj.GetName(), obj.GetNamespace())

			if _, ok := called[name]; !ok {
				called[name] = true
				return mockErr
			}

			return clientImpl.Delete(ctx, obj, opts...)
		}).AnyTimes()

	mock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, obj client.ObjectList, opts ...client.ListOption) error {
			name := getNamespacedName("list", obj.GetObjectKind().GroupVersionKind().String(), "", "")

			if _, ok := called[name]; !ok {
				called[name] = true
				return mockErr
			}

			return clientImpl.List(ctx, obj, opts...)
		}).AnyTimes()

	mock.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
			name := getNamespacedName("patch", obj.GetObjectKind().GroupVersionKind().String(), obj.GetName(), obj.GetNamespace())

			if _, ok := called[name]; !ok {
				called[name] = true
				return mockErr
			}

			return clientImpl.Patch(ctx, obj, patch, opts...)
		}).AnyTimes()

	mock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
			name := getNamespacedName("get", obj.GetObjectKind().GroupVersionKind().String(), obj.GetName(), obj.GetNamespace())

			if _, ok := called[name]; !ok {
				called[name] = true
				return mockErr
			}

			return clientImpl.Get(ctx, key, obj)
		}).AnyTimes()

	statusWriter.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
			name := getNamespacedName("patch-status", obj.GetObjectKind().GroupVersionKind().String(), obj.GetName(), obj.GetNamespace())

			if _, ok := called[name]; !ok {
				called[name] = true
				return mockErr
			}

			return clientImpl.Status().Patch(ctx, obj, patch, opts...)
		}).AnyTimes()

	statusWriter.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
			name := getNamespacedName("update-status", obj.GetObjectKind().GroupVersionKind().String(), obj.GetName(), obj.GetNamespace())

			if _, ok := called[name]; !ok {
				called[name] = true
				return mockErr
			}

			return clientImpl.Status().Update(ctx, obj, opts...)
		}).AnyTimes()

	mock.EXPECT().Status().Return(statusWriter).AnyTimes()

	return mock
}
