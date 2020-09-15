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
	"linkedcare.io/linkedcare/pkg/apiserver/certs"
	"linkedcare.io/linkedcare/pkg/apiserver/runtime"
	"linkedcare.io/linkedcare/pkg/constants"
	"linkedcare.io/linkedcare/pkg/models/cert"
)

const GroupName = "cert.linkedcare.io"

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}

var (
	WebServiceBuilder = runtime.NewContainerBuilder(addWebService)
	AddToContainer    = WebServiceBuilder.AddToContainer
)

func addWebService(c *restful.Container) error {

	ok := "ok"
	webservice := runtime.NewWebService(GroupVersion)

	webservice.Route(webservice.GET("/certs").
		To(certs.GetAllCerts).
		Doc("List all applications in an namespace").
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ClusterResourcesTag}).
		Returns(http.StatusOK, ok, []cert.CertTime{}))

	c.Add(webservice)

	return nil
}
