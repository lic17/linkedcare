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
	restfulspec "github.com/emicklei/go-restful-openapi"
	"k8s.io/apimachinery/pkg/runtime/schema"
	servicemeshv1alpha1 "linkedcare.io/linkedcare/pkg/apis/servicemesh/v1alpha1"
	"linkedcare.io/linkedcare/pkg/apiserver/runtime"
	"linkedcare.io/linkedcare/pkg/apiserver/servicemesh"
	"linkedcare.io/linkedcare/pkg/constants"
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

	//
	//add router for servicepolicy
	//
	webservice.Route(webservice.GET("/namespace/{namespace}/servicepolicy").
		To(servicemesh.GetAllServicepolicies).
		Doc("List all servicepolicies in an namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ClusterResourcesTag}).
		Returns(http.StatusOK, ok, servicemeshv1alpha1.ServicePolicyList{}).
		Param(webservice.PathParameter("namespace", "the namespace of the servicepolicies")))

	webservice.Route(webservice.GET("/namespace/{namespace}/name/{name}/servicepolicy").
		To(servicemesh.GetServicepolicy).
		Doc("Get a servicepolicy in an namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, servicemeshv1alpha1.ServicePolicy{}).
		Param(webservice.PathParameter("namespace", "the namespace of the servicepolicy")).
		Param(webservice.PathParameter("name", "the name of the servicepolicy")))

	webservice.Route(webservice.DELETE("/namespace/{namespace}/name/{name}/servicepolicy").
		To(servicemesh.DeleteServicepolicy).
		Doc("Delete a servicepolicy in a specified namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, servicemeshv1alpha1.ServicePolicy{}).
		Param(webservice.PathParameter("namespace", "the namespace of the servicepolicy")).
		Param(webservice.PathParameter("name", "the name of the servicepolicy")))

	webservice.Route(webservice.POST("/namespace/{namespace}/servicepolicy").
		To(servicemesh.CreateServicepolicy).
		Doc("Create a servicepolicy in a namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, servicemeshv1alpha1.ServicePolicy{}).
		Param(webservice.PathParameter("namespace", "the namespace of the servicepolicy")))

	webservice.Route(webservice.PUT("/namespace/{namespace}/name/{name}/servicepolicy").
		To(servicemesh.UpdateServicepolicy).
		Doc("Update a servicepolicy  in a namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, servicemeshv1alpha1.ServicePolicy{}).
		Param(webservice.PathParameter("namespace", "the namespace of the servicepolicy")).
		Param(webservice.PathParameter("name", "the name of the servicepolicy")))
	//
	//add router for strategy
	//
	webservice.Route(webservice.GET("/namespace/{namespace}/strategy").
		To(servicemesh.GetAllStrategies).
		Doc("List all strategies in an namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ClusterResourcesTag}).
		Returns(http.StatusOK, ok, servicemeshv1alpha1.StrategyList{}).
		Param(webservice.PathParameter("namespace", "the namespace of the strategies")))

	webservice.Route(webservice.GET("/namespace/{namespace}/name/{name}/strategy").
		To(servicemesh.GetStrategy).
		Doc("Get a strategy in an namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, servicemeshv1alpha1.Strategy{}).
		Param(webservice.PathParameter("namespace", "the namespace of the strategy")).
		Param(webservice.PathParameter("name", "the name of the strategy")))

	webservice.Route(webservice.DELETE("/namespace/{namespace}/name/{name}/strategy").
		To(servicemesh.DeleteStrategy).
		Doc("Delete a strategy in a specified namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, servicemeshv1alpha1.Strategy{}).
		Param(webservice.PathParameter("namespace", "the namespace of the strategy")).
		Param(webservice.PathParameter("name", "the name of the strategy")))

	webservice.Route(webservice.POST("/namespace/{namespace}/strategy").
		To(servicemesh.CreateStrategy).
		Doc("Create a strategy in a namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, servicemeshv1alpha1.Strategy{}).
		Param(webservice.PathParameter("namespace", "the namespace of the strategy")))

	webservice.Route(webservice.PUT("/namespace/{namespace}/name/{name}/strategy").
		To(servicemesh.UpdateStrategy).
		Doc("Update a strategy  in a namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceResourcesTag}).
		Returns(http.StatusOK, ok, servicemeshv1alpha1.Strategy{}).
		Param(webservice.PathParameter("namespace", "the namespace of the strategy")).
		Param(webservice.PathParameter("name", "the name of the strategy")))

	c.Add(webservice)

	return nil
}
