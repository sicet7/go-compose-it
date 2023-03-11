package middleware

import (
	"github.com/sicet7/go-compose-it/pkg/config"
	"net/http"
)

func Global(next http.Handler) http.Handler {
	output := next

	output = ProxyHeadersMiddleware(output, config.Get().Http.Net.GetTrustedProxies())
	output = CompressionMiddleware(output, 9)
	output = AccessLogMiddleware(output)

	// Recovery middleware should always be the last middleware added
	// this is to make sure all unhandled panics are handled here
	output = RecoveryMiddleware(output)
	return output
}
