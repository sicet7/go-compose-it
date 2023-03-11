package env

import "go-compose-it/pkg/utils"

var (
	environment Environment
)

type Environment struct {
	DatabaseUrl   string
	HttpAddress   string
	ShutdownWait  int
	CertFile      string
	KeyFile       string
	AccessLogging bool
}

func init() {
	environment = Environment{
		DatabaseUrl:   utils.ReadStringEnv("DATABASE_URL", "sqlite:data.db"),
		HttpAddress:   utils.ReadStringEnv("HTTP_ADDR", "0.0.0.0:8080"),
		ShutdownWait:  utils.ReadIntEnv("SHUTDOWN_WAIT", 60),
		CertFile:      utils.ReadStringEnv("CERT_FILE", ""),
		KeyFile:       utils.ReadStringEnv("KEY_FILE", ""),
		AccessLogging: utils.ReadBoolEnv("ACCESS_LOGGING", false),
	}
}

func Get() *Environment {
	return &environment
}
