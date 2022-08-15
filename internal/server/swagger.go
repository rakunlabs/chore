package server

import (
	"path"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/swagger"

	"github.com/worldline-go/chore/docs"
	"github.com/worldline-go/chore/internal/config"
)

func routerSwagger(f fiber.Router) {
	// information
	docs.SwaggerInfo.Title = config.AppName
	docs.SwaggerInfo.Version = config.AppVersion

	docs.SwaggerInfo.BasePath = path.Join(config.Application.BasePath, docs.SwaggerInfo.BasePath)

	// swagger documentation
	f.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("./swagger/index.html") //nolint:wrapcheck
	})

	f.Get("/swagger/*", swagger.HandlerDefault) // default
}
