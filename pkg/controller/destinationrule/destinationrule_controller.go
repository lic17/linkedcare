package destinationrule

import (
	"context"
	"fmt"
	"reflect"

	apinetworkingv1beta1 "istio.io/api/networking/v1beta1"
	clientgonetworkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	log "k8s.io/klog"
	servicemeshv1alpha1 "linkedcare.io/linkedcare/pkg/apis/servicemesh/v1alpha1"
	"linkedcare.io/linkedcare/pkg/controller/virtualservice/util"

	"time"

	istioclientset "istio.io/client-go/pkg/clientset/versioned"
	istioinformers "istio.io/client-go/pkg/informers/externalversions/networking/v1beta1"
	istiolisters "istio.io/client-go/pkg/listers/networking/v1beta1"
	informersv1 "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	listersv1 "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	servicemeshclient "linkedcare.io/linkedcare/pkg/client/clientset/versioned"
	servicemeshinformers "linkedcare.io/linkedcare/pkg/client/informers/externalversions/servicemesh/v1alpha1"
	servicemeshlisters "linkedcare.io/linkedcare/pkg/client/listers/servicemesh/v1alpha1"
)

const (
	// maxRetries is the number of times a service will be retried before it is dropped out of the queue.
	// With the current rate-limiter in use (5ms*2^(maxRetries-1)) the following numbers represent the
	// sequence of delays between successive queuings of a service.
	//
	// 5ms, 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, 640ms, 1.3s, 2.6s, 5.1s, 10.2s, 20.4s, 41s, 82s
	maxRetries = 15
)

type DestinationRuleController struct {
	client clientset.Interface

	destinationRuleClient istioclientset.Interface
	servicemeshClient     servicemeshclient.Interface

	eventBroadcaster record.EventBroadcaster
	eventRecorder    record.EventRecorder

	serviceLister corelisters.ServiceLister
	serviceSynced cache.InformerSynced

	deploymentLister listersv1.DeploymentLister
	deploymentSynced cache.InformerSynced

	servicePolicyLister servicemeshlisters.ServicePolicyLister
	servicePolicySynced cache.InformerSynced

	destinationRuleLister istiolisters.DestinationRuleLister
	destinationRuleSynced cache.InformerSynced

	queue workqueue.RateLimitingInterface

	workerLoopPeriod time.Duration
}

func NewDestinationRuleController(deploymentInformer informersv1.DeploymentInformer,
	destinationRuleInformer istioinformers.DestinationRuleInformer,
	serviceInformer coreinformers.ServiceInformer,
	servicePolicyInformer servicemeshinformers.ServicePolicyInformer,
	client clientset.Interface,
	destinationRuleClient istioclientset.Interface,
	servicemeshClient servicemeshclient.Interface) *DestinationRuleController {

	broadcaster := record.NewBroadcaster()
	broadcaster.StartLogging(func(format string, args ...interface{}) {
		log.Info(fmt.Sprintf(format, args))
	})
	broadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: client.CoreV1().Events("")})
	recorder := broadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "destinationrule-controller"})

	v := &DestinationRuleController{
		client:                client,
		destinationRuleClient: destinationRuleClient,
		servicemeshClient:     servicemeshClient,
		queue:                 workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "destinationrule"),
		workerLoopPeriod:      time.Second,
	}

	v.deploymentLister = deploymentInformer.Lister()
	v.deploymentSynced = deploymentInformer.Informer().HasSynced

	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    v.addDeployment,
		DeleteFunc: v.deleteDeployment,
		UpdateFunc: func(old, cur interface{}) {
			v.addDeployment(cur)
		},
	})

	v.serviceLister = serviceInformer.Lister()
	v.serviceSynced = serviceInformer.Informer().HasSynced

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    v.enqueueService,
		DeleteFunc: v.enqueueService,
		UpdateFunc: func(old, cur interface{}) {
			v.enqueueService(cur)
		},
	})

	v.destinationRuleLister = destinationRuleInformer.Lister()
	v.destinationRuleSynced = destinationRuleInformer.Informer().HasSynced

	v.servicePolicyLister = servicePolicyInformer.Lister()
	v.servicePolicySynced = servicePolicyInformer.Informer().HasSynced

	servicePolicyInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: v.addServicePolicy,
		UpdateFunc: func(old, cur interface{}) {
			v.addServicePolicy(cur)
		},
		DeleteFunc: v.addServicePolicy,
	})

	v.eventBroadcaster = broadcaster
	v.eventRecorder = recorder

	return v

}

func (v *DestinationRuleController) Start(stopCh <-chan struct{}) error {
	return v.Run(5, stopCh)
}

func (v *DestinationRuleController) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer v.queue.ShutDown()

	log.Info("starting destinationrule controller")
	defer log.Info("shutting down destinationrule controller")

	if !cache.WaitForCacheSync(stopCh, v.serviceSynced, v.destinationRuleSynced, v.deploymentSynced, v.servicePolicySynced) {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	for i := 0; i < workers; i++ {
		go wait.Until(v.worker, v.workerLoopPeriod, stopCh)
	}

	<-stopCh
	return nil
}

func (v *DestinationRuleController) enqueueService(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}

	v.queue.Add(key)
}

func (v *DestinationRuleController) worker() {
	for v.processNextWorkItem() {

	}
}

func (v *DestinationRuleController) processNextWorkItem() bool {
	eKey, quit := v.queue.Get()
	if quit {
		return false
	}

	defer v.queue.Done(eKey)

	err := v.syncService(eKey.(string))
	v.handleErr(err, eKey)

	return true
}

// main function of the reconcile for destinationrule
// destinationrule's name is same with the service that created it
func (v *DestinationRuleController) syncService(key string) error {
	startTime := time.Now()
	defer func() {
		log.V(4).Info("Finished syncing service destinationrule.", "key", key, "duration", time.Since(startTime))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	service, err := v.serviceLister.Services(namespace).Get(name)
	if err != nil {
		// delete the corresponding destinationrule if there is any, as the service has been deleted.
		err = v.destinationRuleClient.NetworkingV1beta1().DestinationRules(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			log.Error(err, "delete destination rule failed", "namespace", namespace, "name", name)
			return err
		}

		// delete orphan service policy if there is any
		err = v.servicemeshClient.ServicemeshV1alpha1().ServicePolicies(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			log.Error(err, "delete orphan service policy failed", "namespace", namespace, "name", name)
			return err
		}

		return nil
	}

	if len(service.Labels) < len(util.ApplicationLabels) ||
		!util.IsApplicationComponent(service.Labels) ||
		!util.IsServicemeshEnabled(service.Annotations) ||
		len(service.Spec.Ports) == 0 {
		// services don't have enough labels to create a virtualservice
		// or they don't have necessary labels
		// or they don't have servicemesh enabled
		// or they don't have any ports defined
		return nil
	}

	appName := util.GetComponentName(&service.ObjectMeta)

	// fetch all deployments that match with service selector
	deployments, err := v.deploymentLister.Deployments(namespace).List(labels.Set(service.Spec.Selector).AsSelectorPreValidated())
	if err != nil {
		return err
	}

	subsets := make([]*apinetworkingv1beta1.Subset, 0)
	for _, deployment := range deployments {

		// not a valid deployment we required
		if !util.IsApplicationComponent(deployment.Labels) ||
			!util.IsApplicationComponent(deployment.Spec.Selector.MatchLabels) ||
			deployment.Status.ReadyReplicas == 0 ||
			!util.IsServicemeshEnabled(deployment.Annotations) {
			continue
		}

		version := util.GetComponentVersion(&deployment.ObjectMeta)

		if len(version) == 0 {
			log.V(4).Info("Deployment doesn't have a version label", "key", types.NamespacedName{Namespace: deployment.Namespace, Name: deployment.Name}.String())
			continue
		}

		subset := &apinetworkingv1beta1.Subset{
			Name: util.NormalizeVersionName(version),
			Labels: map[string]string{
				util.VersionLabel: version,
			},
		}

		subsets = append(subsets, subset)
	}

	currentDestinationRule, err := v.destinationRuleLister.DestinationRules(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			currentDestinationRule = &clientgonetworkingv1beta1.DestinationRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:   service.Name,
					Labels: service.Labels,
				},
				Spec: apinetworkingv1beta1.DestinationRule{
					Host: name,
				},
			}
		} else {
			log.Error(err, "Couldn't get destinationrule for service", "key", key)
			return err
		}
	}

	// fetch all servicepolicies associated to this service
	servicePolicies, err := v.servicePolicyLister.ServicePolicies(namespace).List(labels.SelectorFromSet(map[string]string{util.AppLabel: appName}))
	if err != nil {
		log.Error(err, "could not list service policies is namespace with component name", "namespace", namespace, "name", appName)
		return err
	}

	dr := currentDestinationRule.DeepCopy()
	dr.Spec.TrafficPolicy = nil
	dr.Spec.Subsets = subsets
	//
	if len(servicePolicies) > 0 {
		if len(servicePolicies) > 1 {
			err = fmt.Errorf("more than one service policy associated with service %s/%s is forbidden", namespace, name)
			log.Error(err, "")
			return err
		}

		sp := servicePolicies[0]
		if sp.Spec.Template.Spec.TrafficPolicy != nil {
			dr.Spec.TrafficPolicy = sp.Spec.Template.Spec.TrafficPolicy
		}

		for _, subset := range sp.Spec.Template.Spec.Subsets {
			for i := range dr.Spec.Subsets {
				if subset.Name == dr.Spec.Subsets[i].Name && subset.TrafficPolicy != nil {
					dr.Spec.Subsets[i].TrafficPolicy = subset.TrafficPolicy
				}
			}
		}
	}

	createDestinationRule := len(currentDestinationRule.ResourceVersion) == 0

	if !createDestinationRule && reflect.DeepEqual(currentDestinationRule.Spec, dr.Spec) &&
		reflect.DeepEqual(currentDestinationRule.Labels, service.Labels) {
		log.V(5).Info("destinationrule are equal, skipping update", "key", types.NamespacedName{Namespace: service.Namespace, Name: service.Name}.String())
		return nil
	}

	newDestinationRule := currentDestinationRule.DeepCopy()
	newDestinationRule.Spec = dr.Spec
	newDestinationRule.Labels = service.Labels
	if newDestinationRule.Annotations == nil {
		newDestinationRule.Annotations = make(map[string]string)
	}

	if createDestinationRule {
		_, err = v.destinationRuleClient.NetworkingV1beta1().DestinationRules(namespace).Create(context.TODO(), newDestinationRule, metav1.CreateOptions{})
	} else {
		_, err = v.destinationRuleClient.NetworkingV1beta1().DestinationRules(namespace).Update(context.TODO(), newDestinationRule, metav1.UpdateOptions{})
	}

	if err != nil {
		if createDestinationRule && errors.IsForbidden(err) {
			// A request is forbidden primarily for two reasons:
			// 1. namespace is terminating, endpoint creation is not allowed by default.
			// 2. policy is misconfigured, in which case no service would function anywhere.
			// Given the frequency of 1, we log at a lower level.
			log.V(5).Info("Forbidden from creating endpoints", "error", err)
		}

		if createDestinationRule {
			v.eventRecorder.Event(newDestinationRule, v1.EventTypeWarning, "FailedToCreateDestinationRule", fmt.Sprintf("Failed to create destinationrule for service %v/%v: %v", service.Namespace, service.Name, err))
		} else {
			v.eventRecorder.Event(newDestinationRule, v1.EventTypeWarning, "FailedToUpdateDestinationRule", fmt.Sprintf("Failed to update destinationrule for service %v/%v: %v", service.Namespace, service.Name, err))
		}

		return err
	}

	return nil
}

// When a destinationrule is added, figure out which service it will be used
// and enqueue it. obj must have *appsv1.Deployment type
func (v *DestinationRuleController) addDeployment(obj interface{}) {
	deploy := obj.(*appsv1.Deployment)

	// not a application component
	if !util.IsApplicationComponent(deploy.Labels) || !util.IsApplicationComponent(deploy.Spec.Selector.MatchLabels) {
		return
	}

	services, err := v.getDeploymentServiceMemberShip(deploy)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("unable to get deployment %s/%s's service memberships", deploy.Namespace, deploy.Name))
		return
	}

	for key := range services {
		v.queue.Add(key)
	}

	return
}

func (v *DestinationRuleController) deleteDeployment(obj interface{}) {
	if _, ok := obj.(*appsv1.Deployment); ok {
		v.addDeployment(obj)
		return
	}

	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
		return
	}

	deploy, ok := tombstone.Obj.(*appsv1.Deployment)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a deployment %#v", obj))
		return
	}

	v.addDeployment(deploy)
}

func (v *DestinationRuleController) getDeploymentServiceMemberShip(deployment *appsv1.Deployment) (sets.String, error) {
	set := sets.String{}

	allServices, err := v.serviceLister.Services(deployment.Namespace).List(labels.Everything())
	if err != nil {
		return set, err
	}

	for i := range allServices {
		service := allServices[i]
		if service.Spec.Selector == nil ||
			!util.IsApplicationComponent(service.Labels) ||
			!util.IsServicemeshEnabled(service.Annotations) {
			// services with nil selectors match nothing, not everything.
			continue
		}
		selector := labels.Set(service.Spec.Selector).AsSelectorPreValidated()
		if selector.Matches(labels.Set(deployment.Spec.Selector.MatchLabels)) {
			key, err := cache.MetaNamespaceKeyFunc(service)
			if err != nil {
				return nil, err
			}
			set.Insert(key)
		}
	}

	return set, nil
}

func (v *DestinationRuleController) addServicePolicy(obj interface{}) {
	servicePolicy := obj.(*servicemeshv1alpha1.ServicePolicy)

	appName := servicePolicy.Labels[util.AppLabel]

	services, err := v.serviceLister.Services(servicePolicy.Namespace).List(labels.SelectorFromSet(map[string]string{util.AppLabel: appName}))
	if err != nil {
		log.Error(err, "cannot list services", "namespace", servicePolicy.Namespace, "name", appName)
		utilruntime.HandleError(fmt.Errorf("cannot list services in namespace %s, with component name %v", servicePolicy.Namespace, appName))
		return
	}

	set := sets.String{}
	for _, service := range services {
		key, err := cache.MetaNamespaceKeyFunc(service)
		if err != nil {
			utilruntime.HandleError(err)
			continue
		}
		set.Insert(key)
	}

	// avoid enqueue a key multiple times
	for key := range set {
		v.queue.Add(key)
	}
}

func (v *DestinationRuleController) handleErr(err error, key interface{}) {
	if err == nil {
		v.queue.Forget(key)
		return
	}

	if v.queue.NumRequeues(key) < maxRetries {
		log.V(2).Info("Error syncing virtualservice for service, retrying.", "key", key, "error", err)
		v.queue.AddRateLimited(key)
		return
	}

	log.V(4).Info("Dropping service out of the queue", "key", key, "error", err)
	v.queue.Forget(key)
	utilruntime.HandleError(err)
}