package dbconnectionfx

import "go.uber.org/fx"

var Module = fx.Module("dbconnectionfx",
	fx.Provide(NewDialector),
	fx.Provide(NewConnection),
)
