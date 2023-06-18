package trustedproxymiddlewarefx

func New(p Params) (Result, error) {
	return Result{
		Middleware: &Middleware{
			trustedSubnets: p.TrustedSubnets,
		},
	}, nil
}
