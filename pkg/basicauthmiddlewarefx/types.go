package basicauthmiddlewarefx

import "go.uber.org/fx"

type Params struct {
	fx.In
	Provider UserProvider
}

type Result struct {
	fx.Out
	Middleware *Middleware
}

type Middleware struct {
	provider UserProvider
}

type User interface {
	RequestTag() string
	VerifyPassword(string) bool
}

type UserProvider interface {
	FindUserByUsername(string) (User, error)
}
