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

package manifests

import (
	"context"
	"fmt"

	emperrors "emperror.dev/errors"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/patch"
	. "github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/reconcileutils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type createOrUpdateFactoryItemAction struct {
	*BaseAction
	object      client.Object
	factoryFunc func() (client.Object, error)
	owner       client.Object
	patcher     patch.Patcher
}

type CreateOrUpdateFactoryItemArgs struct {
	Owner   client.Object
	Patcher patch.Patcher
}

func CreateOrUpdateFactoryItemAction(
	newObj client.Object,
	factoryFunc func() (client.Object, error),
	args CreateOrUpdateFactoryItemArgs,
) *createOrUpdateFactoryItemAction {
	return &createOrUpdateFactoryItemAction{
		BaseAction:  NewBaseAction("createOrUpdateFactoryItem"),
		object:      newObj,
		factoryFunc: factoryFunc,
		owner:       args.Owner,
		patcher:     args.Patcher,
	}
}

func (a *createOrUpdateFactoryItemAction) Bind(result *ExecResult) {
	a.SetLastResult(result)
}

func (a *createOrUpdateFactoryItemAction) Exec(ctx context.Context, c *ClientCommand) (*ExecResult, error) {
	reqLogger := a.GetReqLogger(c)
	result, err := a.factoryFunc()

	if err != nil {
		reqLogger.Error(err, "failure creating factory obj")
		return NewExecResult(Error, reconcile.Result{Requeue: true}, a.BaseAction, err), emperrors.Wrap(err, "error with patch")
	}

	key := client.ObjectKeyFromObject(result)

	if err != nil {
		reqLogger.Error(err, "failure getting factory obj name")
		return NewExecResult(Error, reconcile.Result{Requeue: true}, a.BaseAction, err), emperrors.Wrap(err, "error with patch")
	}

	cmd := HandleResult(
		GetAction(key, a.object),
		OnNotFound(CreateAction(result, CreateWithPatch(a.patcher))),
		OnContinue(Call(func() (ClientAction, error) {
			if a.owner != nil {
				if meta, ok := a.owner.(metav1.Object); ok {
					if err := controllerutil.SetControllerReference(
						a.owner.(metav1.Object),
						meta,
						c.GetScheme()); err != nil {
					}
				}
			}

			// handle case if original config is missing
			if orig, _ := a.patcher.GetOriginalConfiguration(a.object); orig == nil {
				data, _ := a.patcher.GetModifiedConfiguration(a.object, false)
				a.patcher.SetOriginalConfiguration(a.object, data)
			}

			patch, err := a.patcher.Calculate(a.object, result)
			if err != nil {
				return nil, emperrors.Wrap(err, "error creating patch")
			}

			if patch.IsEmpty() {
				return nil, nil
			}

			err = a.patcher.SetLastAppliedAnnotation(result)
			if err != nil {
				return nil, emperrors.Wrap(err, "error creating patch")
			}

			patch, err = a.patcher.Calculate(a.object, result)
			if err != nil {
				return nil, emperrors.Wrap(err, "error creating patch")
			}

			if patch.IsEmpty() {
				return nil, nil
			}

			reqLogger.Info("updating with patch",
				"patch", string(patch.Patch),
			)

			jsonPatch, err := jsonmergepatch.CreateThreeWayJSONMergePatch(patch.Original, patch.Modified, patch.Current)
			if err != nil {
				return nil, emperrors.Wrap(err, "Failed to generate merge patch")
			}

			return UpdateWithPatchAction(a.object, types.MergePatchType, jsonPatch), nil
		})))
	cmd.Bind(a.GetLastResult())
	return c.Do(ctx, cmd)
}

type createIfNotExistsAction struct {
	*BaseAction
	factoryFunc         func() (client.Object, error)
	newObject           client.Object
	createActionOptions []CreateActionOption
}

func CreateIfNotExistsFactoryItem(
	newObj client.Object,
	factoryFunc func() (client.Object, error),
	opts ...CreateActionOption,
) *createIfNotExistsAction {
	return &createIfNotExistsAction{
		newObject:           newObj,
		createActionOptions: opts,
		factoryFunc:         factoryFunc,
		BaseAction:          NewBaseAction("createIfNotExistsAction"),
	}
}

func (a *createIfNotExistsAction) Bind(result *ExecResult) {
	a.SetLastResult(result)
}

func (a *createIfNotExistsAction) Exec(ctx context.Context, c *ClientCommand) (*ExecResult, error) {
	reqLogger := a.GetReqLogger(c)

	result, err := a.factoryFunc()

	if err != nil {
		reqLogger.Error(err, "failure creating factory obj")
		return NewExecResult(Error, reconcile.Result{Requeue: true}, a.BaseAction, err), emperrors.Wrap(err, "error with create")
	}

	key := client.ObjectKeyFromObject(result)
	reqLogger = reqLogger.WithValues("requestType", fmt.Sprintf("%T", a.newObject), "key", key)

	reqLogger.V(0).Info("Creating object if not found", "object", result)
	return c.Do(
		ctx,
		HandleResult(
			GetAction(key, a.newObject),
			OnNotFound(
				HandleResult(
					CreateAction(result, a.createActionOptions...),
					OnRequeue(ContinueResponse()),
				),
			),
		),
	)
}
