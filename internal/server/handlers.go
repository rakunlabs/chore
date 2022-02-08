package server

import (
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/api"

	"github.com/gofiber/fiber/v2"
)

// @description Storage and Send API
// @description First login with user and use authorization as "Bearer JWTTOKEN"
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func setHandlers(app *fiber.App) {
	apiRouter := app.Group("/api")                                // /api
	v1Router := apiRouter.Group("/v1", func(c *fiber.Ctx) error { // middleware for /api/v1
		c.Set("Version", "v1")

		return c.Next() //nolint:wrapcheck
	})

	// set swagger
	routerSwagger(v1Router)

	// set routers
	api.Auth(v1Router)
	api.Template(v1Router)
	api.User(v1Router)
	api.Login(v1Router)
	api.Token(v1Router)
	api.Bind(v1Router)

	// testing
	// apitest.Test(v1Router)

	// set send api
	api.Send(v1Router)
}
