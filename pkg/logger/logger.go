package logger

import (
	"github.com/rs/zerolog"
	"go-compose-it/pkg/config"
	"os"
	"time"
)

var (
	loggerList = map[string]zerolog.Logger{}
)

func makeRootLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(config.Get().LogLevel)
	var logFile *os.File
	logFile = os.Stdout
	if config.Get().LogFile != "" {
		openLogFile, err := os.OpenFile(config.Get().LogFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err == nil {
			logFile = openLogFile
		}
	}
	return zerolog.New(logFile).With().Timestamp().Logger()
}

func getLogger(name string) *zerolog.Logger {
	rootLogger, exist := loggerList["root"]
	if !exist {
		rootLogger = makeRootLogger()
		loggerList["root"] = rootLogger
	}
	namedLogger, namedExists := loggerList[name]
	if !namedExists {
		namedLogger = rootLogger.With().Str("logger", name).Logger()
		loggerList[name] = namedLogger
	}
	return &namedLogger
}

func Get(name string) *zerolog.Logger {
	return getLogger(name)
}
