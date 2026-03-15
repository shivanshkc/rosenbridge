package rest

import (
	"log/slog"
	"net/http"

	"github.com/shivanshkc/rosenbridge/pkg/utils/httputils"
)

func (h *Handler) getConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Make sure credentials are correct.
	username, err := h.authenticateUser(r)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	// Upgrade and persist the connection.
	if err := h.wsManager.UpgradeAndAddConnection(w, r, username); err != nil {
		slog.ErrorContext(ctx, "error in UpgradeAndAddConnection call", "error", err)
		// Response is already written.
	}
}
