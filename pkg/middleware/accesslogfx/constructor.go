package accesslogfx

func New(p Params) (Result, error) {
	return Result{
		Middleware: Middleware{
			handler: p.Handler,
		},
	}, nil
}
