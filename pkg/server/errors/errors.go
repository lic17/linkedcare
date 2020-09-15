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
package errors

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

type Error struct {
	Message string `json:"message" description:"error message"`
}

var None = Error{Message: "success"}

func (e *Error) Error() string {
	return e.Message
}

func Wrap(err error) Error {
	return Error{Message: err.Error()}
}

func New(message string) Error {
	return Error{Message: message}
}

func ParseSvcErr(err error, resp *restful.Response) {
	if svcErr, ok := err.(restful.ServiceError); ok {
		resp.WriteServiceError(svcErr.Code, svcErr)
	} else {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, Wrap(err))
	}
}
