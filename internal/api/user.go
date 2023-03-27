package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/worldline-go/chore/internal/server/middleware"
	"github.com/worldline-go/chore/models"
	"github.com/worldline-go/chore/models/apimodels"
	"github.com/worldline-go/chore/pkg/registry"
	"github.com/worldline-go/chore/pkg/sec"
)

type UserDataID struct {
	models.UserData
	apimodels.ID
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
// @Param limit query int false "set the limit, default is 20"
// @Param offset query int false "set the offset, default is 0"
// @Param search query string string "search item"
// @Success 200 {object} apimodels.DataMeta{data=[]UserDataID{},meta=apimodels.Meta{}}
// @failure 400 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func listUsers(c *fiber.Ctx) error {
	users := []UserDataID{}

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

	query := reg.DB.WithContext(c.UserContext()).Model(&models.User{}).Limit(meta.Limit).Offset(meta.Offset)

	if meta.Search != "" {
		query = query.Where("name LIKE ?", meta.Search+"%")
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

	// get counts
	query = reg.DB.WithContext(c.UserContext()).Model(&models.User{})
	if meta.Search != "" {
		query = query.Where("name LIKE ?", meta.Search+"%")
	}

	query.Count(&meta.Count)

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
// @Success 200 {object} apimodels.Data{data=UserDataID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func getUser(c *fiber.Ctx) error {
	id := c.Query("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: apimodels.ErrRequiredIDName.Error(),
			},
		)
	}

	user := new(UserDataID)

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.User{})
	if id != "" {
		query = query.Where("id = ?", id)
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
// @Success 204 "No Content"
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func deleteUser(c *fiber.Ctx) error {
	id := c.Query("id")

	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: apimodels.ErrRequiredIDName.Error(),
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext())
	if id != "" {
		query = query.Where("id = ?", id)
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

	if body.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "password is required",
			},
		)
	}

	// hash password
	if hashedPassword, err := sec.HashPassword([]byte(body.Password)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	} else { //nolint:golint // required for value scope
		body.Password = string(hashedPassword)
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
			ModelCU: apimodels.ModelCU{
				ID: apimodels.ID{ID: id},
			},
		},
	)

	// check write error
	if result.Error != nil && errors.Is(result.Error, gorm.ErrDuplicatedKey) {
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

	if v, ok := body["id"].(string); !ok || v == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: "id is required and cannot be empty",
			},
		)
	}

	// hash password
	if v, ok := body["password"].(string); ok {
		if hashedPassword, err := sec.HashPassword([]byte(v)); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				apimodels.Error{
					Error: err.Error(),
				},
			)
		} else { //nolint:golint // required for value scope
			body["password"] = hashedPassword
		}
	}

	if body["groups"] != nil {
		var err error

		body["groups"], err = json.Marshal(body["groups"])
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				apimodels.Error{
					Error: err.Error(),
				},
			)
		}
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.User{}).Where("id = ?", body["id"])

	result := query.Updates(body)

	// check write error
	if result.Error != nil && errors.Is(result.Error, gorm.ErrDuplicatedKey) {
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
	resultData["id"] = body["id"]

	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: resultData,
		},
	)
}

func User(router fiber.Router) {
	router.Get("/users", middleware.JWTCheck([]string{"admin"}, nil), listUsers)
	router.Get("/user", middleware.JWTCheck([]string{"admin"}, middleware.IDFromQuery), getUser)
	router.Post("/user", middleware.JWTCheck([]string{"admin"}, nil), postUser)
	router.Patch("/user", middleware.JWTCheck([]string{"admin"}, middleware.IDFromBody), patchUser)
	router.Delete("/user", middleware.JWTCheck([]string{"admin"}, middleware.IDFromQuery), deleteUser)
}
