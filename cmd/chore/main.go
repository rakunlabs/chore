package main

import (
	"context"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/logz"

	"github.com/worldline-go/chore/cmd/chore/args"
)

func main() {
	logz.InitializeLog(logz.WithCaller(false))

	if err := args.Execute(context.Background()); err != nil {
		if !errors.Is(err, args.ErrShutdown) {
			log.Error().Err(err).Msg("failed to execute command")
		}
		os.Exit(1)
	}
}
