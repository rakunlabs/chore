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
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/sec"
)

type UserData struct {
	ID    string `json:"id" example:"cf8a07d4-077e-402e-a46b-ac0ed50989ec"`
	Name  string `json:"name" example:"username"`
	Admin bool   `json:"admin" example:"true"`
}

type UserPureID struct {
	models.UserPure
	apimodels.ID
}

type UserMetaAdmin struct {
	Admin bool `json:"admin,omitempty" query:"admin"`
	apimodels.Meta
}

// @Summary List users
// @Tags user
// @Description Get list of the users
// @Security ApiKeyAuth
// @Router /users [get]
// @Param admin query bool false "set admin rights"
// @Param limit query int false "set the limit, default is 20"
// @Param offset query int false "set the offset, default is 0"
// @Success 200 {object} apimodels.DataMeta{data=[]UserData{},meta=apimodels.Meta}
// @failure 400 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func listUsers(c *fiber.Ctx) error {
	users := []UserData{}

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

	// Select("id", "name", "admin")
	query := reg.DB.WithContext(c.UserContext()).Model(&models.User{}).Limit(meta.Limit).Offset(meta.Offset)

	if c.Query("admin", "") != "" {
		query = query.Where(
			"admin = ?", meta.Admin,
		)
	}

	result := query.Find(&users)

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
			Data: apimodels.Data{Data: users},
		},
	)
}

// @Summary Get user
// @Tags user
// @Description Get one user with id or name
// @Security ApiKeyAuth
// @Router /user [get]
// @Param id query string false "get by id"
// @Param name query string false "get by name"
// @Success 200 {object} apimodels.Data{data=UserData{}}
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func getUser(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Query("name")

	if id == "" && name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: apimodels.ErrRequiredIDName.Error(),
			},
		)
	}

	user := new(UserData)

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.User{})
	if id != "" {
		query = query.Where("id = ?", id)
	}

	if name != "" {
		query = query.Where("name = ?", name)
	}

	result := query.First(&user)

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
			Data: user,
		},
	)
}

// @Summary Delete user
// @Tags user
// @Description Delete with id or name
// @Security ApiKeyAuth
// @Router /user [delete]
// @Param id query string false "get by id"
// @Param name query string false "get by name"
// @Success 204 "No Content"
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func deleteUser(c *fiber.Ctx) error {
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
	result := query.Unscoped().Delete(&models.User{})

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

// @Summary New user
// @Tags user
// @Description Send and record new user
// @Security ApiKeyAuth
// @Router /user [post]
// @Param payload body models.UserPure{} false "send user object"
// @Success 200 {object} apimodels.Data{data=apimodels.ID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postUser(c *fiber.Ctx) error {
	body := new(models.UserPure)
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

	// hash password
	if hashedPassword, err := sec.HashPassword(body.Password); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	} else { //nolint:golint // required for value scope
		body.Password = hashedPassword
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
		&models.User{
			UserPure: *body,
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

// @Summary Replace user
// @Tags user
// @Description Replace with new data, id or name must exist in request
// @Security ApiKeyAuth
// @Router /user [patch]
// @Param payload body UserPureID{} false "send part of the user object"
// @Success 200 {object} apimodels.Data{data=apimodels.ID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func patchUser(c *fiber.Ctx) error {
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

	// hash password
	if v, ok := body["password"].(string); ok {
		if hashedPassword, err := sec.HashPassword(v); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				apimodels.Error{
					Error: err.Error(),
				},
			)
		} else { //nolint:golint // required for value scope
			body["password"] = hashedPassword
		}
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.User{})

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

func User(router fiber.Router) {
	router.Get("/users", middleware.JWTCheck(""), listUsers)
	router.Get("/user", middleware.JWTCheck(""), getUser)
	router.Post("/user", middleware.JWTCheck(""), postUser)
	router.Patch("/user", middleware.JWTCheck(""), patchUser)
	router.Delete("/user", middleware.JWTCheck(""), deleteUser)
}
