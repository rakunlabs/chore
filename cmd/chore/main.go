package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/cmd/chore/cmd"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/config"
)

func main() {
	config.InitializeLogger()

	if err := cmd.Execute(context.Background()); err != nil {
		log.Error().Err(err).Msg("failed to execute command")
		os.Exit(1)
	}
}
