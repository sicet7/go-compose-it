package gormconfigfx

import "go.uber.org/fx"

var Module = fx.Module("gormconfigfx",
	fx.Provide(New),
)
