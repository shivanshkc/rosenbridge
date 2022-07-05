package handlers

import (
	"net/http"

	"github.com/shivanshkc/rosenbridge/src/configs"
	"github.com/shivanshkc/rosenbridge/src/core"
	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// BasicHandler serves the base route, which is usually "/api".
func BasicHandler(writer http.ResponseWriter, req *http.Request) {
	// Prerequisites.
	ctx, conf, log := req.Context(), configs.Get(), logger.Get()

	// Response body.
	body := map[string]interface{}{
		"code":    core.CodeOK,
		"name":    conf.Application.Name,
		"version": conf.Application.Version,
	}

	response := &httputils.ResponseDTO{Status: http.StatusOK, Body: body}
	httputils.WriteAndLog(ctx, writer, response, log)
}
