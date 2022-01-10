package server

import (
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/server/handler"

	"github.com/gofiber/fiber/v2"
)

// @description storage and send API
// @BasePath /api/v1
func setHandlers(app *fiber.App) {
	api := app.Group("/api")                          // /api
	v1 := api.Group("/v1", func(c *fiber.Ctx) error { // middleware for /api/v1
		c.Set("Version", "v1")

		return c.Next() //nolint:wrapcheck
	})

	// set swagger
	handler.RouterSwagger(v1)

	// send new request
	v1.Post("/send", handler.Send)

	// KV API
	v1KVGroup := v1.Group("/kv")
	handler.RouterKV(v1KVGroup, "/templates")
	handler.RouterKV(v1KVGroup, "/auths")
	handler.RouterKV(v1KVGroup, "/binds")
}
