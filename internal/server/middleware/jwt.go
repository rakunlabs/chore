package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/registry"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
)

func JWTCheck(allow string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		JWT := registry.Reg().Get(c.Locals("registry").(string)).JWT
		if JWT == nil {
			return c.Next()
		}

		var token string

		headerValue := c.Get("authorization")

		if headerValue == "" {
			return c.Status(http.StatusForbidden).JSON(
				apimodels.Error{
					Error: "forbidden authorization header not found",
				},
			)
		}

		components := strings.SplitN(headerValue, " ", 2)

		if len(components) != 2 || !strings.EqualFold(components[0], "Bearer") {
			return c.Status(http.StatusForbidden).JSON(
				apimodels.Error{
					Error: "forbidden Bearer not found",
				},
			)
		}

		token = components[1]
		if token == "" {
			return c.Status(http.StatusForbidden).JSON(
				apimodels.Error{
					Error: "forbidden token not found",
				},
			)
		}

		tokenValues, err := JWT.Validate(token)
		if err != nil {
			return c.Status(http.StatusForbidden).JSON(
				apimodels.Error{
					Error: err.Error(),
				},
			)
		}

		if check := allow == "admin"; check {
			if tokenAdmin, ok := tokenValues["admin"].(string); ok {
				isTokenAdmin, err := strconv.ParseBool(tokenAdmin)
				if err != nil || isTokenAdmin != check {
					return c.Status(http.StatusForbidden).JSON(
						apimodels.Error{
							Error: "not allowed",
						},
					)
				}
			}
		}

		return c.Next()
	}
}
