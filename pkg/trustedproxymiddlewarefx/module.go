package trustedproxymiddlewarefx

import "go.uber.org/fx"

var Module = fx.Module("trustedproxymiddlewarefx",
	fx.Provide(New),
)
