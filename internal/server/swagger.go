package server

import (
	"github.com/gofiber/fiber/v2"

	swagger "github.com/arsmn/fiber-swagger/v2"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/docs"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/config"
)

func routerSwagger(f fiber.Router) {
	// information
	docs.SwaggerInfo.Title = config.AppName
	docs.SwaggerInfo.Version = config.AppVersion

	// swagger documentation
	f.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("./swagger/index.html") //nolint:wrapcheck
	})
	f.Get("/swagger/*", swagger.Handler) // default
}
