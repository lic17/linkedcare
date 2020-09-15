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

func GetAllStrategies(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")

	st, err := servicemesh.GetAllStrategies(namespace)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson(st)
}

func GetStrategy(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	st, err := servicemesh.GetStrategy(namespace, name)

	if err != nil {
		if k8serr.IsNotFound(err) {
			response.WriteHeaderAndEntity(http.StatusNotFound, lkerr.Wrap(err))
		} else {
			response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		}
		return
	}

	response.WriteAsJson(st)
}

func CreateStrategy(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")

	newStrategy := servicemesh.Strategy{}
	err := request.ReadEntity(&newStrategy)

	if err != nil {
		response.WriteAsJson(err)
		return
	}

	strategy, err := strategyParseParameter(newStrategy)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusBadRequest, lkerr.Wrap(fmt.Errorf("wrong name or version")))
		return
	}

	st, err := servicemesh.CreateStrategy(namespace, strategy)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson(st)

}

func DeleteStrategy(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")

	st, err := servicemesh.DeleteStrategy(namespace, name)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson(st)

}

func UpdateStrategy(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")

	newStrategy := servicemesh.Strategy{}
	err := request.ReadEntity(&newStrategy)

	if err != nil {
		response.WriteAsJson(err)
		return
	}

	strategy, err := strategyParseParameter(newStrategy)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusBadRequest, lkerr.Wrap(fmt.Errorf("wrong name or version")))
		return
	}

	st, err := servicemesh.UpdateStrategy(namespace, name, strategy)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson(st)
}

func strategyParseParameter(st servicemesh.Strategy) (servicemesh.Strategy, error) {
	var err error
	err = nil
	if st.Application.Name == "" {
		err = errors.New("the name of application not set")
		return st, err
	}
	if st.Application.Version == "" {
		err = errors.New("the version of application not set")
		return st, err
	}

	return st, err
}
