package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgconn"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/registry"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/server/middleware"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
)

type TokenRet struct {
	Token string `json:"token" example:"tokenJWT"`
}

// @Summary New token
// @Tags token
// @Description Send and record PAT token
// @Security ApiKeyAuth
// @Router /token [post]
// @Param new query bool false "generate new token"
// @Param payload body TokenRet{} false "send valid token"
// @Success 200 {object} apimodels.Data{data=TokenRet{}}
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postToken(c *fiber.Ctx) error {
	body := new(TokenRet)
	reg := registry.Reg().Get(c.Locals("registry").(string))

	isNew, err := strconv.ParseBool(c.Query("new", "false"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	//nolint: nestif // TODO improve future
	if isNew {
		// generate new PAT token
		var err error
		body.Token, err = reg.JWT.Generate(
			map[string]interface{}{
				"type": models.TypePersonalAccessToken,
			}, 0)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				apimodels.Error{
					Error: err.Error(),
				},
			)
		}
	} else {
		if err := c.BodyParser(body); err != nil {
			return c.Status(http.StatusBadRequest).JSON(
				apimodels.Error{
					Error: err.Error(),
				},
			)
		}

		if body.Token == "" {
			return c.Status(http.StatusBadRequest).JSON(
				apimodels.Error{
					Error: "token is required",
				},
			)
		}

		// check token is valid
		tokenValues, err := reg.JWT.Validate(body.Token)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(
				apimodels.Error{
					Error: fmt.Sprintf("token not valid: %v", err.Error()),
				},
			)
		}

		if t, ok := tokenValues["type"].(string); !ok || t != models.TypePersonalAccessToken {
			return c.Status(http.StatusBadRequest).JSON(
				apimodels.Error{
					Error: "not a PAT type token",
				},
			)
		}
	}

	result := reg.DB.WithContext(c.UserContext()).Create(
		&models.Token{
			TokenPure: models.TokenPure{
				Token: body.Token,
			},
		},
	)

	// check write error
	var pErr *pgconn.PgError

	errors.As(result.Error, &pErr)

	if pErr != nil && pErr.Code == "23505" {
		return c.Status(http.StatusConflict).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	// return recorded data's id
	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: TokenRet{Token: body.Token},
		},
	)
}

// @Summary Delete token
// @Tags token
// @Description Delete with token
// @Security ApiKeyAuth
// @Router /token [delete]
// @Param token query string false "get by token"
// @Success 204 "No Content"
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func deleteToken(c *fiber.Ctx) error {
	tokenQuery := c.Query("token")

	if tokenQuery == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "token is required",
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	// delete directly in DB
	result := reg.DB.WithContext(c.UserContext()).Where("token = ?", tokenQuery).Unscoped().Delete(&models.Token{})

	if result.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(
			apimodels.Error{
				Error: "not found any releated data",
			},
		)
	}

	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	//nolint:wrapcheck // checking before
	return c.SendStatus(http.StatusNoContent)
}

// @Summary Check token
// @Tags token
// @Description Send token to check validation, if not valid return 401
// @Router /token/check [post]
// @Param payload body TokenRet{} false "send a token"
// @Success 200 {object} apimodels.Data{data=TokenRet{}}
// @failure 400 {object} apimodels.Error{}
// @failure 401 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postTokenCheck(c *fiber.Ctx) error {
	body := new(TokenRet)
	reg := registry.Reg().Get(c.Locals("registry").(string))

	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	if body.Token == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "token is required",
			},
		)
	}

	// check token is valid
	_, err := reg.JWT.Validate(body.Token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(
			apimodels.Error{
				Error: fmt.Sprintf("token not valid: %v", err.Error()),
			},
		)
	}

	// return recorded data's id
	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: TokenRet{Token: body.Token},
		},
	)
}

func Token(router fiber.Router) {
	router.Post("/token/check", postTokenCheck)
	router.Post("/token", middleware.JWTCheck(""), postToken)
	router.Delete("/token", middleware.JWTCheck(""), deleteToken)
}
