package basicauthfx

import (
	"context"
	"net/http"
)

func (m *Middleware) Mount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			user, err := m.provider.FindUserByUsername(username)
			if err == nil && user.VerifyPassword(password) {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user.RequestTag())))
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
