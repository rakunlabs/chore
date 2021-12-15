package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logW = zerolog.ConsoleWriter{
	Out: os.Stderr,
	FormatTimestamp: func(i interface{}) string {
		parse, _ := time.Parse(time.RFC3339, i.(string))

		return parse.Format("2006-01-02 15:04:05")
	},
}

func InitializeLogger() {
	log.Logger = zerolog.New(logW).With().Timestamp().Caller().Logger()
}

func SetLogLevel(level string) {
	zLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Warn().Err(err).Str("component", "log").Msgf("zerolog unknown level %s", level)

		return
	}

	zerolog.SetGlobalLevel(zLevel)
}
