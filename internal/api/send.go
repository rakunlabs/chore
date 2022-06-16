package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
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
// @failure 405 {object} apimodels.Error{}
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

	log.Ctx(c.UserContext()).Info().Msgf("call control=[%s] endpoint=[%s]", control.Name, endpoint)

	nodesReg, err := flow.StartFlow(c.UserContext(), control.Name, endpoint, c.Method(), content, reg, c.Body())
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

// EndpointCheck middleware is checking endpoint.
func EndpointCheck(c *fiber.Ctx) error {
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

	v := models.Endpoints{}

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

	endpoints := make(map[string]models.ControlEndpoint)
	if err := json.Unmarshal(v.Endpoints, &endpoints); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	endpointSpec, ok := endpoints[endpoint]
	if !ok {
		return c.Status(http.StatusNotFound).JSON(
			apimodels.Error{
				Error: fmt.Sprintf("endpoint %s not found", endpoint),
			},
		)
	}

	// method check
	allowMethod := false

	for _, endpointMethod := range endpointSpec.Methods {
		if endpointMethod == c.Method() {
			allowMethod = true

			break
		}
	}

	if !allowMethod {
		c.Response().Header.Set("Allow", strings.ToUpper(strings.Join(endpointSpec.Methods, ", ")))

		return c.Status(http.StatusMethodNotAllowed).JSON(
			apimodels.Error{
				Error: fmt.Sprintf("method %s not allowed", c.Method()),
			},
		)
	}

	// public check
	if !endpointSpec.Public {
		//nolint:wrapcheck // next middleware
		return c.Next()
	}

	c.Locals("skip-middleware-jwt", true)

	//nolint:wrapcheck // next middleware
	return c.Next()
}

func Send(router fiber.Router) {
	router.All("/send", EndpointCheck, middleware.JWTCheck(nil, nil), send)
}
