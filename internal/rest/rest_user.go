package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shivanshkc/rosenbridge/pkg/utils/httputils"
)

// createUser is the API handler for the POST /api/user route.
func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Anonymous struct variable to decode request body.
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Read request body.
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.ErrorContext(ctx, "failed to read request body", "error", err)
		httputils.WriteError(w, httputils.BadRequest().WithReasonStr("failed to read request body"))
		return
	}

	if err := validateUsername(body.Username); err != nil {
		slog.ErrorContext(ctx, "invalid username", "error", err)
		httputils.WriteError(w, httputils.BadRequest().WithReasonErr(err))
		return
	}

	if err := validatePassword(body.Password); err != nil {
		slog.ErrorContext(ctx, "invalid password", "error", err)
		httputils.WriteError(w, httputils.BadRequest().WithReasonErr(err))
		return
	}

	httputils.WriteJson(w, http.StatusNotImplemented, nil, map[string]string{"error": "not implemented"})
	// TODO: Database.
}
