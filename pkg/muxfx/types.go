package muxfx

import (
	"go.uber.org/fx"
	"net/http"
)

type Middleware interface {
	Mount(http.Handler) http.Handler
}

type Route interface {
	http.Handler

	Pattern() string
	Middleware() Middleware
}

type Params struct {
	fx.In
	Routes           []Route    `group:"routes"`
	GlobalMiddleware Middleware `name:"global-middleware"`
}

type Result struct {
	fx.Out
	Mux *http.ServeMux
}
