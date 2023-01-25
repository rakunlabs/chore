package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog/log"
	"github.com/rytsh/liz/utils/templatex"
	"gorm.io/gorm"

	"github.com/worldline-go/chore/internal/config"
	"github.com/worldline-go/chore/models/apimodels"
	"github.com/worldline-go/chore/pkg/registry"
	"github.com/worldline-go/chore/pkg/sec"
)

var shutdownTimeout = 5 * time.Second

func Serve(ctx context.Context, wg *sync.WaitGroup, name string, db *gorm.DB) error {
	app := fiber.New(fiber.Config{
		AppName:               config.AppName,
		DisableStartupMessage: true,
		ReadBufferSize:        config.Application.Server.ReadBufferSize,
		WriteBufferSize:       config.Application.Server.WriteBufferSize,
	})

	appStore := &registry.AppStore{
		DB:       db,
		Template: templatex.New(),
		App:      app,
		JWT: sec.NewJWT(
			[]byte(config.Application.Secret),
			func() int64 {
				return time.Now().Add(time.Hour).Unix()
			},
		),
	}

	registry.Reg(registry.WithWaitGroup(wg)).Set(name, appStore)

	// middlewares
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(cors.New())

	// compression for gzip, deflate, brotli
	app.Use(compress.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("registry", name)

		return c.Next() //nolint:wrapcheck // not need
	})

	app.Use(requestid.New(), func(c *fiber.Ctx) error {
		// set request id to logger context
		requestID := c.Locals("requestid").(string)
		logRequest := log.With().Str("requestid", requestID).Logger()
		logCtx := logRequest.WithContext(ctx)

		c.SetUserContext(logCtx)

		return c.Next() //nolint:wrapcheck // not need
	})

	if config.Application.BasePath != "" {
		log.Info().Msgf("application BasePath: %s", config.Application.BasePath)
	}

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
		if err := reg.App.ShutdownWithTimeout(shutdownTimeout); err != nil {
			return fmt.Errorf("failed to shutdown app: %v", err)
		}
	}

	return nil
}
