package muxfx

import (
	"go.uber.org/fx"
	"net/http"
)

type MiddlewareStack struct {
	stack []Middleware
}

func NewMiddlewareStack() MiddlewareStack {
	return MiddlewareStack{}
}

func (s *MiddlewareStack) Add(middleware Middleware) *MiddlewareStack {
	s.stack = append(s.stack, middleware)
	return s
}

func (s *MiddlewareStack) Mount(next http.Handler) http.Handler {
	output := next
	for _, middleware := range s.stack {
		output = middleware.Mount(output)
	}
	return output
}

func AsGlobalMiddleware(f any, paramTags ...string) any {
	if len(paramTags) > 0 {
		return fx.Annotate(
			f,
			fx.ParamTags(paramTags...),
			fx.As(new(Middleware)),
			fx.ResultTags(`name:"global-middleware"`),
		)
	}
	return fx.Annotate(
		f,
		fx.As(new(Middleware)),
		fx.ResultTags(`name:"global-middleware"`),
	)
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
