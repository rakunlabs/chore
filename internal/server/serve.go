package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/config"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/registry"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/request"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/sec"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/translate"
)

var timeOut = 5 * time.Second

func Serve(ctx context.Context, name string, DB *gorm.DB) {
	app := fiber.New(fiber.Config{
		AppName:               config.AppName,
		DisableStartupMessage: true,
		ReadTimeout:           timeOut,
		WriteTimeout:          timeOut,
	})

	appStore := &registry.AppStore{
		DB:       DB,
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
	app.Use(compress.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("registry", name)
		c.SetUserContext(ctx)

		return c.Next() //nolint:wrapcheck // not need
	})

	setHandlers(app)
	setFileHandler(app)

	// 404
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
		log.Error().Err(err).Msg("server cannot start")
	}
}

func Shutdown() {
	registry.Reg().Iter(func(reg *registry.AppStore) {
		if err := reg.App.Shutdown(); err != nil {
			log.Error().Err(err).Msg("failed to shutdown app")
		}
	})
}
