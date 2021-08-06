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

package razeeengine

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func CreateNodeListWatch(kubeClient clientset.Interface) func(string) cache.ListerWatcher {
	return func(ns string) cache.ListerWatcher {
		return &cache.ListWatch{
			ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
				return kubeClient.CoreV1().Nodes().List(context.TODO(), opts)
			},
			WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
				return kubeClient.CoreV1().Nodes().Watch(context.TODO(), opts)
			},
		}
	}
}

func CreateServiceListWatch(kubeClient clientset.Interface) func(string) cache.ListerWatcher {
	return func(ns string) cache.ListerWatcher {
		return &cache.ListWatch{
			ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
				return kubeClient.CoreV1().Services(ns).List(context.TODO(), opts)
			},
			WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
				return kubeClient.CoreV1().Services(ns).Watch(context.TODO(), opts)
			},
		}
	}
}
