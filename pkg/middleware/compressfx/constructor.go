package compressfx

func New(p Params) (Result, error) {
	return Result{
		Middleware: Middleware{
			level: p.Config.Level,
		},
	}, nil
}
