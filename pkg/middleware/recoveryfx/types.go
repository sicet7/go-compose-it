package recoveryfx

import (
	"go.uber.org/fx"
	"net/http"
)

type RecoveryHandler interface {
	Handle(err any, w http.ResponseWriter, r *http.Request)
}

type Params struct {
	fx.In
	RecoveryHandler RecoveryHandler
}

type Result struct {
	fx.Out
	Middleware Middleware
}

type Middleware struct {
	handler RecoveryHandler
}
