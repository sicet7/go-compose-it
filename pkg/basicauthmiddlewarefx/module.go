package basicauthmiddlewarefx

import "go.uber.org/fx"

var Module = fx.Module("basicauthmiddlewarefx",
	fx.Provide(New),
)
