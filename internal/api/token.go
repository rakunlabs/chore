package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/server/middleware"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
)

type TokenDataByID struct {
	models.TokenDataBy
	apimodels.ID
}

// @Summary List tokens
// @Tags token
// @Description Get list of the tokens
// @Security ApiKeyAuth
// @Router /tokens [get]
// @Param limit query int false "set the limit, default is 20"
// @Param offset query int false "set the offset, default is 0"
// @Success 200 {object} apimodels.DataMeta{data=[]TokenDataByID{},meta=apimodels.Meta{}}
// @failure 400 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func listTokens(c *fiber.Ctx) error {
	meta := apimodels.Meta{Limit: apimodels.Limit}

	if err := c.QueryParser(&meta); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	tokens := []TokenDataByID{}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.Token{}).Limit(meta.Limit).Offset(meta.Offset)

	if result := query.Find(&tokens); result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	// get counts
	reg.DB.WithContext(c.UserContext()).Model(&models.Token{}).Count(&meta.Count)

	return c.Status(http.StatusOK).JSON(
		apimodels.DataMeta{
			Meta: meta,
			Data: apimodels.Data{Data: tokens},
		},
	)
}

// @Summary Get token
// @Tags token
// @Description Get token with token id
// @Security ApiKeyAuth
// @Router /token [get]
// @Param id query string false "get by token id"
// @Success 200 {object} apimodels.Data{data=models.Token{}}
// @failure 400 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func getToken(c *fiber.Ctx) error {
	id := c.Query("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "id is required",
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	token := models.Token{}
	result := reg.DB.WithContext(c.UserContext()).Where("id = ?", id).Find(&token)

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

	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: token,
		},
	)
}

// @Summary New token
// @Tags token
// @Description Send and record PAT token
// @Security ApiKeyAuth
// @Router /token [post]
// @Param payload body models.TokenData{} false "token parameters"
// @Success 200 {object} apimodels.Data{data=models.Token{}}
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postToken(c *fiber.Ctx) error {
	body := new(models.TokenData)
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
				Error: "name is empty",
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	var groups []string

	if body.Groups.Groups != nil {
		if err := json.Unmarshal(body.Groups.Groups, &groups); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				apimodels.Error{
					Error: err.Error(),
				},
			)
		}
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	var unixLastDate int64
	if body.Date != nil && !body.Date.IsZero() {
		unixLastDate = body.Date.Unix()

		if unixLastDate-time.Now().Unix() <= 0 {
			return c.Status(http.StatusBadRequest).JSON(
				apimodels.Error{
					Error: "date should be older than now",
				},
			)
		}
	}

	createdBy, ok := c.Locals("id").(uuid.UUID)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: "locals id not uuid type",
			},
		)
	}

	// generate new PAT token
	token, err := reg.JWT.Generate(
		map[string]interface{}{
			"id":     id,
			"user":   createdBy,
			"type":   models.TypePersonalAccessToken,
			"groups": groups,
		}, unixLastDate)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	createToken := models.Token{
		ModelC: apimodels.ModelC{
			ID: apimodels.ID{ID: id},
		},
		TokenPure: models.TokenPure{
			TokenDataBy: models.TokenDataBy{
				TokenData: *body,
				CreatedBy: createdBy,
			},
			TokenPrivate: models.TokenPrivate{
				Token: token,
			},
		},
	}

	result := reg.DB.WithContext(c.UserContext()).Create(&createToken)

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
			Data: createToken,
		},
	)
}

// @Summary Delete token
// @Tags token
// @Description Delete with token
// @Security ApiKeyAuth
// @Router /token [delete]
// @Param id query string false "get by token id"
// @Success 204 "No Content"
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func deleteToken(c *fiber.Ctx) error {
	id := c.Query("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "id is required",
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	// delete directly in DB
	result := reg.DB.WithContext(c.UserContext()).Where("id = ?", id).Unscoped().Delete(&models.Token{})

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
// @Param payload body models.TokenPrivate{} false "send a token"
// @Success 200 {object} apimodels.Data{data=models.TokenPrivate{}}
// @failure 400 {object} apimodels.Error{}
// @failure 401 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postTokenCheck(c *fiber.Ctx) error {
	body := new(models.TokenPrivate)
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
			Data: models.TokenPrivate{Token: body.Token},
		},
	)
}

// @Summary Renew token
// @Tags token
// @Description Get new token based on old token
// @Router /token/renew [get]
// @Success 200 {object} apimodels.Data{data=models.TokenPrivate{}}
// @failure 400 {object} apimodels.Error{}
// @failure 401 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func getTokenRenew(c *fiber.Ctx) error {
	reg := registry.Reg().Get(c.Locals("registry").(string))

	// check token is valid
	token, err := reg.JWT.Renew(c.Locals("token").(string), reg.JWT.DefExpFunc())
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(
			apimodels.Error{
				Error: fmt.Sprintf("renew failed: %v", err.Error()),
			},
		)
	}

	// return recorded data's id
	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: models.TokenPrivate{Token: token},
		},
	)
}

func Token(router fiber.Router) {
	router.Post("/token/check", postTokenCheck)
	router.Get("/token/renew", middleware.JWTCheck(nil, nil), getTokenRenew)
	router.Get("/tokens", middleware.JWTCheck(nil, nil), listTokens)
	router.Get("/token", middleware.JWTCheck(nil, nil), getToken)
	router.Post("/token", middleware.JWTCheck(nil, nil), postToken)
	router.Delete("/token", middleware.JWTCheck(nil, nil), deleteToken)
}
