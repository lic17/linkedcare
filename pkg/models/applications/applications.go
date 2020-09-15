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

package applications

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog"
	"linkedcare.io/linkedcare/pkg/simple/client"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	Ingressv1beta1 "k8s.io/api/networking/v1beta1"
	"linkedcare.io/linkedcare/pkg/informers"
	appv1beta1 "sigs.k8s.io/application/api/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	servicemeshEnabled = "servicemesh.linkedcare.io/enabled"
)

type NewApp struct {
	Name    string
	Version string
}

type ApplicationDetails struct {
	NewApp *NewApp
	App    *appv1beta1.Application
	//Ingresses   []*Ingressv1beta1.Ingress
	Deployments []*v1.Deployment
	Services    []*corev1.Service
}

var appTemplates appv1beta1.Application

// Load yamls
func init() {

	groupkinds := []metav1.GroupKind{
		{Group: "", Kind: "Service"},
		{Group: "apps", Kind: "Deployment"},
		{Group: "apps", Kind: "StatefulSet"},
		{Group: "extensions", Kind: "Ingress"},
		{Group: "servicemesh.linkedcare.io", Kind: "Strategy"},
		{Group: "servicemesh.linkedcare.io", Kind: "ServicePolicy"},
	}

	appLabes := map[string]string{
		"app.linkedcare.io/name":    "",
		"app.linkedcare.io/version": "",
	}

	selector := &metav1.LabelSelector{
		MatchLabels: appLabes,
	}

	appTemplates.Spec.ComponentGroupKinds = groupkinds
	appTemplates.Spec.AddOwnerRef = true
	appTemplates.Spec.Selector = selector
	appTemplates.Labels = appLabes

	appTemplates.Annotations = make(map[string]string)
	appTemplates.Annotations[servicemeshEnabled] = "true"
}

// Get all applications from a namespace
func GetAllApplications(namespace string) ([]*appv1beta1.Application, error) {

	appLister := informers.AppSharedInformerFactory().App().V1beta1().Applications().Lister()
	apps, err := appLister.Applications(namespace).List(labels.Everything())

	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return apps, nil
}

// Get an application from a namespace
func GetApplication(namespace, name string) (*appv1beta1.Application, error) {
	app, err := getApplication(namespace, name)
	return app, err
}

func getApplication(namespace, name string) (*appv1beta1.Application, error) {
	appLister := informers.AppSharedInformerFactory().App().V1beta1().Applications().Lister()
	app, err := appLister.Applications(namespace).Get(name)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NewNotFound(appv1beta1.Resource("applications"), name)
		}
		klog.Error(err)
		return nil, err
	}
	return app, nil
}

// Get an application detail from a namespace
func GetApplicationDetails(namespace, name string) (*ApplicationDetails, error) {
	app, err := getApplicationDetails(namespace, name)
	return app, err
}

func getApplicationDetails(namespace, name string) (*ApplicationDetails, error) {
	appLister := informers.AppSharedInformerFactory().App().V1beta1().Applications().Lister()
	application, err := appLister.Applications(namespace).Get(name)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NewNotFound(appv1beta1.Resource("applications"), name)
		}
		klog.Error(err)
		return nil, err
	}
	app := application.DeepCopy()
	deployments, _ := getDeployments(app)
	services, _ := getServices(app)
	//ingresses, _ := getIngresses(app)
	details := &ApplicationDetails{
		App:         app,
		Deployments: deployments,
		Services:    services,
		//	Ingresses:   ingresses,
	}
	return details, nil
}

func getDeployments(app *appv1beta1.Application) ([]*v1.Deployment, error) {
	depLister := informers.SharedInformerFactory().Apps().V1().Deployments().Lister()

	labelSelector := labels.Set(app.Labels).AsSelector()
	deployments, err := depLister.Deployments(app.Namespace).List(labelSelector)
	if err != nil {
		return deployments, err
	}
	return deployments, nil
}

func getServices(app *appv1beta1.Application) ([]*corev1.Service, error) {
	svcLister := informers.SharedInformerFactory().Core().V1().Services().Lister()

	labelSelector := labels.Set(app.Labels).AsSelector()
	services, err := svcLister.Services(app.Namespace).List(labelSelector)
	if err != nil {
		return services, err
	}
	return services, nil
}

func CreateApplicationDetails(details *ApplicationDetails, namespace string) (*ApplicationDetails, error) {
	return createApplicationDetails(details, namespace)
}

func createApplicationDetails(details *ApplicationDetails, namespace string) (*ApplicationDetails, error) {

	// Create Application
	appCp := details.App.DeepCopy()
	appName := appCp.Name
	appVersion := ""
	if v, ok := appCp.Labels["app.linkedcare.io/version"]; ok {
		appVersion = v
	}
	if details.NewApp != nil {
		if details.NewApp.Name != "" {
			appName = details.NewApp.Name
		}
		if details.NewApp.Version != "" {
			appVersion = details.NewApp.Version
		}
	}
	appCp.Labels["app.linkedcare.io/version"] = appVersion
	appCp.Labels["app.linkedcare.io/name"] = appName
	appCp.Spec.Selector.MatchLabels["app.linkedcare.io/version"] = appVersion
	appCp.Spec.Selector.MatchLabels["app.linkedcare.io/name"] = appName
	appCp.Name = appName

	app := &appv1beta1.Application{}
	app.SetAnnotations(appCp.Annotations)
	app.SetLabels(appCp.Labels)
	app.SetName(appCp.Name)
	app.SetNamespace(namespace)
	app.Spec = appCp.Spec

	k8sClient := client.ClientSets().K8s()

	app, err := k8sClient.Application().AppV1beta1().Applications(namespace).Create(context.TODO(), app, metav1.CreateOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	// Create Deployments
	for _, deployment := range details.Deployments {
		depCp := deployment.DeepCopy()

		depCp.Labels["app.linkedcare.io/version"] = appVersion
		depCp.Labels["app.linkedcare.io/name"] = appName
		depCp.Spec.Selector.MatchLabels["app.linkedcare.io/version"] = appVersion
		depCp.Spec.Selector.MatchLabels["app.linkedcare.io/name"] = appName
		depCp.Spec.Template.Labels["app.linkedcare.io/version"] = appVersion
		depCp.Spec.Template.Labels["app.linkedcare.io/name"] = appName

		dep := &v1.Deployment{}
		dep.SetAnnotations(depCp.Annotations)
		dep.SetLabels(depCp.Labels)
		dep.SetName(depCp.Name)
		dep.SetNamespace(namespace)
		dep.Spec = depCp.Spec

		k8sClient := client.ClientSets().K8s()

		dep, err := k8sClient.Kubernetes().AppsV1().Deployments(namespace).Create(context.TODO(), dep, metav1.CreateOptions{})
		if err != nil {
			klog.Error(err)
			return nil, err
		}
	}

	// Create Services
	for _, service := range details.Services {
		svcCp := service.DeepCopy()

		svcCp.Labels["app.linkedcare.io/version"] = appVersion
		svcCp.Labels["app.linkedcare.io/name"] = appName

		svc := &corev1.Service{}
		svc.SetAnnotations(svcCp.Annotations)
		svc.SetLabels(svcCp.Labels)
		svc.SetName(svcCp.Name)
		svc.SetNamespace(namespace)

		svc.Spec.Ports = svcCp.Spec.Ports
		for i, _ := range svc.Spec.Ports {
			svc.Spec.Ports[i].NodePort = 0
		}
		svc.Spec.Selector = svcCp.Spec.Selector
		svc.Spec.Selector["app.linkedcare.io/version"] = appVersion
		svc.Spec.Selector["app.linkedcare.io/name"] = appName
		svc.Spec.SessionAffinityConfig = svcCp.Spec.SessionAffinityConfig.DeepCopy()

		k8sClient := client.ClientSets().K8s()

		svc, err := k8sClient.Kubernetes().CoreV1().Services(namespace).Create(context.TODO(), svc, metav1.CreateOptions{})
		if err != nil {
			klog.Error(err)
			return nil, err
		}
	}

	return details, nil
}

func getIngresses(app *appv1beta1.Application) ([]*Ingressv1beta1.Ingress, error) {
	ingressLister := informers.SharedInformerFactory().Networking().V1beta1().Ingresses().Lister()

	labelSelector := labels.Set(app.Labels).AsSelector()
	ingresses, err := ingressLister.Ingresses(app.Namespace).List(labelSelector)
	if err != nil {
		return ingresses, err
	}
	return ingresses, nil
}

// Create an application in a namespace
func CreateApplication(namespace, name, version string) (*appv1beta1.Application, error) {

	app, err := createApplication(namespace, name, version)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return app, nil
}

// DeleteApplication is used to delete the application in namespace
func DeleteApplication(namespace, name string) (*appv1beta1.Application, error) {
	app, err := deleteApplication(namespace, name)

	if err != nil {
		klog.Error(err)
		return app, err
	}
	return app, nil
}

// RestartApplication is used to restart the application in namespace
func RestartApplication(namespace, name string) error {
	err := restartApplication(namespace, name)

	if err != nil {
		klog.Error(err)
		return err
	}
	return nil
}

// RestartService is used to restart the service in namespace
func RestartService(namespace, name string) error {
	err := restartService(namespace, name)

	if err != nil {
		klog.Error(err)
		return err
	}
	return nil
}

func createApplication(namespace, name, version string) (*appv1beta1.Application, error) {

	app := appTemplates.DeepCopy()

	k8sClient := client.ClientSets().K8s()
	app.Name = name

	// Modify selector
	app.Spec.Selector.MatchLabels["app.linkedcare.io/name"] = name
	app.Spec.Selector.MatchLabels["app.linkedcare.io/version"] = version
	// Modify labels
	app.Labels["app.linkedcare.io/name"] = name
	app.Labels["app.linkedcare.io/version"] = version

	app, err := k8sClient.Application().AppV1beta1().Applications(namespace).Create(context.TODO(), app, metav1.CreateOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return app, nil
}

func deleteApplication(namespace, name string) (*appv1beta1.Application, error) {

	app, err := getApplication(namespace, name)

	if err != nil {
		klog.Error(err)
		return app, err
	}

	k8sClient := client.ClientSets().K8s()

	// delete application
	deleteOptions := metav1.DeleteOptions{}

	err = k8sClient.Application().AppV1beta1().Applications(namespace).Delete(context.TODO(), name, deleteOptions)
	if err != nil {
		klog.Error(err)
		return app, err
	}

	return app, nil
}

func restartApplication(namespace, name string) error {

	app, err := getApplication(namespace, name)

	if err != nil {
		klog.Error(err)
		return err
	}

	labelSelector := labels.FormatLabels(app.Spec.Selector.MatchLabels)
	k8sClient := client.ClientSets().K8s()

	// delete pod
	deleteOptions := metav1.DeleteOptions{}
	listOpts := metav1.ListOptions{
		LabelSelector: labelSelector,
	}

	err = k8sClient.Kubernetes().CoreV1().Pods(namespace).DeleteCollection(context.TODO(), deleteOptions, listOpts)
	if err != nil {
		klog.Error(err)
		return err
	}

	return nil
}

func restartService(namespace, name string) error {

	svc, err := getService(namespace, name)

	if err != nil {
		klog.Error(err)
		return err
	}

	labelSelector := labels.FormatLabels(svc.Spec.Selector)
	k8sClient := client.ClientSets().K8s()

	// delete pod
	deleteOptions := metav1.DeleteOptions{}
	listOpts := metav1.ListOptions{
		LabelSelector: labelSelector,
	}

	err = k8sClient.Kubernetes().CoreV1().Pods(namespace).DeleteCollection(context.TODO(), deleteOptions, listOpts)
	if err != nil {
		klog.Error(err)
		return err
	}

	return nil
}

func getService(namespace, name string) (*corev1.Service, error) {
	svcLister := informers.SharedInformerFactory().Core().V1().Services().Lister()
	svc, err := svcLister.Services(namespace).Get(name)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NewNotFound(corev1.Resource("servicers"), name)
		}
		klog.Error(err)
		return nil, err
	}
	return svc, nil
}
