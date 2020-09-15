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

package servicemesh

import (
	"context"

	"errors"

	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog"
	"linkedcare.io/linkedcare/pkg/apiserver/applications"
	"linkedcare.io/linkedcare/pkg/simple/client"

	"linkedcare.io/linkedcare/pkg/informers"

	servicemeshv1alpha1 "linkedcare.io/linkedcare/pkg/apis/servicemesh/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServicePolicy struct {
	Component   string                                          `json:"component"`
	Application applications.Application                        `json:"application"`
	Template    servicemeshv1alpha1.DestinationRuleSpecTemplate `json:"template,omitempty"`
}

var servicepolicyTemplates servicemeshv1alpha1.ServicePolicy

// Load yamls
func init() {

	servicepolicyLabes := map[string]string{
		"app":                       "",
		"app.linkedcare.io/name":    "",
		"app.linkedcare.io/version": "",
	}

	selector := &metav1.LabelSelector{
		MatchLabels: servicepolicyLabes,
	}

	servicepolicyTemplates.Labels = servicepolicyLabes
	servicepolicyTemplates.Spec.Selector = selector
}

// Get all servicepolicys from a namespace
func GetAllServicepolicies(namespace string) ([]*servicemeshv1alpha1.ServicePolicy, error) {

	spLister := informers.LcSharedInformerFactory().Servicemesh().V1alpha1().ServicePolicies().Lister()
	sps, err := spLister.ServicePolicies(namespace).List(labels.Everything())

	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return sps, nil
}

// Get an servicepolicy from a namespace
func GetServicepolicy(namespace, name string) (*servicemeshv1alpha1.ServicePolicy, error) {
	app, err := getServicepolicy(namespace, name)
	return app, err
}

func getServicepolicy(namespace, name string) (*servicemeshv1alpha1.ServicePolicy, error) {
	spLister := informers.LcSharedInformerFactory().Servicemesh().V1alpha1().ServicePolicies().Lister()
	sp, err := spLister.ServicePolicies(namespace).Get(name)

	if err != nil {
		if k8serr.IsNotFound(err) {
			return nil, k8serr.NewNotFound(servicemeshv1alpha1.Resource("servicepolicys"), name)
		}
		klog.Error(err)
		return nil, err
	}
	return sp, nil
}

// Create an servicepolicy in a namespace
func CreateServicepolicy(namespace string, servicePolicy ServicePolicy) (*servicemeshv1alpha1.ServicePolicy, error) {

	app, err := createServicepolicy(namespace, servicePolicy)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return app, nil
}

// DeleteServicepolicy is used to delete the servicepolicy in namespace
func DeleteServicepolicy(namespace, name string) (*servicemeshv1alpha1.ServicePolicy, error) {
	app, err := deleteServicepolicy(namespace, name)

	if err != nil {
		klog.Error(err)
		return app, err
	}
	return app, nil
}

func createServicepolicy(namespace string, servicePolicy ServicePolicy) (*servicemeshv1alpha1.ServicePolicy, error) {

	k8sClient := client.ClientSets().K8s()

	theSp := servicePolicyParseParameter(namespace, servicePolicy)

	sp, err := k8sClient.Linkedcare().ServicemeshV1alpha1().ServicePolicies(namespace).Create(context.TODO(), theSp, metav1.CreateOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return sp, nil
}

func updateServicepolicy(namespace, name string, servicePolicy ServicePolicy) (*servicemeshv1alpha1.ServicePolicy, error) {

	oldServicePolicy, err := getServicepolicy(namespace, name)
	if err != nil {
		return nil, err
	}

	theSp := servicePolicyParseParameter(namespace, servicePolicy)
	if theSp.Name != name {
		return nil, errors.New("the name is not right")
	}

	newServicePolicy := oldServicePolicy.DeepCopy()
	newServicePolicy.Spec = theSp.Spec

	k8sClient := client.ClientSets().K8s()

	sp, err := k8sClient.Linkedcare().ServicemeshV1alpha1().ServicePolicies(namespace).Update(context.TODO(), newServicePolicy, metav1.UpdateOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return sp, nil
}

func servicePolicyParseParameter(namespace string, servicePolicy ServicePolicy) *servicemeshv1alpha1.ServicePolicy {
	sp := servicepolicyTemplates.DeepCopy()

	sp.Name = servicePolicy.Component

	// Modify selector
	sp.Spec.Selector.MatchLabels["app"] = servicePolicy.Component
	sp.Spec.Selector.MatchLabels["app.linkedcare.io/name"] = servicePolicy.Application.Name
	sp.Spec.Selector.MatchLabels["app.linkedcare.io/version"] = servicePolicy.Application.Version
	// Modify labels
	sp.Labels["app"] = servicePolicy.Component
	sp.Labels["app.linkedcare.io/name"] = servicePolicy.Application.Name
	sp.Labels["app.linkedcare.io/version"] = servicePolicy.Application.Version

	// Modify template
	sp.Spec.Template = servicePolicy.Template
	return sp
}

func deleteServicepolicy(namespace, name string) (*servicemeshv1alpha1.ServicePolicy, error) {

	app, err := getServicepolicy(namespace, name)

	if err != nil {
		klog.Error(err)
		return app, err
	}

	k8sClient := client.ClientSets().K8s()

	// delete servicepolicy
	deleteOptions := metav1.DeleteOptions{}

	err = k8sClient.Linkedcare().ServicemeshV1alpha1().ServicePolicies(namespace).Delete(context.TODO(), name, deleteOptions)
	if err != nil {
		klog.Error(err)
		return app, err
	}

	return app, nil
}

// Create an servicepolicy in a namespace
func UpdateServicepolicy(namespace, name string, servicePolicy ServicePolicy) (*servicemeshv1alpha1.ServicePolicy, error) {

	sp, err := updateServicepolicy(namespace, name, servicePolicy)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return sp, nil
}
