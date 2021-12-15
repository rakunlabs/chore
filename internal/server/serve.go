package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/config"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/server/registry"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/store/inf"
)

var timeOut = 5 * time.Second

func Serve(name string, storeHandler inf.CRUD) {
	app := fiber.New(fiber.Config{
		AppName:               config.Application.AppName,
		DisableStartupMessage: true,
		ReadTimeout:           timeOut,
		WriteTimeout:          timeOut,
	})

	appStore := &registry.AppStore{
		StoreHandler: storeHandler,
		App:          app,
	}

	registry.GetRegistry().Set(name, appStore)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("storeHandler", &(appStore.StoreHandler))

		return c.Next() //nolint:wrapcheck // not need
	})

	setHandlers(app)
	setFileHandler(app)

	// Custom host
	hostPort := config.Application.Host + ":" + config.Application.Port
	log.Logger.Info().Msg("server starting [" + hostPort + "]")

	if err := app.Listen(hostPort); err != nil {
		log.Logger.Error().Err(err).Msg("server cannot start")
	}
}

func Shutdown() {
	registry.GetRegistry().Iter(func(app *fiber.App) {
		if err := app.Shutdown(); err != nil {
			log.Logger.Error().Err(err).Msg("failed to shutdown app")
		}
	})
}
