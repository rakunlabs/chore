package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/worldline-go/auth/pkg/authecho"
	"github.com/worldline-go/chore/internal/server/claims"
	"github.com/worldline-go/chore/models/apimodels"
)

var (
	AdminRoleKey = "chore_admin"
	UserRoleKey  = "chore_user"

	AdminRole = authecho.MiddlewareRole(authecho.WithRoles("chore_admin"))
	UserRole  = authecho.MiddlewareRole(authecho.WithRoles("chore_user"))
)

// JWTCheck comes after the auth middleware to check ID's can have do that.
func JWTCheck(getID func(c echo.Context) string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// noop authetication skip
			if v, ok := c.Get(authecho.KeyAuthNoop).(bool); ok && v {
				return next(c)
			}

			id := getID(c)
			if id == "" {
				return c.JSON(http.StatusForbidden, apimodels.Error{Error: "id not found"})
			}

			// get token ID
			claim, ok := c.Get(authecho.KeyClaims).(*claims.Custom)
			if !ok {
				return c.JSON(http.StatusForbidden, apimodels.Error{Error: "claims not found"})
			}

			if claim.Subject == "" {
				return c.JSON(http.StatusForbidden, apimodels.Error{Error: "claims user id not found"})
			}

			if id == claim.Subject {
				// set disables
				c.Set(authecho.DisableRoleCheckKey, true)
				c.Set(authecho.DisableScopeCheckKey, true)
			}

			return next(c)
		}
	}
}

func IDFromBody(c echo.Context) string {
	body := struct {
		ID string `json:"id"`
	}{}
	if err := c.Bind(&body); err != nil {
		return ""
	}

	return body.ID
}

func IDFromQuery(c echo.Context) string {
	return c.QueryParam("id")
}