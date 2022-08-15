package args

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/worldline-go/chore/internal/config"
	"github.com/worldline-go/chore/internal/server"
	"github.com/worldline-go/chore/internal/store"

	// Add flow nodes to register in control flow algorithm.
	_ "github.com/worldline-go/chore/pkg/flow/nodes"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/worldline-go/igconfig"
	"github.com/worldline-go/igconfig/loader"
	"github.com/worldline-go/logz"
)

type overrideHold struct {
	Memory *string
	Value  string
}

var rootCmd = &cobra.Command{
	Use:     "chore",
	Short:   "control flow runner",
	Long:    config.Banner("request with templates"),
	Version: config.AppVersion,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := logz.SetLogLevel(config.Application.LogLevel); err != nil {
			return err //nolint:wrapcheck // no need
		}

		return nil
	},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// load configuration
		if err := loadConfig(cmd.Context(), cmd.Flags().Visit); err != nil {
			return err
		}

		if err := runRoot(cmd.Context()); err != nil {
			return err
		}

		return nil
	},
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx) //nolint:wrapcheck // no need
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

func loadConfig(ctx context.Context, visit func(fn func(*pflag.Flag))) error {
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

	loader.VaultSecretAdditionalPaths = append(loader.VaultSecretAdditionalPaths,
		loader.AdditionalPath{Map: "migrate", Name: "migrate"},
		loader.AdditionalPath{Map: "migrate", Name: "migrations"},
	)

	if err := igconfig.LoadWithLoadersWithContext(ctxConfig, "", &config.LoadConfig, loaders[3]); err != nil {
		return fmt.Errorf("unable to load prefix settings: %v", err)
	}

	loader.ConsulConfigPathPrefix = config.LoadConfig.Prefix.Consul
	loader.VaultSecretBasePath = config.LoadConfig.Prefix.Vault

	if err := igconfig.LoadWithLoadersWithContext(ctxConfig, config.LoadConfig.AppName, &config.Application, loaders...); err != nil {
		return fmt.Errorf("unable to load configuration settings: %v", err)
	}

	// override used cmd values
	visit(func(f *pflag.Flag) {
		if v, ok := overrideValues[f.Name]; ok {
			*v.Memory = v.Value
		}
	})

	// set log again to get changes
	if err := logz.SetLogLevel(config.Application.LogLevel); err != nil {
		return err //nolint:wrapcheck // no need
	}

	// print loaded object
	log.Info().Object("config", igconfig.Printer{Value: config.Application}).Msg("loaded config")

	return nil
}

func runRoot(ctxParent context.Context) (err error) {
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

		if errShutdown := server.Shutdown(); errShutdown != nil {
			// set error
			err = errShutdown
		}

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
		"timeZone": config.Application.Store.TimeZone,
		"dsn":      config.Application.Store.DBDataSource,
	})
	if err != nil {
		return fmt.Errorf("cannot open db: %v", err)
	}

	// get generic interface and close in defer
	dbGeneric, err := dbConn.DB()
	if err != nil {
		return fmt.Errorf("cannot get generic interface of gorm: %v", err)
	}

	defer dbGeneric.Close()

	// migrate database
	if err := store.AutoMigrate(ctx, dbConn); err != nil {
		return fmt.Errorf("auto migration: %v", err)
	}

	// server wait
	if err := server.Serve(ctx, "main", dbConn); err != nil {
		return err //nolint:wrapcheck // no need
	}

	return nil
}
