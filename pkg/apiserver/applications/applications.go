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
	"errors"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	k8serr "k8s.io/apimachinery/pkg/api/errors"

	lkerr "linkedcare.io/linkedcare/pkg/server/errors"

	"linkedcare.io/linkedcare/pkg/models/applications"
)

type Application struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Get all application
func GetAllApplications(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")
	applications, err := applications.GetAllApplications(namespace)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson(applications)
}

// Get an application for specified namespace
func GetApplication(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	application, err := applications.GetApplication(namespace, name)

	if err != nil {
		if k8serr.IsNotFound(err) {
			response.WriteHeaderAndEntity(http.StatusNotFound, lkerr.Wrap(err))
		} else {
			response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		}
		return
	}

	response.WriteAsJson(application)
}

// Get an application details for specified namespace
func GetApplicationDetails(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	application, err := applications.GetApplicationDetails(namespace, name)

	if err != nil {
		if k8serr.IsNotFound(err) {
			response.WriteHeaderAndEntity(http.StatusNotFound, lkerr.Wrap(err))
		} else {
			response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		}
		return
	}

	response.WriteAsJson(application)
}

// Create an application details for specified namespace
func CreateApplicationDetails(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")

	newApplicationDetails := applications.ApplicationDetails{}
	err := request.ReadEntity(&newApplicationDetails)

	if err != nil {
		response.WriteAsJson(err)
		return
	}

	application, err := applications.CreateApplicationDetails(&newApplicationDetails, namespace)

	if err != nil {
		if k8serr.IsNotFound(err) {
			response.WriteHeaderAndEntity(http.StatusNotFound, lkerr.Wrap(err))
		} else {
			response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		}
		return
	}

	response.WriteAsJson(application)
}

// Resatrt an application for specified namespace
func RestartApplication(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")

	err := applications.RestartApplication(namespace, name)
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson("restarting successfully")
}

// Resatrt an service for specified namespace
func RestartService(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")

	err := applications.RestartService(namespace, name)
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson("restarting successfully")
}

// Create application
func CreateApplication(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")

	newApplication := Application{}
	err := request.ReadEntity(&newApplication)

	if err != nil {
		response.WriteAsJson(err)
		return
	}

	name, version, err := parseParameter(newApplication)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusBadRequest, lkerr.Wrap(fmt.Errorf("wrong name or version")))
		return
	}

	app, err := applications.CreateApplication(namespace, name, version)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson(app)
}

// Delete application
func DeleteApplication(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")

	app, err := applications.DeleteApplication(namespace, name)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
		return
	}

	response.WriteAsJson(app)
}

func parseParameter(app Application) (name, version string, err error) {

	err = nil
	if app.Name == "" {
		err = errors.New("name not set")
		return
	}
	if app.Version == "" {
		err = errors.New("version not set")
		return
	}
	name = app.Name
	version = app.Version

	return
}
