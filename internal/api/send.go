package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/server/middleware"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
)

// @Summary Send request
// @Description Send request with bind id or name
// @Security ApiKeyAuth
// @Router /send [post]
// @Router /send [get]
// @Param endpoint query string true "set endpoint"
// @Param control query string true "set control"
// @Param payload body string false "send key values" SchemaExample()
// @Accept plain
// @Success 200 {object} interface{} "respond from related server"
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func send(c *fiber.Ctx) error {
	endpoint := c.Locals("endpoint").(string)
	name := c.Locals("control").(string)

	control := models.Control{}

	reg := registry.Reg().Get(c.Locals("registry").(string))
	query := reg.DB.WithContext(c.UserContext()).Where("name = ?", name)
	result := query.First(&control)

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

	// file, err := c.FormFile("document")
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(
	// 		apimodels.Error{
	// 			Error: err.Error(),
	// 		},
	// 	)
	// }

	content, err := base64.StdEncoding.DecodeString(control.Content)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	nodesReg, err := flow.StartFlow(c.UserContext(), control.Name, endpoint, content, reg, c.Body())
	if errors.Is(err, flow.ErrEndpointNotFound) {
		return c.Status(http.StatusNotFound).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	respondChan := nodesReg.GetChan()
	if respondChan == nil {
		return c.SendStatus(http.StatusAccepted)
	}

	valueChan := <-respondChan

	for k, v := range valueChan.Header {
		c.Response().Header.Set(k, fmt.Sprint(v))
	}

	return c.Status(valueChan.Status).Send(valueChan.Data)
}

// PublicCheck middleware is checking endpoint.
func PublicCheck(c *fiber.Ctx) error {
	endpoint := c.Query("endpoint")
	name := c.Query("control")

	if endpoint == "" || name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "endpoint and control parameters cannot be empty",
			},
		)
	}

	c.Locals("endpoint", endpoint)
	c.Locals("control", name)

	v := models.ControlPublicEndpoint{}

	reg := registry.Reg().Get(c.Locals("registry").(string))
	query := reg.DB.WithContext(c.UserContext()).Model(&models.Control{}).Where("name = ?", name)
	result := query.First(&v)

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

	if v.PublicEndpoints == nil {
		//nolint:wrapcheck // next middleware
		return c.Next()
	}

	publicEndpoints := make([]string, 0)
	if err := json.Unmarshal(v.PublicEndpoints, &publicEndpoints); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	public := false

	for _, e := range publicEndpoints {
		if e == endpoint {
			public = true

			break
		}
	}

	if public {
		// next middleware, order important
		_ = c.Next()
	}

	//nolint:wrapcheck // next middleware
	return c.Next()
}

func Send(router fiber.Router) {
	router.Post("/send", PublicCheck, middleware.JWTCheck(nil, nil), send)
	router.Get("/send", PublicCheck, middleware.JWTCheck(nil, nil), send)
}
