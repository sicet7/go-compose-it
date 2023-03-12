package logger

import (
	"github.com/rs/zerolog"
	"github.com/sicet7/go-compose-it/src/config"
	"os"
	"time"
)

func NewLogger(conf *config.Configuration) *zerolog.Logger {
	var logFile *os.File
	logFile = os.Stdout
	if conf.Log.File != "" {
		openLogFile, err := os.OpenFile(conf.Log.File, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		logFile = openLogFile
	}
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(conf.Log.Level())
	logger := zerolog.New(logFile).With().Timestamp().Logger()
	return &logger
}
