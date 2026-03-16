package rest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/shivanshkc/rosenbridge/internal/database"
	"github.com/shivanshkc/rosenbridge/pkg/utils/httputils"

	"golang.org/x/crypto/bcrypt"
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

	// Hash password.
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.ErrorContext(ctx, "failed to hash password", "error", err)
		httputils.WriteError(w, httputils.InternalServerError())
		return
	}

	// Insert in database.
	user := database.User{Username: body.Username, PasswordHash: string(passwordHashBytes)}
	if err := h.dbase.InsertUser(ctx, user); err != nil {
		if errors.Is(err, database.ErrUserAlreadyExists) {
			slog.ErrorContext(ctx, "user already exists", "error", err)
			httputils.WriteError(w, httputils.Conflict().WithReasonStr("user already exists"))
			return
		}
		slog.ErrorContext(ctx, "unexpected error in user insertion", "error", err)
		httputils.WriteError(w, httputils.InternalServerError())
		return
	}

	httputils.WriteJson(w, http.StatusCreated, nil, map[string]string{"username": user.Username})
}
