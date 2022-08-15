package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/logz"

	"github.com/worldline-go/chore/cmd/chore/args"
)

func main() {
	logz.InitializeLog(nil)

	if err := args.Execute(context.Background()); err != nil {
		log.Error().Err(err).Msg("failed to execute command")
		os.Exit(1)
	}
}
