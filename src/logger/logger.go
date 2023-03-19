package logger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

type LogConfig interface {
	LogFile() *os.File
	LogLevel() zerolog.Level
}

func NewLogger(conf LogConfig) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(conf.LogLevel())
	logger := zerolog.New(conf.LogFile()).With().Timestamp().Logger()
	return &logger
}
