package database

import (
	"database/sql"
	"errors"
	"github.com/rs/zerolog"
	"go-compose-it/pkg/env"
	myLogger "go-compose-it/pkg/logger"
	"golang.org/x/exp/maps"
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
	logger     zerolog.Logger
	connection *gorm.DB
	supported  = []string{
		"sqlite",
		"mssql",
		"pgsql",
		"mysql",
	}
)

func init() {
	logger = myLogger.Get("database")
	newConn, err := newConnection(env.Get().DatabaseUrl)

	if err != nil {
		logger.Fatal().Msgf("failed to connect to database: %v", err)
	}
	connection = newConn
}

func Conn() *gorm.DB {
	return connection
}

func newConnection(databaseUrl string) (*gorm.DB, error) {

	parts := strings.SplitN(databaseUrl, ":", 2)

	if !slices.Contains(supported, parts[0]) {
		return nil, errors.New("unknown or unsupported database type")
	}

	dbType := parts[0]

	dsn := parts[1]

	gormLogging := gormLogger.New(&logger, gormLogger.Config{})
	var newConn *gorm.DB
	switch dbType {
	case "sqlite":
		dbCon, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: gormLogging,
		})
		if err != nil {
			return nil, err
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
			return nil, err
		}
		newConn = dbCon
		break
	case "pgsql":
		dbCon, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogging,
		})
		if err != nil {
			return nil, err
		}
		newConn = dbCon
		break
	case "mysql":
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}
		dbCon, err := gorm.Open(mysql.New(mysql.Config{
			Conn:              sqlDB,
			DefaultStringSize: 256,
		}), &gorm.Config{
			Logger: gormLogging,
		})
		if err != nil {
			return nil, err
		}
		newConn = dbCon
		break
	default:
		return nil, errors.New("unknown or unsupported database type")
	}

	return newConn, nil
}

func RunMigrations(models map[string]interface{}) error {
	return Conn().AutoMigrate(maps.Values(models)...)
}
