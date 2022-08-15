package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/worldline-go/chore/internal/utils"
	"github.com/worldline-go/chore/models/apimodels"
	"github.com/worldline-go/chore/pkg/registry"
)

// JWTCheck first control token and second control groups and getID.
func JWTCheck(groups []string, getID func(c *fiber.Ctx) string) func(*fiber.Ctx) error {
	var checkGroups map[string]struct{}

	if len(groups) > 0 {
		checkGroups = make(map[string]struct{}, len(groups))
		for _, g := range groups {
			checkGroups[g] = struct{}{}
		}
	}

	return func(c *fiber.Ctx) error {
		if v, _ := c.Locals("skip-middleware-jwt").(bool); v {
			return c.Next()
		}

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

		// set userid, it will use in next handler
		userID, err := uuid.Parse(tokenValues["user"].(string))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				apimodels.Error{
					Error: err.Error(),
				},
			)
		}

		tokenGroups, ok := tokenValues["groups"].([]interface{})
		if !ok {
			tokenGroups = []interface{}{}
		}

		// set locals
		c.Locals("id", userID)
		c.Locals("token", token)
		c.Locals("groups", tokenGroups)

		if checkGroups == nil && getID == nil {
			return c.Next()
		}

		// should be inside at least one group
		if checkGroups != nil {
			for _, t := range tokenGroups {
				ts, ok := t.(string)
				if !ok {
					continue
				}

				if _, ok := checkGroups[ts]; ok {
					return c.Next()
				}
			}
		}

		// if allow true check id requirement
		if getID != nil && checkAllowID(c, getID) {
			return c.Next()
		}

		return c.Status(http.StatusForbidden).JSON(
			apimodels.Error{
				Error: "access denied",
			},
		)
	}
}

func checkAllowID(c *fiber.Ctx, fn func(c *fiber.Ctx) string) bool {
	getID := fn(c)

	// authorization check
	return utils.Allow(c, getID)
}

func IDFromBody(c *fiber.Ctx) string {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return ""
	}

	v, _ := body["id"].(string)

	return v
}

func IDFromQuery(c *fiber.Ctx) string {
	return c.Query("id")
}
