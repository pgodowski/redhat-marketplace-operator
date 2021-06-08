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

package dictionary

import (
	"context"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	"github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/filter"
	pkgtypes "github.com/redhat-marketplace/redhat-marketplace-operator/metering/v2/pkg/types"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/apis/marketplace/v1beta1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/client"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/managers"
	"github.com/sasha-s/go-deadlock"
	"golang.org/x/time/rate"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ cache.Store = &MeterDefinitionDictionary{}
var _ metav1.Object = &MeterDefinitionExtended{}

type MeterDefinitionExtended struct {
	v1beta1.MeterDefinition
	Filter filter.MeterDefinitionLookupFilter
}

var _ metav1.Object = &MeterDefinitionExtended{}

type MeterDefinitionDictionary struct {
	cache       cache.Store
	keyFunc     cache.KeyFunc
	delta       *cache.DeltaFIFO
	starterList *v1beta1.MeterDefinitionList

	meterDefinitionsSeen MeterDefinitionsSeenStore

	findOwner *client.FindOwnerHelper

	log logr.Logger

	rateLimits map[types.UID]*rate.Limiter
	deadlock.RWMutex
}

func NewMeterDefinitionDictionary(
	ctx context.Context,
	kubeClient clientset.Interface,
	findOwner *client.FindOwnerHelper,
	namespaces pkgtypes.Namespaces,
	log logr.Logger,
	list *v1beta1.MeterDefinitionList,
	meterDefinitionsSeen MeterDefinitionsSeenStore,
) *MeterDefinitionDictionary {
	return &MeterDefinitionDictionary{
		log:                  log.WithName("mdef_dictionary"),
		keyFunc:              cache.MetaNamespaceKeyFunc,
		cache:                cache.NewStore(cache.MetaNamespaceKeyFunc),
		delta:                cache.NewDeltaFIFO(cache.MetaNamespaceKeyFunc, nil),
		findOwner:            findOwner,
		rateLimits:           map[types.UID]*rate.Limiter{},
		starterList:          list,
		meterDefinitionsSeen: meterDefinitionsSeen,
	}
}

func (w *MeterDefinitionDictionary) Start(ctx context.Context) error {
	// for i := range w.starterList.Items {
	// 	localMdef := w.starterList.Items[i]
	// 	err := w.Add(&localMdef)

	// 	if err != nil {
	// 		return errors.WithStack(err)
	// 	}
	// }

	return nil
}

func (def *MeterDefinitionDictionary) FindObjectMatches(
	obj interface{},
	results *[]filter.Result,
	skipCache bool,
) error {
	filters, err := def.ListFilters()

	if err != nil {
		def.log.Error(err, "error listing filters")
		return err
	}

	for i := range filters {
		localLookup := &filters[i]

		var (
			ok *bool
		)

		if !skipCache {
			ok = lookupCache.Get(localLookup, obj)
		}

		if ok == nil {
			lookupOk, err := localLookup.Matches(obj)

			if err != nil {
				def.log.Error(err, "error finding matches")
				return err
			}

			err = lookupCache.Set(localLookup, obj, lookupOk)
			if err != nil {
				def.log.Error(err, "error saving cache")
				return err
			}

			ok = &lookupOk
		}

		if ok != nil && *ok {
			*results = append(*results, filter.Result{
				Ok:     *ok,
				Lookup: localLookup,
			})
		}
	}

	return nil
}

func (def *MeterDefinitionDictionary) ListFilters() ([]filter.MeterDefinitionLookupFilter, error) {
	def.RLock()
	defer def.RUnlock()

	filters := []filter.MeterDefinitionLookupFilter{}

	objs := def.cache.List()

	for _, obj := range objs {
		meterdef, ok := obj.(*MeterDefinitionExtended)

		if !ok {
			return nil, errors.New("encountered unexpected type")
		}

		filters = append(filters, meterdef.Filter)
	}

	return filters, nil
}

// Add adds the given object to the accumulator associated with the given object's key
func (def *MeterDefinitionDictionary) Add(obj interface{}) error {
	def.Lock()
	defer def.Unlock()

	addObj, err := def.newMeterDefinitionExtended(obj)

	if err != nil {
		def.log.Error(err, "error exending obj")
		return err
	}

	if !def.allow(addObj) {
		def.log.Info("rate limited, skipping")
		return nil
	}

	def.log.Info("recording obj", "obj", fmt.Sprintf("%+v", obj))

	if err := def.delta.Add(addObj); err != nil {
		def.log.Error(err, "error adding obj to delta")
		return err
	}

	if err := def.cache.Add(addObj); err != nil {
		def.log.Error(err, "error adding obj to cache")
		return err
	}

	return nil
}

// Update updates the given object in the accumulator associated with the given object's key
func (def *MeterDefinitionDictionary) Update(obj interface{}) error {
	def.Lock()
	defer def.Unlock()

	addObj, err := def.newMeterDefinitionExtended(obj)

	if err != nil {
		def.log.Error(err, "error exending obj")
		return err
	}

	if !def.allow(addObj) {
		def.log.Info("rate limited, skipping")
		return nil
	}

	def.log.Info("updating obj", "obj", fmt.Sprintf("%+v", obj))

	if err := def.delta.Update(addObj); err != nil {
		def.log.Error(err, "error adding obj to delta")
		return err
	}

	if err := def.cache.Update(addObj); err != nil {
		def.log.Error(err, "error adding obj to cache")
		return err
	}

	return nil
}

const (
	// 10 every 5 seconds
	limitRateBucket = 10
	limitRate       = 5 * time.Second
)

func (def *MeterDefinitionDictionary) allow(addObj metav1.Object) bool {
	var (
		rateLimiter *rate.Limiter
		ok          bool
	)

	if rateLimiter, ok = def.rateLimits[addObj.GetUID()]; !ok {
		limit := rate.Every(limitRate)
		rateLimiter = rate.NewLimiter(limit, limitRateBucket)
		def.rateLimits[addObj.GetUID()] = rateLimiter
	}

	return rateLimiter.Allow()
}

// Delete deletes the given object from the accumulator associated with the given object's key
func (def *MeterDefinitionDictionary) Delete(obj interface{}) error {
	def.Lock()
	defer def.Unlock()

	addObj, err := def.newMeterDefinitionExtended(obj)

	if err != nil {
		def.log.Error(err, "error exending obj")
		return err
	}

	if err := def.delta.Delete(addObj); err != nil {
		return err
	}

	o, err := meta.Accessor(obj)
	if err != nil {
		def.log.Error(err, "error converting obj")
		return err
	}

	if _, ok := def.rateLimits[o.GetUID()]; ok {
		delete(def.rateLimits, o.GetUID())
	}

	return def.cache.Delete(addObj)
}

// List returns a list of all the currently non-empty accumulators
func (def *MeterDefinitionDictionary) List() []interface{} {
	def.RLock()
	defer def.RUnlock()

	return def.cache.List()
}

// ListKeys returns a list of all the keys currently associated with non-empty accumulators
func (def *MeterDefinitionDictionary) ListKeys() []string {
	def.RLock()
	defer def.RUnlock()

	return def.cache.ListKeys()
}

// Get returns the accumulator associated with the given object's key
func (def *MeterDefinitionDictionary) Get(obj interface{}) (item interface{}, exists bool, err error) {
	def.RLock()
	defer def.RUnlock()

	return def.cache.Get(obj)
}

// GetByKey returns the accumulator associated with the given key
func (def *MeterDefinitionDictionary) GetByKey(key string) (item interface{}, exists bool, err error) {
	def.RLock()
	defer def.RUnlock()

	return def.cache.GetByKey(key)
}

func (def *MeterDefinitionDictionary) AddIfNotPresent(obj interface{}) error {
	def.Lock()
	defer def.Unlock()

	return def.delta.AddIfNotPresent(obj)
}

func (def *MeterDefinitionDictionary) Close() {
	def.delta.Close()
}

func (def *MeterDefinitionDictionary) HasSynced() bool {
	return def.delta.HasSynced()
}

func (def *MeterDefinitionDictionary) Pop(process cache.PopProcessFunc) (interface{}, error) {
	return def.delta.Pop(process)
}

// Replace will delete the contents of the store, using instead the
// given list. Store takes ownership of the list, you should not reference
// it after calling this function.
func (def *MeterDefinitionDictionary) Replace(in []interface{}, str string) error {
	def.Lock()
	defer def.Unlock()

	objs := []interface{}{}

	for _, obj := range in {
		localObj := obj
		addObj, err := def.newMeterDefinitionExtended(localObj)
		if err != nil {
			return err
		}

		objs = append(objs, addObj)
	}

	if err := def.cache.Replace(objs, str); err != nil {
		return err
	}

	if err := def.delta.Replace(objs, str); err != nil {
		return err
	}

	return nil
}

// Resync is meaningless in the terms appearing here but has
// meaning in some implementations that have non-trivial
// additional behavior (e.g., DeltaFIFO).
func (def *MeterDefinitionDictionary) Resync() error {
	list := def.meterDefinitionsSeen.List()
	for i := range list {
		obj := list[i]
		_, exists, err := def.Get(obj)

		if err != nil {
			return err
		}

		if !exists {
			if err := def.Add(obj); err != nil {
				return err
			}
		}
	}

	objects := def.List()

	for i := range objects {
		enobj := objects[i].(*MeterDefinitionExtended)
		_, exists, err := def.meterDefinitionsSeen.Get(&enobj.MeterDefinition)

		if err != nil {
			return err
		}

		if !exists {
			err := def.Delete(&enobj.MeterDefinition)
			if err != nil {
				return err
			}
		}
	}

	return def.delta.Resync()
}

func (def *MeterDefinitionDictionary) newMeterDefinitionExtended(obj interface{}) (*MeterDefinitionExtended, error) {
	meterdef, ok := obj.(*v1beta1.MeterDefinition)

	if !ok {
		return nil, errors.New("expected meter definition")
	}

	lookup, err := filter.NewMeterDefinitionLookupFilter(meterdef, def.findOwner, def.log)

	if err != nil {
		return nil, err
	}

	return &MeterDefinitionExtended{
		MeterDefinition: *meterdef,
		Filter:          *lookup,
	}, nil
}

func ProvideMeterDefinitionList(
	cache managers.CacheIsStarted,
	client runtimeClient.Client,
) (*v1beta1.MeterDefinitionList, error) {
	obj := v1beta1.MeterDefinitionList{}
	err := client.List(context.TODO(), &obj, &runtimeClient.ListOptions{
		Namespace: "",
	})
	return &obj, err
}

type MeterDefinitionsSeenStore cache.Store

func NewMeterDefinitionsSeenStore() MeterDefinitionsSeenStore {
	return cache.NewStore(cache.MetaNamespaceKeyFunc)
}
