package muxfx

import "net/http"

func New(p Params) (Result, error) {
	mux := http.NewServeMux()
	for _, route := range p.Routes {
		mux.Handle(route.Pattern(), p.GlobalMiddleware.Mount(route.Middleware().Mount(route)))
	}
	return Result{
		Mux: mux,
	}, nil
}
