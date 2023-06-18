package muxfx

import "go.uber.org/fx"

var Module = fx.Module("muxfx",
	fx.Provide(New),
)
