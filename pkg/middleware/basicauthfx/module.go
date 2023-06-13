package basicauthfx

import "go.uber.org/fx"

var Module = fx.Module("basicauthfx",
	fx.Provide(New),
)
