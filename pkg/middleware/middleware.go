package middleware

import (
	"net/http"
)

func Global(next http.Handler) http.Handler {
	output := next

	// Recovery middleware should always be the last middleware added
	// this is to make sure all unhandled panics are handled here
	output = RecoveryMiddleware(output)
	return output
}
