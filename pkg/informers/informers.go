/*

 Copyright 2019 The Linkedcare Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

*/
package informers

import (
	"sync"
	"time"

	k8sinformers "k8s.io/client-go/informers"
	lcinformers "linkedcare.io/linkedcare/pkg/client/informers/externalversions"
	"linkedcare.io/linkedcare/pkg/simple/client"
	applicationinformers "sigs.k8s.io/application/pkg/client/informers/externalversions"
)

const defaultResync = 600 * time.Second

var (
	k8sOnce            sync.Once
	lcOnce             sync.Once
	appOnce            sync.Once
	informerFactory    k8sinformers.SharedInformerFactory
	lcInformerFactory  lcinformers.SharedInformerFactory
	appInformerFactory applicationinformers.SharedInformerFactory
)

func SharedInformerFactory() k8sinformers.SharedInformerFactory {
	k8sOnce.Do(func() {
		k8sClient := client.ClientSets().K8s().Kubernetes()
		informerFactory = k8sinformers.NewSharedInformerFactory(k8sClient, defaultResync)
	})
	return informerFactory
}

func LcSharedInformerFactory() lcinformers.SharedInformerFactory {
	lcOnce.Do(func() {
		k8sClient := client.ClientSets().K8s().Linkedcare()
		lcInformerFactory = lcinformers.NewSharedInformerFactory(k8sClient, defaultResync)
	})
	return lcInformerFactory
}

func AppSharedInformerFactory() applicationinformers.SharedInformerFactory {
	appOnce.Do(func() {
		appClient := client.ClientSets().K8s().Application()
		appInformerFactory = applicationinformers.NewSharedInformerFactory(appClient, defaultResync)
	})
	return appInformerFactory
}
