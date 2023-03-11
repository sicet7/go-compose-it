package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/sicet7/go-compose-it/pkg/config"
	"os"
	"sync"
	"time"
)

var (
	loggerList = map[string]zerolog.Logger{}
	makeLock   = sync.Mutex{}
)

func makeLogger(name string) zerolog.Logger {
	makeLock.Lock()
	defer makeLock.Unlock()
	logger, exist := loggerList[name]
	if exist {
		return logger
	}
	rootLogger, rootExist := loggerList["root"]
	if !rootExist {
		rootLogger = makeRootLogger()
		loggerList["root"] = rootLogger
	}

	if name == "root" {
		return rootLogger
	} else {
		logger = rootLogger.With().Str("logger", name).Logger()
	}
	loggerList[name] = logger
	return logger
}

func makeRootLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(config.Get().GetLogLevel())
	var logFile *os.File
	logFile = os.Stdout
	if config.Get().Log.File != "" {
		openLogFile, err := os.OpenFile(config.Get().Log.File, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("failed to open log file: \"%s\"\n", config.Get().Log.File)
			os.Exit(1)
		}
		logFile = openLogFile
	}
	return zerolog.New(logFile).With().Timestamp().Logger()
}

func getLogger(name string) *zerolog.Logger {
	namedLogger, namedExists := loggerList[name]
	if !namedExists {
		namedLogger = makeLogger(name)
	}
	return &namedLogger
}

func Get(name string) *zerolog.Logger {
	return getLogger(name)
}
