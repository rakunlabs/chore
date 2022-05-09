package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgconn"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/parser"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/sec"
)

type LoginModel struct {
	Login    string `json:"login" example:"admin"`
	Password string `json:"password" example:"admin"`
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

	if body.Login == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "required username, email address or user id",
			},
		)
	}

	return loginAndGetToken(c, *body, false)
}

// @Summary Login
// @Description Get JWT token for 1 hour
// @Tags public
// @Router /login [get]
// @Param raw query bool false "raw token output"
// @Security BasicAuth
// @Success 200 {object} apimodels.Data{data=LoginToken{}}
// @failure 400 {object} apimodels.Error{}
// @failure 401 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func getLogin(c *fiber.Ctx) error {
	// parse authorization basic if it is exist
	login, err := parser.GetAuthorizationBasic(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	raw, err := parser.GetQueryBool(c, "raw")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	body := LoginModel{
		Login:    login.User,
		Password: login.Pass,
	}

	return loginAndGetToken(c, body, raw)
}

func loginAndGetToken(c *fiber.Ctx, body LoginModel, raw bool) error {
	// declare user for result
	user := UserPureID{}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.User{})
	query = query.Where("name = ?", body.Login).Or("email = ?", body.Login).Or("id = ?", body.Login)

	result := query.First(&user)

	if !sec.CheckHashPassword([]byte(user.Password), []byte(body.Password)) {
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
			"user":   user.ID.ID,
			"groups": user.Groups.Groups,
			"type":   models.TypeAccessToken,
		}, reg.JWT.DefExpFunc())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	if raw {
		return c.Status(http.StatusOK).Send([]byte(token))
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
	router.Get("/login", getLogin)
}
