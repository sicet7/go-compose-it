package recoverymiddlewarefx

import "go.uber.org/fx"

var Module = fx.Module("recoverymiddlewarefx",
	fx.Provide(New),
)
