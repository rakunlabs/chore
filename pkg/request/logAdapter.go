package request

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog"
)

type LogZ struct {
	Log zerolog.Logger
}

var _ retryablehttp.LeveledLogger = LogZ{}

func (l LogZ) Error(msg string, keysAndValues ...interface{}) {
	l.Log.Error().Interface("value", keysAndValues).Msg(msg)
}

func (l LogZ) Info(msg string, keysAndValues ...interface{}) {
	l.Log.Info().Interface("value", keysAndValues).Msg(msg)
}

func (l LogZ) Debug(msg string, keysAndValues ...interface{}) {
	l.Log.Debug().Interface("value", keysAndValues).Msg(msg)
}

func (l LogZ) Warn(msg string, keysAndValues ...interface{}) {
	l.Log.Warn().Interface("value", keysAndValues).Msg(msg)
}
