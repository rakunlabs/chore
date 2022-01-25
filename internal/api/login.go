package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgconn"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/registry"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/sec"
)

type LoginModel struct {
	Name     string `json:"name" example:"userX"`
	Password string `json:"password" example:"pass1234"`
}

type LoginToken struct {
	Token string `json:"token" example:"tokenJWT"`
}

// @Summary Login
// @Description Get JWT token for 1 hour
// @Tags public
// @Router /login [post]
// @Param payload body LoginModel{} false "send login object"
// @Success 200 {object} apimodels.Data{data=LoginToken{}}
// @failure 400 {object} apimodels.Error{}
// @failure 401 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postLogin(c *fiber.Ctx) error {
	body := new(LoginModel)
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

	// declare user for result
	user := UserPureID{}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	result := reg.DB.WithContext(c.UserContext()).Model(&models.User{}).Where(
		"name = ?", body.Name,
	).First(&user)

	if !sec.CheckHashPassword(user.Password, body.Password) {
		return c.Status(http.StatusUnauthorized).JSON(
			apimodels.Error{
				Error: "name or password wrong",
			},
		)
	}

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

	// generate JWT token
	token, err := reg.JWT.Generate(
		map[string]interface{}{
			"id":    user.ID.ID.String(),
			"admin": user.Admin,
			"type":  models.TypeAccessToken,
		}, reg.JWT.DefExpFunc())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	// return recorded data's id
	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: LoginToken{
				Token: token,
			},
		},
	)
}

func Login(router fiber.Router) {
	router.Post("/login", postLogin)
}
