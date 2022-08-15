package api

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/worldline-go/chore/internal/config"
)

type information struct {
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	StartDate  time.Time `json:"startDate"`
	ServerDate time.Time `json:"serverDate"`
}

// @Summary Information
// @Description Get information of the server
// @Tags public
// @Router /info [get]
// @Success 200 {object} information{}
func getInfo(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(
		information{
			Name:       config.AppName,
			Version:    config.AppVersion,
			StartDate:  config.StartDate,
			ServerDate: time.Now(),
		},
	)
}

// @Summary Ping server
// @Description Check server is active
// @Tags public
// @Router /ping [get]
// @Success 200
func getPing(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func Info(router fiber.Router) {
	router.Get("/info", getInfo)
	router.Get("/ping", getPing)
}
