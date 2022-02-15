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
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/store"
	"gitlab.test.igdcs.com/finops/nextgen/utils/basics/igconfig.git/v2"
	"gitlab.test.igdcs.com/finops/nextgen/utils/basics/igconfig.git/v2/loader"

	// Add flow nodes to register in control flow algorithm.
	_ "gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow/nodes"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type overrideHold struct {
	Memory *string
	Value  string
}

var rootCmd = &cobra.Command{
	Use:     "chore",
	Short:   "custom request sender",
	Long:    config.Banner("request with templates"),
	Version: config.AppVersion,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.SetLogLevel(config.Application.LogLevel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// load configuration
		loadConfig(cmd.Context(), cmd.Flags().Visit)

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

// override function hold first values of definitions.
// Use with pflag visit function.
func override(ow map[string]overrideHold) {
	ow["host"] = overrideHold{&config.Application.Host, config.Application.Host}
	ow["port"] = overrideHold{&config.Application.Port, config.Application.Port}
	ow["log-level"] = overrideHold{&config.Application.LogLevel, config.Application.LogLevel}
}

func loadConfig(ctx context.Context, visit func(fn func(*pflag.Flag))) {
	overrideValues := make(map[string]overrideHold)
	override(overrideValues)

	logConfig := log.With().Str("component", "config").Logger()
	ctxConfig := logConfig.WithContext(ctx)

	loaders := []loader.Loader{
		&loader.Consul{},
		&loader.Vault{},
		&loader.File{},
		&loader.Env{},
	}

	if err := igconfig.LoadWithLoadersWithContext(ctxConfig, config.GetLoadName(), &config.Application, loaders...); err != nil {
		log.Error().Err(err).Msg("unable to load configuration settings")

		return
	}

	// override used cmd values
	visit(func(f *pflag.Flag) {
		if v, ok := overrideValues[f.Name]; ok {
			*v.Memory = v.Value
		}
	})

	// set log again to get changes
	config.SetLogLevel(config.Application.LogLevel)

	// print loaded object
	log.Info().EmbedObject(igconfig.Printer{Value: config.Application}).Msg("loaded config")
}

func runRoot(ctxParent context.Context) {
	// appname and version
	// config.Banner("custom send request with templates"),
	log.Info().Msgf("%s [%s]", strings.ToTitle(config.AppName), config.AppVersion)

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
			log.Info().Msg("Gracefully shutting down...")
		}

		server.Shutdown()

		wg.Done()
	}()

	defer func() {
		close(c)
		wg.Wait()
	}()

	// open db connection
	dbConn, err := store.OpenConnection(config.Application.Store.Type, map[string]interface{}{
		"host":     config.Application.Store.Host,
		"port":     config.Application.Store.Port,
		"password": config.Application.Store.Password,
		"user":     config.Application.Store.User,
		"dbName":   config.Application.Store.DBName,
		"schema":   config.Application.Store.Schema,
	})
	if err != nil {
		log.Error().Err(err).Msg("cannot open db")

		return
	}

	// get generic interface and close in defer
	dbGeneric, err := dbConn.DB()
	if err != nil {
		log.Error().Err(err).Msg("cannot get generic interface of gorm")

		return
	}

	defer dbGeneric.Close()

	// migrate database
	if err := store.AutoMigrate(ctx, dbConn); err != nil {
		log.Error().Err(err).Msg("auto migration")

		return
	}

	// server wait
	server.Serve(ctx, "main", dbConn)
}
