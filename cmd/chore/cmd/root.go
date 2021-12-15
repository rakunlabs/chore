package cmd

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/config"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/server"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/store/kv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "chore",
	Short:   "custom request sender",
	Long:    "custom send request with templates",
	Version: config.Application.AppVersion,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.SetLogLevel(config.Application.LogLevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		runRoot(cmd.Context())
	},
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx) //nolint:wrapcheck
}

//nolint:gochecknoinits
func init() {
	rootCmd.Flags().StringVarP(&config.Application.Host, "host", "H", config.Application.Host, "Host to listen on")
	rootCmd.Flags().StringVarP(&config.Application.Port, "port", "P", config.Application.Port, "Port to listen on")
	rootCmd.PersistentFlags().StringVarP(&config.Application.LogLevel, "log-level", "l", config.Application.LogLevel, "Log level")
}

func runRoot(ctxParent context.Context) {
	// appname and version
	log.Logger.Info().Msgf("%s [%s]", strings.ToTitle(config.Application.AppName), config.Application.AppVersion)

	wg := &sync.WaitGroup{}

	ctx, ctxCancel := context.WithCancel(ctxParent)
	defer ctxCancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)

	go func() {
		v := <-c

		ctxCancel()

		if v != nil {
			log.Logger.Info().Msg("Gracefully shutting down...")
		}

		server.Shutdown()

		wg.Done()
	}()

	// get store handler
	storageHandler, err := kv.NewConsul(ctx, os.Getenv("APP_NAME"))
	if err == nil {
		// server wait
		server.Serve("main", storageHandler)
	}

	close(c)
	wg.Wait()
}
