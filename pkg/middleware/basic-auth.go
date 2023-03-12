package middleware

import (
	"context"
	"net/http"
)

type BasicAuthUser interface {
	RequestTag() string
	VerifyPassword(string) bool
}

type BasicAuthUserProvider interface {
	FindUserByUsername(string) (BasicAuthUser, error)
}

func BasicAuthMiddleware(next http.Handler, provider BasicAuthUserProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			user, err := provider.FindUserByUsername(username)
			if err == nil && user.VerifyPassword(password) {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user.RequestTag())))
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
