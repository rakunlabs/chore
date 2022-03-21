package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/server/middleware"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
)

// @Summary Get Settings
// @Tags settings
// @Description Get whole settings
// @Security ApiKeyAuth
// @Router /settings [get]
// @Success 200 {object} apimodels.Data{data=UserDataID{}}
// @failure 500 {object} apimodels.Error{}
func getSettings(c *fiber.Ctx) error {
	settings := new(models.EmailPure)

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.Settings{}).Where("namespace = ?", "application")

	result := query.First(settings)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusOK).JSON(
			apimodels.Data{
				Data: settings,
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

	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: settings,
		},
	)
}

// @Summary Replace settings
// @Tags settings
// @Description Replace with new data
// @Security ApiKeyAuth
// @Router /settings [patch]
// @Param payload body models.Email{} false "send part of the settings object"
// @Success 204 "No Content"
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func patchSettings(c *fiber.Ctx) error {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.Settings{})

	queryUpdate := query.Where("namespace = ?", "application")
	result := queryUpdate.Updates(body)

	if result.RowsAffected == 0 {
		body["namespace"] = "application"
		query.Create(body)
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

func Settings(router fiber.Router) {
	router.Get("/settings", middleware.JWTCheck([]string{"admin"}, nil), getSettings)
	router.Patch("/settings", middleware.JWTCheck([]string{"admin"}, nil), patchSettings)
}
