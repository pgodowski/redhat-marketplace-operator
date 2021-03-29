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

package runnables

import (
	"github.com/google/wire"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Runnables []manager.Runnable

var RunnableSet = wire.NewSet(
	ProvideRunnables,
	wire.Struct(new(CRDUpdater), "*"),
)

func ProvideRunnables(
	crdUpdater *CRDUpdater,
) Runnables {
	return []manager.Runnable{
		crdUpdater,
	}
}
