package recoverymiddlewarefx

func New(p Params) (Result, error) {
	return Result{
		Middleware: &Middleware{
			handler: p.RecoveryHandler,
		},
	}, nil
}
