package apitest

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/rs/zerolog/log"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
)

// @Summary New Test
// @Tags test
// @Description Send and record new user
// @Router /test [post]
// @Param payload body models.TestPure{} false "send test object"
// @Success 200 {object} apimodels.Data{data=apimodels.ID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postTest(c *fiber.Ctx) error {
	body := new(models.TestPure)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	if body.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "name is required",
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	tx := reg.DB.Begin()

	test := models.Test{
		TestPure: models.TestPure{
			Name: body.Name,
		},
	}

	result := tx.WithContext(c.UserContext()).Create(&test)

	if result.Error != nil {
		tx.Rollback()
	}

	log.Debug().Msgf("TEST-ID: %v", test.ID)

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

	tx.Commit()

	// return recorded data's id
	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: apimodels.ID{ID: uuid.Nil},
		},
	)
}

func Test(router fiber.Router) {
	router.Post("/test", postTest)
}
