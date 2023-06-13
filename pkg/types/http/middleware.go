package http

import "net/http"

type Middleware interface {
	Mount(http.Handler) http.Handler
}
