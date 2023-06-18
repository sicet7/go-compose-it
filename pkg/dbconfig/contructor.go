package dbconfig

import (
	"errors"
	"golang.org/x/exp/slices"
	"strings"
)

func New(databaseUrl string) (DatabaseConfig, error) {
	parts := strings.SplitN(databaseUrl, ":", 2)

	if !slices.Contains(supported, parts[0]) {
		return DatabaseConfig{}, errors.New("unknown or unsupported database type")
	}

	return DatabaseConfig{
		typeName: parts[0],
		dsn:      parts[1],
	}, nil
}
