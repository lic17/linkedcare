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
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"

	"errors"

	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"linkedcare.io/linkedcare/pkg/models/servicemesh"
	lkerr "linkedcare.io/linkedcare/pkg/server/errors"
)

func GetAllServicepolicies(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")

	sp, err := servicemesh.GetAllServicepolicies(namespace)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson(sp)
}

func GetServicepolicy(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	sp, err := servicemesh.GetServicepolicy(namespace, name)

	if err != nil {
		if k8serr.IsNotFound(err) {
			response.WriteHeaderAndEntity(http.StatusNotFound, lkerr.Wrap(err))
		} else {
			response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		}
		return
	}

	response.WriteAsJson(sp)
}

func CreateServicepolicy(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")

	newServicePolicy := servicemesh.ServicePolicy{}
	err := request.ReadEntity(&newServicePolicy)

	if err != nil {
		response.WriteAsJson(err)
		return
	}

	servicePolicy, err := servicepolicyParseParameter(newServicePolicy)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusBadRequest, lkerr.Wrap(fmt.Errorf("wrong name or version")))
		return
	}

	sp, err := servicemesh.CreateServicepolicy(namespace, servicePolicy)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson(sp)

}

func DeleteServicepolicy(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")

	app, err := servicemesh.DeleteServicepolicy(namespace, name)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson(app)

}

func UpdateServicepolicy(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")

	newServicePolicy := servicemesh.ServicePolicy{}
	err := request.ReadEntity(&newServicePolicy)

	if err != nil {
		response.WriteAsJson(err)
		return
	}

	servicePolicy, err := servicepolicyParseParameter(newServicePolicy)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusBadRequest, lkerr.Wrap(fmt.Errorf("wrong name or version")))
		return
	}

	sp, err := servicemesh.UpdateServicepolicy(namespace, name, servicePolicy)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson(sp)
}

func servicepolicyParseParameter(sp servicemesh.ServicePolicy) (servicemesh.ServicePolicy, error) {

	var err error
	err = nil
	if sp.Application.Name == "" {
		err = errors.New("the name of application not set")
		return sp, err
	}
	if sp.Application.Version == "" {
		err = errors.New("the version of application not set")
		return sp, err
	}

	return sp, err
}
