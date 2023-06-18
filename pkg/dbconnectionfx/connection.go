package dbconnectionfx

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ConnectionParams struct {
	fx.In
	Dialector gorm.Dialector
	Config    *gorm.Config `optional:"true"`
}

type ConnectionResult struct {
	fx.Out
	Connection *gorm.DB
}

func NewConnection(p ConnectionParams) (ConnectionResult, error) {
	newConn, err := gorm.Open(p.Dialector, p.Config)
	if err != nil {
		return ConnectionResult{}, err
	}
	return ConnectionResult{
		Connection: newConn,
	}, nil
}
