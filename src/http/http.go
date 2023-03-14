package http

import (
	"github.com/sicet7/go-compose-it/src/http/middleware"
	"go.uber.org/fx"
	"net/http"
)

type Route interface {
	http.Handler

	Pattern() string
	Middleware() middleware.Middleware
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
		mux.Handle(route.Pattern(), route.Middleware().Mount(route))
	}
	return mux
}
