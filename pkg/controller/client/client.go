package client

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	linkedcareVersioned "linkedcare.io/linkedcare/pkg/client/clientset/versioned"
)

type ClientSet struct {
	LinkedcareClient *linkedcareVersioned.Clientset
	K8sClient        *kubernetes.Clientset
	Config           *rest.Config
}

func CreateClientSet(kubeconfig string, stopper <-chan struct{}) *ClientSet {
	// creates the Client
	cfg := newConfig(kubeconfig)
	clientset, err := linkedcareVersioned.NewForConfig(cfg)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	k8sclientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return &ClientSet{
		LinkedcareClient: clientset,
		K8sClient:        k8sclientset,
		Config:           cfg,
	}
	/*factory := externalversions.NewSharedInformerFactory(clientset, 0)

	//
	//handle for servicePolicies
	//
	servicePolicieshd := newServicePoliciesHandler()

	servicePoliciesInformer := factory.Servicemesh().V1alpha1().ServicePolicies()
	servicePoliciesInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				servicePolicieshd.Add(obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				servicePolicieshd.Update(oldObj, newObj)
			},
			DeleteFunc: func(obj interface{}) {
				servicePolicieshd.Delete(obj)
			},
		})
	//
	//handle for Strategy
	//
	strategyhd := newStrategyHandler()

	strategiesInformer := factory.Servicemesh().V1alpha1().Strategies()
	strategiesInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				strategyhd.Add(obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				strategyhd.Update(oldObj, newObj)
			},
			DeleteFunc: func(obj interface{}) {
				strategyhd.Delete(obj)
			},
		})
	go factory.Start(stopper)

	factory.WaitForCacheSync(stopper)
	*/
}

func newConfig(kubeconfig string) (cfg *rest.Config) {
	var err error
	if len(kubeconfig) != 0 {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			klog.Fatal(err)
		}
	} else {
		if cfg, err = rest.InClusterConfig(); err != nil {
			klog.Fatal(err)
		}
	}
	return
}
