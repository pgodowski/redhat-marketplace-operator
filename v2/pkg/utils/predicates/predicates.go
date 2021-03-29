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

package predicates

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func NamespacePredicate(targetNamespace string) predicate.Funcs {
	return predicate.Funcs{
		// Ensures MarketPlaceConfig reconciles only within namespace
		GenericFunc: func(e event.GenericEvent) bool {
			return e.Meta.GetNamespace() == targetNamespace
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.MetaOld.GetNamespace() == targetNamespace && e.MetaNew.GetNamespace() == targetNamespace
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return e.Meta.GetNamespace() == targetNamespace
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return e.Meta.GetNamespace() == targetNamespace
		},
	}
}
