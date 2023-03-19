package http

import (
	"go.uber.org/fx"
	"net/http"
)

type Route interface {
	http.Handler

	Pattern() string
	Middleware() Middleware
}

func AsRoute(f any, paramTags ...string) any {
	if len(paramTags) > 0 {
		return fx.Annotate(
			f,
			fx.ParamTags(paramTags...),
			fx.As(new(Route)),
			fx.ResultTags(`group:"routes"`),
		)
	}
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

func NewServeMux(routes []Route, globalMiddleware Middleware) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), globalMiddleware.Mount(route.Middleware().Mount(route)))
	}
	return mux
}
