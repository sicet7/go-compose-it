package accesslogmiddlewarefx

import (
	"go.uber.org/fx"
	"net/http"
	"time"
)

type Params struct {
	fx.In
	Handler AccessLogHandler
}

type Result struct {
	fx.Out
	Middleware *Middleware
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
