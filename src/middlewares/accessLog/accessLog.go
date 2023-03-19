package accessLog

import (
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type AccessLogMiddleware struct {
	logger *zerolog.Logger
}

func (m *AccessLogMiddleware) Mount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		newWriter := responseWrite{internalWriter: w, code: 200}
		next.ServeHTTP(&newWriter, r)
		duration := time.Since(startTime)
		m.logger.Info().
			Str("type", "access-log").
			Str("host", r.Host).
			Int("code", newWriter.Code()).
			Str("user-agent", r.UserAgent()).
			Int64("duration", duration.Milliseconds()).
			Str("ip", r.RemoteAddr).
			Int64("content-length", r.ContentLength).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Msg(http.StatusText(newWriter.Code()))
	})
}

func (*AccessLogMiddleware) Priority() int {
	return 10000
}

func NewAccessLogMiddleware(logger *zerolog.Logger) *AccessLogMiddleware {
	return &AccessLogMiddleware{
		logger: logger,
	}
}

type responseWrite struct {
	code           int
	internalWriter http.ResponseWriter
}

func (r *responseWrite) Header() http.Header {
	return r.internalWriter.Header()
}

func (r *responseWrite) Write(b []byte) (int, error) {
	return r.internalWriter.Write(b)
}

func (r *responseWrite) WriteHeader(statusCode int) {
	r.code = statusCode
	r.internalWriter.WriteHeader(statusCode)
}

func (r *responseWrite) Code() int {
	return r.code
}
