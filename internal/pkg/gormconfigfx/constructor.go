package gormconfigfx

import "gorm.io/gorm"

func New(p Params) (Result, error) {
	return Result{
		Config: &gorm.Config{
			Logger: p.Logger,
		},
	}, nil
}
