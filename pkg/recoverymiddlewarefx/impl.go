package recoverymiddlewarefx

import "net/http"

func (m *Middleware) Mount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.handler.Handle(err, w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
