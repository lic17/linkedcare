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
package app

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"linkedcare.io/linkedcare/pkg/controller/application"
	"linkedcare.io/linkedcare/pkg/controller/destinationrule"

	"time"

	"linkedcare.io/linkedcare/pkg/controller/virtualservice"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	istioclientset "istio.io/client-go/pkg/clientset/versioned"
	istioinformers "istio.io/client-go/pkg/informers/externalversions"
	linkedcareclientset "linkedcare.io/linkedcare/pkg/client/clientset/versioned"
	linkedcareinformers "linkedcare.io/linkedcare/pkg/client/informers/externalversions"
	applicationclientset "sigs.k8s.io/application/pkg/client/clientset/versioned"
	applicationinformers "sigs.k8s.io/application/pkg/client/informers/externalversions"
)

const defaultResync = 600 * time.Second

var log = logf.Log.WithName("controller-manager")

func AddControllers(mgr manager.Manager, cfg *rest.Config, stopCh <-chan struct{}) error {

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Error(err, "building kubernetes client failed")
	}

	istioclient, err := istioclientset.NewForConfig(cfg)
	if err != nil {
		log.Error(err, "create istio client failed")
		return err
	}

	applicationClient, err := applicationclientset.NewForConfig(cfg)
	if err != nil {
		log.Error(err, "create application client failed")
		return err
	}
	linkedcareclient, err := linkedcareclientset.NewForConfig(cfg)
	if err != nil {
		log.Error(err, "create linkedcare client failed")
		return err
	}

	informerFactory := informers.NewSharedInformerFactory(kubeClient, defaultResync)
	istioInformer := istioinformers.NewSharedInformerFactory(istioclient, defaultResync)
	applicationInformer := applicationinformers.NewSharedInformerFactory(applicationClient, defaultResync)

	linkedcareInformer := linkedcareinformers.NewSharedInformerFactory(linkedcareclient, defaultResync)

	vsController := virtualservice.NewVirtualServiceController(informerFactory.Core().V1().Services(),
		istioInformer.Networking().V1beta1().VirtualServices(),
		istioInformer.Networking().V1beta1().DestinationRules(),
		linkedcareInformer.Servicemesh().V1alpha1().Strategies(),
		kubeClient,
		istioclient,
		linkedcareclient)

	drController := destinationrule.NewDestinationRuleController(informerFactory.Apps().V1().Deployments(),
		istioInformer.Networking().V1beta1().DestinationRules(),
		informerFactory.Core().V1().Services(),
		linkedcareInformer.Servicemesh().V1alpha1().ServicePolicies(),
		kubeClient,
		istioclient,
		linkedcareclient)

	apController := application.NewApplicationController(informerFactory.Core().V1().Services(),
		informerFactory.Apps().V1().Deployments(),
		informerFactory.Apps().V1().StatefulSets(),
		linkedcareInformer.Servicemesh().V1alpha1().Strategies(),
		linkedcareInformer.Servicemesh().V1alpha1().ServicePolicies(),
		applicationInformer.App().V1beta1().Applications(),
		kubeClient,
		applicationClient)

	linkedcareInformer.Start(stopCh)
	istioInformer.Start(stopCh)
	informerFactory.Start(stopCh)
	applicationInformer.Start(stopCh)

	controllers := map[string]manager.Runnable{
		"virtualservice-controller":  vsController,
		"destinationrule-controller": drController,
		"application-controller":     apController,
	}

	for name, ctrl := range controllers {
		err = mgr.Add(ctrl)
		if err != nil {
			log.Error(err, "add controller to manager failed", "name", name)
			return err
		}
	}

	return nil
}
