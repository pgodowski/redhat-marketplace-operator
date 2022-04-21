// Copyright 2021 IBM Corp.
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

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/filter"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/mailbox"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/stores"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/common"
	marketplacev1beta1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MeterDefinitionRemovalWatcher struct {
	*Processor
	dictionary           *stores.MeterDefinitionDictionary
	meterDefinitionStore *stores.MeterDefinitionStore
	keyFunc              cache.KeyFunc

	nsWatcher   *filter.NamespaceWatcher
	messageChan chan cache.Delta
	log         logr.Logger
	kubeClient  client.Client
}

func ProvideMeterDefinitionRemovalWatcher(
	dictionary *stores.MeterDefinitionDictionary,
	meterDefinitionStore *stores.MeterDefinitionStore,
	mb *mailbox.Mailbox,
	log logr.Logger,
	kubeClient client.Client,
	nsWatcher *filter.NamespaceWatcher,
) *MeterDefinitionRemovalWatcher {
	sp := &MeterDefinitionRemovalWatcher{
		Processor: &Processor{
			log:           log,
			digestersSize: 1,
			retryCount:    3,
			mailbox:       mb,
			channelName:   mailbox.MeterDefinitionChannel,
		},
		dictionary:           dictionary,
		meterDefinitionStore: meterDefinitionStore,
		messageChan:          make(chan cache.Delta),
		log:                  log.WithName("meterDefinitionRemovalWatcher"),
		kubeClient:           kubeClient,
		nsWatcher:            nsWatcher,
	}
	sp.Processor.DeltaProcessor = sp
	return sp
}

func (w *MeterDefinitionRemovalWatcher) Process(ctx context.Context, d cache.Delta) error {
	meterdef, ok := d.Object.(*stores.MeterDefinitionExtended)
	if !ok {
		return errors.New("encountered unexpected type")
	}

	if d.Type == cache.Deleted {
		w.nsWatcher.RemoveNamespace(client.ObjectKeyFromObject(&meterdef.MeterDefinition))
		return w.onDelete(ctx, d)
	}

	if d.Type == cache.Added {
		w.nsWatcher.AddNamespace(client.ObjectKeyFromObject(&meterdef.MeterDefinition))
		return w.onAdd(ctx, d)
	}

	if d.Type == cache.Updated {
		w.nsWatcher.AddNamespace(client.ObjectKeyFromObject(&meterdef.MeterDefinition))
		return w.onUpdate(ctx, d)
	}

	if d.Type == cache.Sync {
		return w.onSync(ctx, d)
	}

	return nil
}

func (w *MeterDefinitionRemovalWatcher) onAdd(ctx context.Context, d cache.Delta) error {
	meterdef, ok := d.Object.(*stores.MeterDefinitionExtended)

	if !ok {
		return errors.New("encountered unexpected type")
	}

	w.log.Info("processing meterdef", "name/namespace", meterdef.Name+"/"+meterdef.Namespace)
	return w.meterDefinitionStore.Resync()
}

func (w *MeterDefinitionRemovalWatcher) onDelete(ctx context.Context, d cache.Delta) error {
	meterdef, ok := d.Object.(*stores.MeterDefinitionExtended)

	if !ok {
		return errors.New("encountered unexpected type")
	}

	w.log.Info("processing meterdef", "name/namespace", meterdef.Name+"/"+meterdef.Namespace)

	key, err := cache.MetaNamespaceKeyFunc(&meterdef.MeterDefinition)

	if err != nil {
		w.log.Error(err, "error creating key")
		return errors.WithStack(err)
	}

	objects, err := w.meterDefinitionStore.ByIndex(stores.IndexMeterDefinition, key)

	if err != nil {
		w.log.Error(err, "error finding data")
		return errors.WithStack(err)
	}

	for i := range objects {
		obj := objects[i]
		w.log.Info("deleting obj", "object", obj)
		err := w.meterDefinitionStore.Delete(obj)

		if err != nil {
			w.log.Error(err, "error deleting data")
			return errors.WithStack(err)
		}
	}

	return nil
}

func (w *MeterDefinitionRemovalWatcher) onUpdate(ctx context.Context, d cache.Delta) error {
	meterdef, ok := d.Object.(*stores.MeterDefinitionExtended)

	if !ok {
		return errors.New("encountered unexpected type")
	}

	w.log.Info("processing meterdef", "name/namespace", meterdef.Name+"/"+meterdef.Namespace)

	key, err := cache.MetaNamespaceKeyFunc(&meterdef.MeterDefinition)

	if err != nil {
		w.log.Error(err, "error creating key")
		return errors.WithStack(err)
	}

	objects, err := w.meterDefinitionStore.ByIndex(stores.IndexMeterDefinition, key)

	if err != nil {
		w.log.Error(err, "error finding data")
		return errors.WithStack(err)
	}

	for i := range objects {
		obj := objects[i]
		err := w.meterDefinitionStore.DeleteFromIndex(obj)

		if err != nil {
			w.log.Error(err, "error deleting data")
			return errors.WithStack(err)
		}
	}

	return w.meterDefinitionStore.Resync()
}

// Clear Status.WorkloadResources on initial sync such that Status is correct when metric-state starts
func (w *MeterDefinitionRemovalWatcher) onSync(ctx context.Context, d cache.Delta) error {
	meterdef, ok := d.Object.(*stores.MeterDefinitionExtended)

	if !ok {
		return errors.New("encountered unexpected type")
	}

	resources := []common.WorkloadResource{}
	mdef := &marketplacev1beta1.MeterDefinition{}

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if err := w.kubeClient.Get(ctx, types.NamespacedName{Name: meterdef.Name, Namespace: meterdef.Namespace}, mdef); err != nil {
			return err
		}

		mdef.Status.WorkloadResources = resources

		return w.kubeClient.Status().Update(ctx, mdef)
	})
	if err != nil {
		w.log.Error(err, "meterdefinition status update failed")
	}

	return nil
}
