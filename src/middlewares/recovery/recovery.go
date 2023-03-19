package recovery

import (
	"github.com/rs/zerolog"
	"net/http"
)

type RecoveryMiddleware struct {
	logger *zerolog.Logger
}

func NewRecoveryMiddleware(logger *zerolog.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		logger: logger,
	}
}

func (m *RecoveryMiddleware) Mount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				m.logger.Error().Str("type", "panic-log").Msgf("unhandled panic %v", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (*RecoveryMiddleware) Priority() int {
	return 100000
}
