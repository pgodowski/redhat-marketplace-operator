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

package processors

import (
	"context"
	"sort"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/mailbox"
	pkgtypes "github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/types"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/common"
	marketplacev1beta1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	"github.com/sasha-s/go-deadlock"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatusProcessor will update the meter definition
// status with the objects that matched it.
type StatusProcessor struct {
	*Processor
	log        logr.Logger
	kubeClient client.Client
	mutex      deadlock.Mutex
	scheme     *runtime.Scheme
}

// NewStatusProcessor is the provider that creates
// the processor.
func ProvideStatusProcessor(
	log logr.Logger,
	kubeClient client.Client,
	mb *mailbox.Mailbox,
	scheme *runtime.Scheme,
) *StatusProcessor {
	sp := &StatusProcessor{
		Processor: &Processor{
			log:           log,
			digestersSize: 1,
			retryCount:    3,
			mailbox:       mb,
			channelName:   mailbox.ObjectChannel,
		},
		log:        log.WithValues("process", "statusProcessor"),
		kubeClient: kubeClient,
		scheme:     scheme,
	}

	sp.Processor.DeltaProcessor = sp
	return sp
}

// Process will receive a new ObjectResourceMessage and find and update the metere
// definition associated with the object. To prevent gaps, it bulk retrieves the
// resources and checks it against the status.
func (u *StatusProcessor) Process(ctx context.Context, inObj cache.Delta) error {
	enhancedObj, ok := inObj.Object.(*pkgtypes.MeterDefinitionEnhancedObject)

	if !ok {
		return errors.WithStack(errors.New("obj is unexpected type"))
	}

	u.log.Info("updating status",
		"obj", enhancedObj.GetName()+"/"+enhancedObj.GetNamespace(),
		"mdefs", len(enhancedObj.MeterDefinitions),
	)

	for i := range enhancedObj.MeterDefinitions {
		key, err := client.ObjectKeyFromObject(&enhancedObj.MeterDefinitions[i])

		if err != nil {
			return errors.WithStack(err)
		}

		err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			u.mutex.Lock()
			defer u.mutex.Unlock()

			resources := []common.WorkloadResource{}
			set := map[types.UID]common.WorkloadResource{}

			mdef := &marketplacev1beta1.MeterDefinition{}

			if err := u.kubeClient.Get(ctx, key, mdef); err != nil {
				return err
			}

			for _, obj := range mdef.Status.WorkloadResources {
				set[obj.UID] = obj
			}

			workload, err := common.NewWorkloadResource(enhancedObj.Object, u.scheme)

			if err != nil {
				return errors.WithStack(err)
			}

			switch inObj.Type {
			case cache.Deleted:
				delete(set, workload.UID)
			case cache.Added:
				fallthrough
			case cache.Sync:
				fallthrough
			case cache.Updated:
				fallthrough
			case cache.Replaced:
				set[workload.UID] = *workload
			default:
				return nil
			}

			for i := range set {
				resources = append(resources, set[i])
			}

			sort.Sort(common.ByAlphabetical(resources))
			mdef.Status.WorkloadResources = resources

			return u.kubeClient.Status().Update(ctx, mdef)
		})

		if err != nil {
			if k8serrors.IsNotFound(err) {
				return nil
			}

			return errors.WithStack(err)
		}
	}

	return nil
}
