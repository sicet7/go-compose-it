package recoveryfx

import "go.uber.org/fx"

var Module = fx.Module("recoveryfx",
	fx.Provide(New),
)
