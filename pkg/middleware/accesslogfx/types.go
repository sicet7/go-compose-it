package accesslogfx

import (
	"net/http"
	"time"
)

type Params struct {
	Handler AccessLogHandler
}

type Result struct {
	Middleware Middleware
}

type AccessLogAction struct {
	ResponseCode int
	Request      *http.Request
	Duration     time.Duration
}

type AccessLogHandler interface {
	LogAction(action AccessLogAction)
}

type Middleware struct {
	handler AccessLogHandler
}

type responseWrite struct {
	code           int
	internalWriter http.ResponseWriter
}
