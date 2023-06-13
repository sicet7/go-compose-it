package accesslogfx

import "go.uber.org/fx"

var Module = fx.Module("accesslogfx",
	fx.Provide(New),
)
