package handlers

import (
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/core/constants"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// GetIntro is the handler for the GetIntro API. It is used to check whether the service is down or running.
func GetIntro(w http.ResponseWriter, r *http.Request) {
	// Body of the HTTP response.
	responseBody := map[string]string{"code": constants.CodeOK}
	// Writing the response.
	httputils.Write(w, http.StatusOK, nil, responseBody)
}
