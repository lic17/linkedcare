package git

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"linkedcare.io/linkedcare/pkg/models/git"
	"linkedcare.io/linkedcare/pkg/server/errors"
)

func GitReadVerify(request *restful.Request, response *restful.Response) {

	authInfo := git.AuthInfo{}

	err := request.ReadEntity(&authInfo)
	ns := request.PathParameter("namespace")
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, errors.Wrap(err))
		return
	}

	err = git.GitReadVerify(ns, authInfo)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, errors.Wrap(err))
		return
	}

	response.WriteAsJson(errors.None)
}
