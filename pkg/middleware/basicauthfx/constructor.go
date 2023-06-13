package basicauthfx

func New(p Params) (Result, error) {
	return Result{
		Middleware: Middleware{
			provider: p.Provider,
		},
	}, nil
}
