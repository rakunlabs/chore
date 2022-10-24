package run

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/worldline-go/chore/internal/server/middleware"
	"github.com/worldline-go/chore/models/apimodels"
	"github.com/worldline-go/chore/pkg/registry"
	"github.com/worldline-go/chore/pkg/script/js"
)

// @Summary Run JS script
// @Tags run
// @Description Run JS script with scripts and input values
// @Security ApiKeyAuth
// @Router /run/js [post]
// @Param payload body runModel true "Script and inputs"
// @Accept plain
// @Success 200 {object} string "result of the script"
// @failure 400 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postJS(c *fiber.Ctx) error {
	body := defaultRunModel()
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	if body.Settings.Timeout != "" {
		duration, err := time.ParseDuration(body.Settings.Timeout)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(
				apimodels.Error{
					Error: err.Error(),
				},
			)
		}

		body.Settings.TimeoutDuration = duration
	}

	log.Ctx(c.UserContext()).Debug().Msgf("script run with API")

	runtime := js.NewGoja()
	parsedInputs := js.ParseInputs(body.Inputs)

	if body.Settings.Async {
		ctx := c.UserContext()
		logX := log.Ctx(ctx)
		wg := registry.Reg().WG

		wg.Add(1)
		go func() {
			defer wg.Done()

			result, err := runtime.RunScript(ctx, string(body.Script), parsedInputs)
			if err != nil {
				logX.Error().Err(err).Msg("error while running script")

				return
			}

			logX.Info().Msgf("%s", result)
		}()

		return c.SendStatus(http.StatusAccepted)
	}

	ctx := log.Ctx(c.UserContext()).WithContext(c.Context())

	result, err := runtime.RunScript(ctx, string(body.Script), parsedInputs)
	if err != nil && !errors.Is(err, js.ErrThrow) {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	// return recorded data's id
	_, err = c.Status(http.StatusOK).Write(result)

	return err
}

// @Summary Render template
// @Tags run
// @Description Render template with input values
// @Security ApiKeyAuth
// @Router /run/template [post]
// @Param payload body string false "send key values" SchemaExample()
// @Accept plain
// @Success 200 {object} apimodels.Data{data=apimodels.ID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postTemplate(c *fiber.Ctx) error {
	return c.Status(http.StatusNotImplemented).JSON(
		apimodels.Error{
			Error: "waiting implementation",
		},
	)
}

func API(router fiber.Router) {
	router.Post("/run/js", middleware.JWTCheck(nil, nil), postJS)
	router.Post("/run/template", middleware.JWTCheck(nil, nil), postTemplate)
}
