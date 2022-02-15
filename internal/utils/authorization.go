package utils

import "github.com/gofiber/fiber/v2"

func Allow(c *fiber.Ctx, id string) bool {
	localGroups, _ := c.Locals("groups").([]interface{})

	for _, v := range localGroups {
		if v.(string) == "admin" {
			return true
		}
	}

	localID, _ := c.Locals("id").(string)

	return id == localID
}
