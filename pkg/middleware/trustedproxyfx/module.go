package trustedproxyfx

import "go.uber.org/fx"

var Module = fx.Module("trustedproxyfx",
	fx.Provide(New),
)
