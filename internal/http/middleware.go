package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/shivanshkc/rosenbridge/v3/pkg/logger"
	"github.com/shivanshkc/rosenbridge/v3/pkg/utils/ctxutils"
)

// Middleware implements all the REST middleware methods.
type Middleware struct {
	Logger *logger.Logger
}

// Recovery is a panic recovery middleware.
func (m *Middleware) Recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           func(c echo.Context) bool { return false },
		StackSize:         middleware.DefaultRecoverConfig.StackSize,
		DisableStackAll:   false,
		DisablePrintStack: false,
		// This allows the usage of our custom logger.
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log := m.Logger.ForContext(c.Request().Context())
			log.Error().Err(err).Bytes("stack", stack).Msg("")
			return err
		},
	})(next)
}

// CORS is a Cross-Origin Resource Sharing (CORS) middleware.
func (m *Middleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          func(c echo.Context) bool { return false },
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"*"},
	})(next)
}

// Secure defends against cross-site scripting (XSS) attack, content type sniffing, clickjacking,
// insecure connection and other code injection attacks.
func (m *Middleware) Secure(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.Secure()(next)
}

// AccessLogger middleware handles access logging.
func (m *Middleware) AccessLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		// This will be used to calculate the total request execution time.
		start := time.Now()
		// Shorthand for the underlying request.
		req := eCtx.Request()

		// Setup the request's context.
		ctxutils.SetRequestCtxInfo(req)
		// Fetch the logger for the updated request context.
		log := m.Logger.ForContext(req.Context())

		// Embedding the writer into the custom-writer to persist status-code for logging.
		cWriter := &responseWriterWithCode{ResponseWriter: eCtx.Response()}
		// Update the underlying response writer.
		eCtx.SetResponse(echo.NewResponse(cWriter, eCtx.Echo()))

		// Request entry log.
		log.Info().Str("method", req.Method).Str("url", req.URL.String()).
			Msg("request received")

		// Release control to the next middleware or handler.
		err := next(eCtx)

		// Request exit log.
		log.Info().Int("code", cWriter.statusCode).Int64("latency", int64(time.Since(start))).
			Msg("request completed")

		return err
	}
}

// responseWriterWithCode is a wrapper for http.ResponseWriter for persisting statusCode.
type responseWriterWithCode struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriterWithCode) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
