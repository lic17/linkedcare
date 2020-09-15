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

type Strategy struct {
	Component        string                                         `json:"component"`
	Application      applications.Application                       `json:"application"`
	PrincipalVersion string                                         `json:"principal,omitempty"`
	GovernorVersion  string                                         `json:"governor"`
	Template         servicemeshv1alpha1.VirtualServiceTemplateSpec `json:"template,omitempty"`
}

var strategyTemplates servicemeshv1alpha1.Strategy

// Load yamls
func init() {

	strategyLabes := map[string]string{
		"app":                       "",
		"app.linkedcare.io/name":    "",
		"app.linkedcare.io/version": "",
	}

	selector := &metav1.LabelSelector{
		MatchLabels: strategyLabes,
	}

	strategyTemplates.Labels = strategyLabes
	strategyTemplates.Spec.Selector = selector
	//TODO add type of bluecreen and mirror
	strategyTemplates.Spec.Type = servicemeshv1alpha1.CanaryType
	strategyTemplates.Spec.StrategyPolicy = servicemeshv1alpha1.PolicyWaitForWorkloadReady
}

// Get all strategys from a namespace
func GetAllStrategies(namespace string) ([]*servicemeshv1alpha1.Strategy, error) {

	stLister := informers.LcSharedInformerFactory().Servicemesh().V1alpha1().Strategies().Lister()
	sts, err := stLister.Strategies(namespace).List(labels.Everything())

	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return sts, nil
}

// Get an strategy from a namespace
func GetStrategy(namespace, name string) (*servicemeshv1alpha1.Strategy, error) {
	st, err := getStrategy(namespace, name)
	return st, err
}

func getStrategy(namespace, name string) (*servicemeshv1alpha1.Strategy, error) {
	stLister := informers.LcSharedInformerFactory().Servicemesh().V1alpha1().Strategies().Lister()
	st, err := stLister.Strategies(namespace).Get(name)

	if err != nil {
		if k8serr.IsNotFound(err) {
			return nil, k8serr.NewNotFound(servicemeshv1alpha1.Resource("strategys"), name)
		}
		klog.Error(err)
		return nil, err
	}
	return st, nil
}

// Create an strategy in a namespace
func CreateStrategy(namespace string, strategy Strategy) (*servicemeshv1alpha1.Strategy, error) {
	st, err := createStrategy(namespace, strategy)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return st, nil
}

// DeleteStrategy is used to delete the strategy in namespace
func DeleteStrategy(namespace, name string) (*servicemeshv1alpha1.Strategy, error) {
	st, err := deleteStrategy(namespace, name)

	if err != nil {
		klog.Error(err)
		return st, err
	}
	return st, nil
}

func createStrategy(namespace string, strategy Strategy) (*servicemeshv1alpha1.Strategy, error) {

	k8sClient := client.ClientSets().K8s()

	theStrategy := strategyParseParameter(namespace, strategy)
	st, err := k8sClient.Linkedcare().ServicemeshV1alpha1().Strategies(namespace).Create(context.TODO(), theStrategy, metav1.CreateOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return st, nil
}

func updateStrategy(namespace, name string, strategy Strategy) (*servicemeshv1alpha1.Strategy, error) {
	oldStrategy, err := getStrategy(namespace, name)
	if err != nil {
		return nil, err
	}

	theStrategy := strategyParseParameter(namespace, strategy)
	if theStrategy.Name != name {
		return nil, errors.New("the name is not right")
	}

	newStrategy := oldStrategy.DeepCopy()
	newStrategy.Spec = theStrategy.Spec

	k8sClient := client.ClientSets().K8s()

	st, err := k8sClient.Linkedcare().ServicemeshV1alpha1().Strategies(namespace).Update(context.TODO(), newStrategy, metav1.UpdateOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return st, nil
}

func strategyParseParameter(namespace string, strategy Strategy) *servicemeshv1alpha1.Strategy {
	st := strategyTemplates.DeepCopy()

	st.Name = strategy.Component

	// Modify selector
	st.Spec.Selector.MatchLabels["app"] = strategy.Component
	st.Spec.Selector.MatchLabels["app.linkedcare.io/name"] = strategy.Application.Name
	st.Spec.Selector.MatchLabels["app.linkedcare.io/version"] = strategy.Application.Version
	// Modify labels
	st.Labels["app"] = strategy.Component
	st.Labels["app.linkedcare.io/name"] = strategy.Application.Name
	st.Labels["app.linkedcare.io/version"] = strategy.Application.Version

	//modify principal version
	st.Spec.PrincipalVersion = strategy.PrincipalVersion

	//modify governor version
	st.Spec.GovernorVersion = strategy.GovernorVersion

	// Modify template
	st.Spec.Template = strategy.Template
	return st
}

func deleteStrategy(namespace, name string) (*servicemeshv1alpha1.Strategy, error) {

	st, err := getStrategy(namespace, name)

	if err != nil {
		klog.Error(err)
		return st, err
	}

	if len(st.Spec.GovernorVersion) > 0 {

		k8sClient := client.ClientSets().K8s()

		// delete strategy
		deleteOptions := metav1.DeleteOptions{}

		err = k8sClient.Linkedcare().ServicemeshV1alpha1().Strategies(namespace).Delete(context.TODO(), name, deleteOptions)
		if err != nil {
			klog.Error(err)
			return st, err
		}
	} else {
		err = k8serr.NewBadRequest("Governor version is null, can't delete the strategy")
	}

	return st, err
}

// Create an strategy in a namespace
func UpdateStrategy(namespace, name string, strategy Strategy) (*servicemeshv1alpha1.Strategy, error) {

	st, err := updateStrategy(namespace, name, strategy)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return st, nil
}
