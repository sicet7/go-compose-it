package accesslogmiddlewarefx

import (
	"net/http"
	"time"
)

func (m *Middleware) Mount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		newWriter := responseWrite{internalWriter: w, code: 200}
		next.ServeHTTP(&newWriter, r)
		duration := time.Since(startTime)
		m.handler.LogAction(AccessLogAction{
			ResponseCode: newWriter.Code(),
			Request:      r,
			Duration:     duration,
		})
	})
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
