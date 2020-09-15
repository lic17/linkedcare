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
package v1alpha1

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"linkedcare.io/linkedcare/pkg/apiserver/runtime"
	"linkedcare.io/linkedcare/pkg/apiserver/servicemesh"
	"linkedcare.io/linkedcare/pkg/server/errors"
)

const GroupName = "servicemesh.linkedcare.io"

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

var (
	WebServiceBuilder = runtime.NewContainerBuilder(addWebService)
	AddToContainer    = WebServiceBuilder.AddToContainer
)

func addWebService(c *restful.Container) error {

	ok := "ok"
	webservice := runtime.NewWebService(GroupVersion)

	webservice.Route(webservice.POST("/servicepolicy/namespace/{namespace}").
		To(servicemesh.CreateServicepolicy).
		Deprecate().
		Doc("Create servicepolicy for service in the namespace").
		Param(webservice.PathParameter("namespace", "the name of the namespace where the servicepolicy create in")).
		Returns(http.StatusOK, ok, errors.Error{}))

	webservice.Route(webservice.POST("/strategy/namespace/{namespace}").
		To(servicemesh.CreateStrategy).
		Doc("Create strategy for service in the namespace").
		Deprecate().
		Param(webservice.PathParameter("namespace", "the name of the namespace where the strategy create in")).
		Param(webservice.QueryParameter("action", "action must be \"rerun\"")).
		Param(webservice.QueryParameter("resourceVersion", "version of job, rerun when the version matches").Required(true)).
		Returns(http.StatusOK, ok, errors.Error{}))
	/*
		webservice.Route(webservice.GET("/servicepolicy").
			To(routers.GetAllRouters).
			Doc("List all routers of all projects").
			Metadata(restfulspec.KeyOpenAPITags, []string{constants.ClusterResourcesTag}).
			Returns(http.StatusOK, ok, corev1.ServiceList{}))

		webservice.Route(webservice.GET("/namespaces/{namespace}/router").
			To(routers.GetRouter).
			Doc("List router of a specified project").
			Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
			Returns(http.StatusOK, ok, corev1.Service{}).
			Param(webservice.PathParameter("namespace", "the name of the project")))

		webservice.Route(webservice.DELETE("/namespaces/{namespace}/router").
			To(routers.DeleteRouter).
			Doc("List router of a specified project").
			Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
			Returns(http.StatusOK, ok, corev1.Service{}).
			Param(webservice.PathParameter("namespace", "the name of the project")))

		webservice.Route(webservice.POST("/namespaces/{namespace}/router").
			To(routers.CreateRouter).
			Doc("Create a router for a specified project").
			Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
			Returns(http.StatusOK, ok, corev1.Service{}).
			Param(webservice.PathParameter("namespace", "the name of the project")))

		webservice.Route(webservice.PUT("/namespaces/{namespace}/router").
			To(routers.UpdateRouter).
			Doc("Update a router for a specified project").
			Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
			Returns(http.StatusOK, ok, corev1.Service{}).
			Param(webservice.PathParameter("namespace", "the name of the project")))
	*/

	c.Add(webservice)

	return nil
}
