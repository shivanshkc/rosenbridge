package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// AccessLogger logs the details of the incoming requests and outgoing responses.
func AccessLogger(next http.Handler) http.Handler {
	log := logger.Get()

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx, entry := request.Context(), time.Now()
		ctxData := httputils.GetReqCtx(ctx)

		// Getting client's IP address for logging.
		clientIP, err := httputils.GetClientIP(request)
		if err != nil {
			clientIP = "unknown"
		}
		// Getting own IP address (IP address of this machine) for logging.
		ownIP, err := httputils.GetOwnIP()
		if err != nil {
			ownIP = "unknown"
		}

		// Logging request entry.
		log.Info(ctx, &logger.Entry{
			Payload: ">> request in",
			Request: &logger.NetworkRequest{
				Protocol:    "http",
				ID:          ctxData.ID,
				Method:      request.Method,
				URL:         request.URL.Path,
				RequestSize: request.ContentLength,
				ServerIP:    ownIP,
				ClientIP:    clientIP,
			},
		})

		// Wrapping the writer with a custom writer for persisting statusCode.
		customWriter := &responseWriterWithCode{ResponseWriter: writer}
		// Releasing the control to following middlewares/handlers.
		next.ServeHTTP(customWriter, request)

		// Calculating response content length.
		var resContentLength int64
		if contentLength := writer.Header().Get("content-length"); contentLength != "" {
			resContentLength, _ = strconv.ParseInt(contentLength, 10, 64)
		}

		// Logging response exit.
		log.Info(ctx, &logger.Entry{
			Timestamp: time.Now(),
			Payload:   "<< request out",
			Request: &logger.NetworkRequest{
				Status:       customWriter.statusCode,
				Protocol:     "http",
				ID:           ctxData.ID,
				Method:       request.Method,
				URL:          request.URL.Path,
				RequestSize:  request.ContentLength,
				ResponseSize: resContentLength,
				Latency:      time.Since(entry),
				ServerIP:     ownIP,
				ClientIP:     clientIP,
			},
		})
	})
}
