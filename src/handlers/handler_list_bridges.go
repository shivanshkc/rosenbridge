package handlers

import (
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/core/constants"
	"github.com/shivanshkc/rosenbridge/src/core/models"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// ListBridges is the handler for the GET Bridge List API of Rosenbridge.
func ListBridges(w http.ResponseWriter, r *http.Request) { // nolint:varnamelen // I like the "w" and "r" names.
	clientIDs := r.URL.Query()["client_id"]
	// Checking the slice of client IDs.
	if err := checkClientIDSlice(clientIDs); err != nil {
		// Converting to HTTP error.
		errHTTP := errutils.BadRequest().WithReasonError(err)
		// Sending back the response.
		httputils.Write(w, errHTTP.Status, nil, errHTTP)
		// Ending execution.
		return
	}

	// Calling the core function.
	bridges, err := core.ListBridges(r.Context(), clientIDs)
	if err != nil {
		// Converting to HTTP error.
		errHTTP := errutils.ToHTTPError(err)
		// Sending back the response.
		httputils.Write(w, errHTTP.Status, nil, errHTTP)
		// Ending execution.
		return
	}

	// Creating the response body.
	responseBody := &struct {
		Code    string              `json:"code"`
		Reason  string              `json:"reason"`
		Bridges []*models.BridgeDoc `json:"bridges"`
	}{
		Code:    constants.CodeOK,
		Reason:  "",
		Bridges: bridges,
	}

	// Sending back the response.
	httputils.Write(w, http.StatusOK, nil, responseBody)
}
