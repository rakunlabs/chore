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
	l.Log.Error().Fields(keysAndValues).Msg(msg)
}

func (l LogZ) Info(msg string, keysAndValues ...interface{}) {
	l.Log.Info().Fields(keysAndValues).Msg(msg)
}

func (l LogZ) Debug(msg string, keysAndValues ...interface{}) {
	l.Log.Debug().Fields(keysAndValues).Msg(msg)
}

func (l LogZ) Warn(msg string, keysAndValues ...interface{}) {
	l.Log.Warn().Fields(keysAndValues).Msg(msg)
}
