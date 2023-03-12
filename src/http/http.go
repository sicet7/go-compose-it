package http

import (
	"go.uber.org/fx"
	"net/http"
)

type Route interface {
	http.Handler

	Pattern() string
	Middlewares(http.Handler) http.Handler
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route.Middlewares(route))
	}
	return mux
}
