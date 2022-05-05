package fn

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetQueryBool(c *fiber.Ctx, s string) (bool, error) {
	rest := false
	restRaw := c.Query(s)

	if restRaw != "" {
		var err error

		rest, err = strconv.ParseBool(restRaw)
		if err != nil {
			return false, err
		}
	}

	return rest, nil
}
