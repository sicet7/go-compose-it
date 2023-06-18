package accesslogmiddlewarefx

import "go.uber.org/fx"

var Module = fx.Module("accesslogmiddlewarefx",
	fx.Provide(New),
)
