package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/config"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/request"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/sec"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/translate"
)

func Serve(ctx context.Context, name string, db *gorm.DB) error {
	app := fiber.New(fiber.Config{
		AppName:               config.AppName,
		DisableStartupMessage: true,
		ReadBufferSize:        config.Application.Server.ReadBufferSize,
		WriteBufferSize:       config.Application.Server.WriteBufferSize,
	})

	appStore := &registry.AppStore{
		DB:       db,
		Template: translate.NewTemplate(),
		Client:   request.NewClient(),
		App:      app,
		JWT: sec.NewJWT(
			[]byte(config.Application.Secret),
			func() int64 {
				return time.Now().Add(time.Hour).Unix()
			},
		),
	}

	registry.Reg().Set(name, appStore)

	// middlewares
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(cors.New())

	// compression for gzip, deflate, brotli
	app.Use(compress.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("registry", name)
		c.SetUserContext(ctx)

		return c.Next() //nolint:wrapcheck // not need
	})

	log.Info().Msgf("Application BasePath: %s", config.Application.BasePath)
	appRouter := app.Group(config.Application.BasePath)

	setHandlers(appRouter)

	setFileHandler(appRouter)

	// 404 if not found any handler
	app.Use(func(c *fiber.Ctx) error {
		//nolint: wrapcheck // middleware
		return c.Status(http.StatusNotFound).JSON(apimodels.API{
			Error: apimodels.Error{Error: "404 not found"},
		})
	})

	// custom host
	hostPort := config.Application.Host + ":" + config.Application.Port
	log.Info().Msg("server starting [" + hostPort + "]")

	if err := app.Listen(hostPort); err != nil {
		return fmt.Errorf("server cannot start: %v", err)
	}

	return nil
}

func Shutdown() error {
	reg := registry.Reg().Get("main")

	// check registry exist and server running
	if reg != nil && reg.App.Server() != nil {
		if err := reg.App.Shutdown(); err != nil {
			return fmt.Errorf("failed to shutdown app: %v", err)
		}
	}

	return nil
}
