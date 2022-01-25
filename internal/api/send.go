package api

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/registry"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/server/middleware"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
)

// @Summary Send request
// @Description Send request with bind id or name
// @Security ApiKeyAuth
// @Router /send [post]
// @Param id query string false "get by id"
// @Param name query string false "get by name"
// @Param payload body map[string]interface{} false "send key values"
// @Success 200 {object} interface{} "respond from related server"
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postSend(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Query("name")

	if id == "" && name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: apimodels.ErrRequiredIDName,
			},
		)
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext())
	if id != "" {
		query = query.Where("binds.id = ?", id)
	}

	if name != "" {
		query = query.Where("binds.name = ?", name)
	}

	bind := models.Bind{}

	result := query.Joins("Template").Joins("Auth").First(&bind)

	log.Debug().Msgf("%+v", bind)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(
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

	content, err := base64.StdEncoding.DecodeString(bind.Template.Content)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	payload, err := reg.Template.Ext(body, string(content))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	data, err := reg.Client.Send(
		c.UserContext(),
		bind.Auth.URL,
		// authentications[authIndex]["URL"].(string)+"/"+params+queryString,
		bind.Auth.Method,
		bind.Auth.Headers,
		payload,
	)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	return c.SendStream(bytes.NewReader(data))
}

func Send(router fiber.Router) {
	router.Post("/send", middleware.JWTCheck(""), postSend)
}
