package routes

import (
	"go-compose-it/pkg/actions"
	"go-compose-it/pkg/middleware"
	"net/http"
)

var routes = map[string]http.HandlerFunc{
	"/api/health": actions.HealthAction,
}

func Mount(handler *http.ServeMux) *http.ServeMux {
	for path, route := range routes {
		handler.Handle(path, middleware.Global(route))
	}
	return handler
}
