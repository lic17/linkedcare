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

package certs

import (
	"github.com/emicklei/go-restful"

	"linkedcare.io/linkedcare/pkg/models/cert"
)

// Get all cert
func GetAllCerts(request *restful.Request, response *restful.Response) {

	certs, _ := cert.GetCertTime()

	//TODO return err
	//if err != nil {
	//	response.WriteHeaderAndEntity(http.StatusInternalServerError, lkerr.Wrap(err))
	//	return
	//}

	response.WriteAsJson(certs)
}
