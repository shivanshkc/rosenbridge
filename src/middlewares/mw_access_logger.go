package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shivanshkc/rosenbridge/src/logger"
	"github.com/shivanshkc/rosenbridge/src/utils/ctxutils"
	"github.com/shivanshkc/rosenbridge/src/utils/errutils"
	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

// AccessLogger logs the details of the incoming requests and outgoing responses.
func AccessLogger(next http.Handler) http.Handler {
	log := logger.Get()

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		ctxData := ctxutils.GetRequestContextData(ctx)
		// An unlikely scenario, but no harm in checking.
		if ctxData == nil {
			log.Error(ctx, &logger.Entry{Payload: "No context data found for request"})
			err := errutils.InternalServerError().AddMessages("no context found")
			httputils.WriteErrAndLog(ctx, writer, err, log)
			return
		}

		// Getting client's IP address for logging.
		clientIP, err := httputils.GetClientIP(request)
		if err != nil {
			log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("failed to get client IP address: %w", err)})
			clientIP = "unknown"
		}

		// Getting own IP address (IP address of this machine) for logging.
		ownIP, err := httputils.GetOwnIP()
		if err != nil {
			log.Error(ctx, &logger.Entry{Payload: fmt.Errorf("failed to get own IP address: %w", err)})
			ownIP = "unknown"
		}

		// Logging request entry.
		log.Info(ctx, &logger.Entry{
			Timestamp: ctxData.EntryTime,
			Payload:   fmt.Sprintf("ENTRY: Request: %s from: %s", ctxData.ID, clientIP),
			Request: &logger.NetworkRequest{
				Protocol:     "http",
				ID:           ctxData.ID,
				Method:       request.Method,
				URL:          request.URL.Path,
				RequestSize:  request.ContentLength,
				ResponseSize: 0,
				Latency:      0,
				ServerIP:     ownIP,
				ClientIP:     clientIP,
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
			Payload:   fmt.Sprintf("EXIT: Request: %s from: %s", ctxData.ID, clientIP),
			Request: &logger.NetworkRequest{
				Status:       customWriter.statusCode,
				Protocol:     "http",
				ID:           ctxData.ID,
				Method:       request.Method,
				URL:          request.URL.Path,
				RequestSize:  request.ContentLength,
				ResponseSize: resContentLength,
				Latency:      time.Since(ctxData.EntryTime),
				ServerIP:     ownIP,
				ClientIP:     clientIP,
			},
		})
	})
}
