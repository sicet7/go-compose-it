package routes

import (
	"github.com/sicet7/go-compose-it/pkg/actions"
	"github.com/sicet7/go-compose-it/pkg/config"
	"github.com/sicet7/go-compose-it/pkg/logger"
	"github.com/sicet7/go-compose-it/pkg/middleware"
	"net/http"
)

var routes = map[string]http.HandlerFunc{
	"/api/health": actions.HealthAction,
}

func Mount(handler *http.ServeMux) *http.ServeMux {
	for path, route := range routes {
		handler.Handle(path, mountGlobalMiddlewares(route))
	}
	return handler
}

func mountGlobalMiddlewares(next http.Handler) http.Handler {
	output := next

	output = middleware.ProxyHeadersMiddleware(output, config.Get().Http.Net.GetTrustedProxies())
	output = middleware.CompressionMiddleware(output, 9)
	output = middleware.AccessLogMiddleware(output, logger.Get("http-access"))

	// Recovery middleware should always be the last middleware added
	// this is to make sure all unhandled panics are handled here
	output = middleware.RecoveryMiddleware(output, logger.Get("recovery"))
	return output
}
