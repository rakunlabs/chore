package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/registry"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/server/middleware"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
)

type BindPureID struct {
	models.BindPure
	apimodels.ID
}

// @Summary List binds
// @Tags bind
// @Description Get list of the binds
// @Security ApiKeyAuth
// @Router /binds [get]
// @Param limit query int false "set the limit, default is 20"
// @Param offset query int false "set the offset, default is 0"
// @Success 200 {object} apimodels.DataMeta{data=[]BindPureID{},meta=apimodels.Meta}
// @failure 400 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func listBinds(c *fiber.Ctx) error {
	binds := []BindPureID{}

	meta := &UserMetaAdmin{
		Meta:  apimodels.Meta{Limit: apimodels.Limit},
		Admin: false,
	}

	if err := c.QueryParser(meta); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.Bind{}).Limit(meta.Limit).Offset(meta.Offset)
	result := query.Find(&binds)

	// check write error
	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	return c.Status(http.StatusOK).JSON(
		apimodels.DataMeta{
			Meta: meta.Meta,
			Data: apimodels.Data{Data: binds},
		},
	)
}

// @Summary Get bind
// @Tags bind
// @Description Get one bind with id or name
// @Security ApiKeyAuth
// @Router /bind [get]
// @Param id query string false "get by id"
// @Param name query string false "get by name"
// @Success 200 {object} apimodels.Data{data=BindPureID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func getBind(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Query("name")

	if id == "" && name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: apimodels.ErrRequiredIDName.Error(),
			},
		)
	}

	getData := new(BindPureID)

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.Auth{})
	if id != "" {
		query = query.Where("id = ?", id)
	}

	if name != "" {
		query = query.Where("name = ?", name)
	}

	result := query.First(&getData)

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

	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: getData,
		},
	)
}

// @Summary New bind
// @Tags bind
// @Description Send and record new bind
// @Security ApiKeyAuth
// @Router /bind [post]
// @Param payload body models.BindPure{} false "send bind object"
// @Success 200 {object} apimodels.Data{data=apimodels.ID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postBind(c *fiber.Ctx) error {
	body := new(models.BindPure)
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

	id, err := uuid.NewUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	result := reg.DB.WithContext(c.UserContext()).Create(
		&models.Bind{
			BindPure: *body,
			ModelS: apimodels.ModelS{
				ID: apimodels.ID{ID: id},
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
			Data: apimodels.ID{ID: id},
		},
	)
}

// @Summary Bind auth
// @Tags bind
// @Description Bind with a few data, id must exist in request
// @Security ApiKeyAuth
// @Router /bind [patch]
// @Param payload body BindPureID{} false "send part of the user object"
// @Success 200 {object} apimodels.Data{data=apimodels.ID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func patchBind(c *fiber.Ctx) error {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	isSetBodyID := false
	if v, ok := body["id"].(string); ok && v != "" {
		isSetBodyID = true
	}

	isSetBodyName := false
	if v, ok := body["name"].(string); ok && v != "" {
		isSetBodyName = true
	}

	if !isSetBodyName && !isSetBodyID {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "id or name is required and cannot be empty",
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.Bind{})

	if isSetBodyID {
		query = query.Where("id = ?", body["id"])
	}

	if isSetBodyName {
		query = query.Where("name = ?", body["name"])
	}

	result := query.Updates(body)

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

	resultData := make(map[string]interface{})
	if isSetBodyID {
		resultData["id"] = body["id"]
	}

	if isSetBodyName {
		resultData["name"] = body["name"]
	}

	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: resultData,
		},
	)
}

// @Summary Delete bind
// @Tags bind
// @Description Delete with id or name
// @Security ApiKeyAuth
// @Router /bind [delete]
// @Param id query string false "get by id"
// @Param name query string false "get by name"
// @Success 204 "No Content"
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func deleteBind(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Query("name")

	if id == "" && name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: apimodels.ErrRequiredIDName,
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext())
	if id != "" {
		query = query.Where("id = ?", id)
	}

	if name != "" {
		query = query.Where("name = ?", name)
	}

	// delete directly in DB
	result := query.Unscoped().Delete(&models.Auth{})

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

func Bind(router fiber.Router) {
	router.Get("/binds", middleware.JWTCheck(""), listBinds)
	router.Get("/bind", middleware.JWTCheck(""), getBind)
	router.Post("/bind", middleware.JWTCheck(""), postBind)
	router.Patch("/bind", middleware.JWTCheck(""), patchBind)
	router.Delete("/bind", middleware.JWTCheck(""), deleteBind)
}
