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
package v1beta1

import (
	"net/http"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"linkedcare.io/linkedcare/pkg/apiserver/applications"
	"linkedcare.io/linkedcare/pkg/apiserver/runtime"
	"linkedcare.io/linkedcare/pkg/constants"
	appModels "linkedcare.io/linkedcare/pkg/models/applications"
	appv1beta1 "sigs.k8s.io/application/api/v1beta1"
)

const GroupName = "app.k8s.io"

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}

var (
	WebServiceBuilder = runtime.NewContainerBuilder(addWebService)
	AddToContainer    = WebServiceBuilder.AddToContainer
)

func addWebService(c *restful.Container) error {

	ok := "ok"
	webservice := runtime.NewWebService(GroupVersion)

	webservice.Route(webservice.GET("/namespace/{namespace}/application").
		To(applications.GetAllApplications).
		Doc("List all applications in an namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ClusterResourcesTag}).
		Returns(http.StatusOK, ok, appv1beta1.ApplicationList{}).
		Param(webservice.PathParameter("namespace", "the namespace of the applications")))

	webservice.Route(webservice.GET("/namespace/{namespace}/name/{name}/application").
		To(applications.GetApplication).
		Doc("Get an application in an namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, appv1beta1.Application{}).
		Param(webservice.PathParameter("namespace", "the namespace of the application")).
		Param(webservice.PathParameter("name", "the name of the application")))

	webservice.Route(webservice.GET("/namespace/{namespace}/name/{name}/application/details").
		To(applications.GetApplicationDetails).
		Doc("Get an application detail in an namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, appModels.ApplicationDetails{}).
		Param(webservice.PathParameter("namespace", "the namespace of the application")).
		Param(webservice.PathParameter("name", "the name of the application")))

	webservice.Route(webservice.POST("/namespace/{namespace}/application/details").
		To(applications.CreateApplicationDetails).
		Doc("Create an application detail in an namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, appModels.ApplicationDetails{}).
		Param(webservice.PathParameter("namespace", "the namespace of the application")))

	webservice.Route(webservice.DELETE("/namespace/{namespace}/name/{name}/application").
		To(applications.DeleteApplication).
		Doc("Delete an application in a specified namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, appv1beta1.Application{}).
		Param(webservice.PathParameter("namespace", "the namespace of the application")).
		Param(webservice.PathParameter("name", "the name of the application")))

	webservice.Route(webservice.PUT("/namespace/{namespace}/name/{name}/application").
		To(applications.RestartApplication).
		Doc("Restart an application in a specified namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, appv1beta1.Application{}).
		Param(webservice.PathParameter("namespace", "the namespace of the application")).
		Param(webservice.PathParameter("name", "the name of the application")))

	webservice.Route(webservice.PUT("/namespace/{namespace}/name/{name}/service").
		To(applications.RestartService).
		Doc("Restart a service in a specified namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, appv1beta1.Application{}).
		Param(webservice.PathParameter("namespace", "the namespace of the service")).
		Param(webservice.PathParameter("name", "the name of the service")))

	webservice.Route(webservice.POST("/namespace/{namespace}/application").
		To(applications.CreateApplication).
		Doc("Create an application in a namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, appv1beta1.Application{}).
		Param(webservice.PathParameter("namespace", "the namespace of the application")))

	c.Add(webservice)

	return nil
}
