package compressfx

import "go.uber.org/fx"

var Module = fx.Module("compressfx",
	fx.Provide(New),
)
