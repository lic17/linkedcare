package api

import (
	"github.com/emicklei/go-restful"
	urlruntime "k8s.io/apimachinery/pkg/util/runtime"

	//	operationsv1alpha2 "linkedcare.io/linkedcare/pkg/api/operations/v1alpha2"
	//	resourcesv1alpha2 "linkedcare.io/linkedcare/pkg/api/resources/v1alpha2"
	//	terminalv1alpha2 "linkedcare.io/linkedcare/pkg/api/terminal/v1alpha2"
	appv1beta1 "linkedcare.io/linkedcare/pkg/api/application/v1beta1"
	certv1beta1 "linkedcare.io/linkedcare/pkg/api/cert/v1beta1"
	servicemeshv1alpha1 "linkedcare.io/linkedcare/pkg/api/servicemesh/v1alpha1"
)

func InstallAPIs(container *restful.Container) {
	//urlruntime.Must(operationsv1alpha2.AddToContainer(container))
	//urlruntime.Must(resourcesv1alpha2.AddToContainer(container))
	//urlruntime.Must(terminalv1alpha2.AddToContainer(container))
	urlruntime.Must(servicemeshv1alpha1.AddToContainer(container))
	urlruntime.Must(appv1beta1.AddToContainer(container))
	urlruntime.Must(certv1beta1.AddToContainer(container))
}
