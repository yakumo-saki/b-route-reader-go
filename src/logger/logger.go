package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/b-route-reader-go/src/config"
)

func Initiallize() {
	switch config.LOG_LEVEL {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	if config.LOG_NO_DATETIME {
		outputNoTimestamp := zerolog.ConsoleWriter{Out: os.Stdout}
		log.Logger = zerolog.New(outputNoTimestamp).With().Caller().Logger()
	} else {
		zerolog.TimeFieldFormat = time.RFC3339Nano
		outputWithTimestamp := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05.000"}
		log.Logger = zerolog.New(outputWithTimestamp).With().Timestamp().Caller().Logger()
	}
}
