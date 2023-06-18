package dbconnectionfx

import (
	"errors"
	"github.com/sicet7/go-compose-it/pkg/dbconfig"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DialectorParams struct {
	fx.In
	Config dbconfig.DatabaseConfig
}

type DialectorResult struct {
	fx.Out
	Dialector gorm.Dialector
}

func NewDialector(p DialectorParams) (DialectorResult, error) {
	switch p.Config.Type() {
	case "sqlite":
		return DialectorResult{
			Dialector: sqlite.Open(p.Config.DSN()),
		}, nil
	case "mssql":
		return DialectorResult{
			Dialector: sqlserver.New(sqlserver.Config{
				DSN: p.Config.DSN(),
			}),
		}, nil
	case "pgsql":
		return DialectorResult{
			Dialector: postgres.New(postgres.Config{
				DSN: p.Config.DSN(),
			}),
		}, nil
	case "mysql":
		return DialectorResult{
			Dialector: mysql.New(mysql.Config{
				DSN: p.Config.DSN(),
			}),
		}, nil
	default:
		return DialectorResult{}, errors.New("unknown or unsupported database connection type")
	}
}
