package database

import (
	"database/sql"
	"errors"
	"github.com/rs/zerolog"
	"github.com/sicet7/go-compose-it/src/config"
	"go.uber.org/fx"
	"golang.org/x/exp/slices"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"strings"
)

var (
	supported = []string{
		"sqlite",
		"mssql",
		"pgsql",
		"mysql",
	}
)

func NewConnection(
	conf *config.Configuration,
	logger *zerolog.Logger,
	lc fx.Lifecycle,
) *gorm.DB {

	parts := strings.SplitN(conf.Database.Url, ":", 2)

	if !slices.Contains(supported, parts[0]) {
		panic(errors.New("unknown or unsupported database type"))
	}

	dbType := parts[0]
	dsn := parts[1]

	scopedLogger := logger.With().Str("type", "database").Logger()

	gormLogging := gormLogger.New(&scopedLogger, gormLogger.Config{})
	var newConn *gorm.DB
	switch dbType {
	case "sqlite":
		dbCon, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: gormLogging,
		})
		if err != nil {
			panic(err)
		}
		newConn = dbCon
		break
	case "mssql":
		dbCon, err := gorm.Open(sqlserver.New(sqlserver.Config{
			DSN:               dsn,
			DefaultStringSize: 256,
		}), &gorm.Config{
			Logger: gormLogging,
		})
		if err != nil {
			panic(err)
		}
		newConn = dbCon
		break
	case "pgsql":
		dbCon, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogging,
		})
		if err != nil {
			panic(err)
		}
		newConn = dbCon
		break
	case "mysql":
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		dbCon, err := gorm.Open(mysql.New(mysql.Config{
			Conn:              sqlDB,
			DefaultStringSize: 256,
		}), &gorm.Config{
			Logger: gormLogging,
		})
		if err != nil {
			panic(err)
		}
		newConn = dbCon
		break
	default:
		panic(errors.New("unknown or unsupported database type"))
	}

	//TODO: Run Migrations here
	//lc.Append(fx.Hook{
	//	OnStart: func(ctx context.Context) error {
	//
	//	},
	//})

	return newConn
}

//func RunMigrations(models map[string]interface{}) error {
//	return Conn().AutoMigrate(maps.Values(models)...)
//}
