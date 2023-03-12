package middleware

import (
	"github.com/rs/zerolog"
	"net/http"
)

func RecoveryMiddleware(next http.Handler, logger *zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error().Msgf("unhandled panic %v", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
