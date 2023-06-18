package trustedproxymiddlewarefx

import (
	"go.uber.org/fx"
	"net"
)

type Params struct {
	fx.In
	TrustedSubnets []net.IPNet `group:"trusted-network-proxies"`
}

type Result struct {
	fx.Out
	Middleware *Middleware
}

type Middleware struct {
	trustedSubnets []net.IPNet
}
