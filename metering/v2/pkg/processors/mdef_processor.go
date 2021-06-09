package processors

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/dictionary"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/mailbox"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/meterdefinition"
	"k8s.io/client-go/tools/cache"
)

type MeterDefinitionRemovalWatcher struct {
	*Processor
	dictionary           *dictionary.MeterDefinitionDictionary
	meterDefinitionStore *meterdefinition.MeterDefinitionStore

	messageChan chan cache.Delta
	log         logr.Logger
}

func ProvideMeterDefinitionRemovalWatcher(
	dictionary *dictionary.MeterDefinitionDictionary,
	meterDefinitionStore *meterdefinition.MeterDefinitionStore,
	mb *mailbox.Mailbox,
	log logr.Logger,
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
	}
	sp.Processor.DeltaProcessor = sp
	return sp
}

func (w *MeterDefinitionRemovalWatcher) Process(ctx context.Context, d cache.Delta) error {
	if d.Type == cache.Deleted {
		return w.onDelete(ctx, d)
	}

	if d.Type == cache.Added {
		return w.onAdd(ctx, d)
	}

	return nil
}


func (w *MeterDefinitionRemovalWatcher) onAdd(ctx context.Context, d cache.Delta) error {
	meterdef, ok := d.Object.(*dictionary.MeterDefinitionExtended)

	if !ok {
		return errors.New("encountered unexpected type")
	}

	w.log.Info("processing meterdef", "name/namespace", meterdef.Name+"/"+meterdef.Namespace)
	return w.meterDefinitionStore.Resync()
}


func (w *MeterDefinitionRemovalWatcher) onDelete(ctx context.Context, d cache.Delta) error {
	meterdef, ok := d.Object.(*dictionary.MeterDefinitionExtended)

	if !ok {
		return errors.New("encountered unexpected type")
	}

	w.log.Info("processing meterdef", "name/namespace", meterdef.Name+"/"+meterdef.Namespace)

	key, err := cache.MetaNamespaceKeyFunc(&meterdef.MeterDefinition)

	if err != nil {
		w.log.Error(err, "error creating key")
		return errors.WithStack(err)
	}

	objects, err := w.meterDefinitionStore.ByIndex(meterdefinition.IndexMeterDefinition, key)

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
