package compressionmiddlewarefx

import "go.uber.org/fx"

var Module = fx.Module("compressionmiddlewarefx",
	fx.Provide(New),
)
