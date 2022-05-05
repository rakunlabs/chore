package fn

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func JSON(c *fiber.Ctx, data interface{}, pretty bool) error {
	var (
		raw []byte
		err error
	)

	if pretty {
		raw, err = json.MarshalIndent(data, "", "  ")
	} else {
		raw, err = c.App().Config().JSONEncoder(data)
	}

	if err != nil {
		return err
	}

	c.Context().Response.SetBodyRaw(raw)
	c.Context().Response.Header.SetContentType(fiber.MIMEApplicationJSON)

	return nil
}
