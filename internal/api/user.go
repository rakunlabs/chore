package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/worldline-go/chore/internal/server/middlewares"
	"github.com/worldline-go/chore/internal/utils"
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
// @Param search query string false "search item"
// @Success 200 {object} apimodels.DataMeta{data=[]UserDataID{},meta=apimodels.Meta{}}
// @failure 400 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func listUsers(c echo.Context) error {
	users := []UserDataID{}

	meta := &UserMetaAdmin{
		Meta:  apimodels.Meta{Limit: apimodels.Limit},
		Admin: false,
	}

	if err := c.Bind(meta); err != nil {
		return c.JSON(http.StatusBadRequest, apimodels.Error{Error: err.Error()})
	}

	query := registry.Reg.DB.WithContext(c.Request().Context()).Model(&models.User{}).Limit(meta.Limit).Offset(meta.Offset)

	if meta.Search != "" {
		query = query.Where("name LIKE ?", meta.Search+"%")
	}

	result := query.Find(&users)

	// check write error
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, apimodels.Error{Error: result.Error.Error()})
	}

	// get counts
	query = registry.Reg.DB.WithContext(c.Request().Context()).Model(&models.User{})
	if meta.Search != "" {
		query = query.Where("name LIKE ?", meta.Search+"%")
	}

	query.Count(&meta.Count)

	return c.JSON(http.StatusOK,
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
func getUser(c echo.Context) error {
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, apimodels.Error{Error: apimodels.ErrRequiredIDName.Error()})
	}

	user := new(UserDataID)

	query := registry.Reg.DB.WithContext(c.Request().Context()).Model(&models.User{})
	if id != "" {
		query = query.Where("id = ?", id)
	}

	result := query.First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, apimodels.Error{Error: result.Error.Error()})
	}

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, apimodels.Error{Error: result.Error.Error()})
	}

	return c.JSON(http.StatusOK,
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
func deleteUser(c echo.Context) error {
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, apimodels.Error{Error: apimodels.ErrRequiredIDName.Error()})
	}

	ctx := utils.Context(c)
	query := registry.Reg.DB.WithContext(ctx)
	if id != "" {
		query = query.Where("id = ?", id)
	}

	// delete directly in DB
	result := query.Unscoped().Delete(&models.User{})

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, apimodels.Error{Error: "not found any releated data"})
	}

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, apimodels.Error{Error: result.Error.Error()})
	}

	//nolint:wrapcheck // checking before
	return c.NoContent(http.StatusNoContent)
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
func postUser(c echo.Context) error {
	body := new(models.UserPure)
	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, apimodels.Error{Error: err.Error()})
	}

	if body.Name == "" {
		return c.JSON(http.StatusBadRequest, apimodels.Error{Error: "name is required"})
	}

	if body.Password == "" {
		return c.JSON(http.StatusBadRequest, apimodels.Error{Error: "password is required"})
	}

	// hash password
	hashedPassword, err := sec.HashPassword([]byte(body.Password))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apimodels.Error{Error: err.Error()})
	}

	body.Password = string(hashedPassword)

	id, err := uuid.NewUUID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apimodels.Error{Error: err.Error()})
	}

	ctx := utils.Context(c)

	result := registry.Reg.DB.WithContext(ctx).Create(
		&models.User{
			UserPure: *body,
			ModelCU: apimodels.ModelCU{
				ID: apimodels.ID{ID: id},
			},
		},
	)

	// check write error
	if result.Error != nil && errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return c.JSON(http.StatusConflict, apimodels.Error{Error: result.Error.Error()})
	}

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, apimodels.Error{Error: result.Error.Error()})
	}

	// return recorded data's id
	return c.JSON(http.StatusOK, apimodels.Data{Data: apimodels.ID{ID: id}})
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
func patchUser(c echo.Context) error {
	var body map[string]interface{}
	if err := c.Bind(&body); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	if v, ok := body["id"].(string); !ok || v == "" {
		return c.JSON(
			http.StatusBadRequest,
			apimodels.Error{
				Error: "id is required and cannot be empty",
			},
		)
	}

	// hash password
	if v, ok := body["password"].(string); ok {
		hashedPassword, err := sec.HashPassword([]byte(v))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, apimodels.Error{Error: err.Error()})
		}

		body["password"] = hashedPassword
	}

	if body["groups"] != nil {
		var err error

		body["groups"], err = json.Marshal(body["groups"])
		if err != nil {
			return c.JSON(http.StatusInternalServerError, apimodels.Error{Error: err.Error()})
		}
	}

	ctx := utils.Context(c)
	query := registry.Reg.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", body["id"])

	result := query.Updates(body)

	// check write error
	if result.Error != nil && errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return c.JSON(http.StatusConflict, apimodels.Error{Error: result.Error.Error()})
	}

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, apimodels.Error{Error: result.Error.Error()})
	}

	resultData := make(map[string]interface{})
	resultData["id"] = body["id"]

	return c.JSON(http.StatusOK,
		apimodels.Data{
			Data: resultData,
		},
	)
}

func User(e *echo.Group, authMiddleware echo.MiddlewareFunc) {
	e.GET("/users", listUsers, authMiddleware, middlewares.AdminRole)
	e.GET("/user", getUser, authMiddleware, middlewares.JWTCheck(middlewares.IDFromQuery), middlewares.AdminRole)
	e.POST("/user", postUser, authMiddleware, middlewares.AdminRole)
	e.PATCH("/user", patchUser, authMiddleware, middlewares.JWTCheck(middlewares.IDFromBody), middlewares.AdminRole)
	e.DELETE("/user", deleteUser, authMiddleware, middlewares.JWTCheck(middlewares.IDFromQuery), middlewares.AdminRole)
}
