package handler

import (
	"github.com/gofiber/fiber/v2"
)

// Send request
// @Summary Send request
// @Description Send request to api.
// @Accept */*
// @Success 201 {object} map[string]interface{}
// @Router /send [get]
func Send(c *fiber.Ctx) error {
	_, err := c.WriteString("Hellooo from send")

	return err
}
