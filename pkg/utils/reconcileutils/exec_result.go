package reconcileutils

import (
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ActionResultStatus string

var (
	Continue ActionResultStatus = "continue"
	NotFound ActionResultStatus = "not_found"
	Requeue  ActionResultStatus = "requeue"
	Error    ActionResultStatus = "error"
)

type ExecResult struct {
	Status          ActionResultStatus
	ReconcileResult reconcile.Result
	Err             error
}

func (e *ExecResult) GetResult() ActionResultStatus {
	return e.Status
}

func (e *ExecResult) Is(comp ActionResultStatus) bool {
	return e.Status == comp
}

func (e *ExecResult) GetReconcile() reconcile.Result {
	return e.ReconcileResult
}

func (e *ExecResult) GetError() error {
	return e.Err
}

func (e *ExecResult) Return() (reconcile.Result, error) {
	return e.ReconcileResult, e.Err
}

func NewExecResult(
	status ActionResultStatus,
	reconcileResult reconcile.Result,
	err error) *ExecResult {
	return &ExecResult{
		Status:          status,
		ReconcileResult: reconcileResult,
		Err:             err,
	}
}
