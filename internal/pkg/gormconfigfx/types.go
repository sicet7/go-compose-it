package gormconfigfx

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Params struct {
	fx.In
	Logger logger.Interface
}

type Result struct {
	fx.Out
	Config *gorm.Config
}
