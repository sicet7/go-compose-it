package middleware

import (
	"github.com/sicet7/go-compose-it/pkg/logger"
	"net/http"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Get("recovery").Error().Msgf("unhandled panic %v", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
