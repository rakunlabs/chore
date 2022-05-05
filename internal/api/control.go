package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/api/fn"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/server/middleware"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
)

type ControlPureContentID struct {
	models.ControlPureContent
	apimodels.ID
}

type ControlPureID struct {
	models.ControlPure
	apimodels.ID
}

// @Summary List controls
// @Tags control
// @Description Get list of the controls
// @Security ApiKeyAuth
// @Router /controls [get]
// @Param limit query int false "set the limit, default is 20"
// @Param offset query int false "set the offset, default is 0"
// @Success 200 {object} apimodels.DataMeta{data=[]ControlPureID{},meta=apimodels.Meta{}}
// @failure 400 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func listControls(c *fiber.Ctx) error {
	controlsPureID := []ControlPureID{}

	meta := &apimodels.Meta{Limit: apimodels.Limit}

	if err := c.QueryParser(meta); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.Control{}).Limit(meta.Limit).Offset(meta.Offset)

	result := query.Find(&controlsPureID)

	// check write error
	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: result.Error.Error(),
			},
		)
	}

	// get counts
	reg.DB.WithContext(c.UserContext()).Model(&models.Control{}).Count(&meta.Count)

	return c.Status(http.StatusOK).JSON(
		apimodels.DataMeta{
			Meta: meta,
			Data: apimodels.Data{Data: controlsPureID},
		},
	)
}

// @Summary Get control
// @Tags control
// @Description Get one control with id
// @Security ApiKeyAuth
// @Router /control [get]
// @Param id query string false "get by id"
// @Param name query string false "get by name"
// @Param nodata query bool false "not get content"
// @Param dump query bool false "get for record values"
// @Param pretty query bool false "pretty output for dump"
// @Success 200 {object} apimodels.Data{data=ControlPureContentID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func getControl(c *fiber.Ctx) error {
	nodata, err := fn.GetQueryBool(c, "nodata")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	id := c.Query("id")
	name := c.Query("name")

	if id == "" && name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: apimodels.ErrRequiredIDName.Error(),
			},
		)
	}

	dump, err := fn.GetQueryBool(c, "dump")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	pretty, err := fn.GetQueryBool(c, "pretty")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	controlContent := new(ControlPureContentID)
	control := new(ControlPureID)

	reg := registry.Reg().Get(c.Locals("registry").(string))

	query := reg.DB.WithContext(c.UserContext()).Model(&models.Control{})

	if id != "" {
		query = query.Where("id = ?", id)
	}

	if name != "" {
		query = query.Where("name = ?", name)
	}

	var result *gorm.DB
	if nodata {
		result = query.First(&control)
	} else {
		result = query.First(&controlContent)
	}

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

	var ret interface{}
	if nodata {
		ret = control
	} else {
		if dump {
			contentRaw, err := base64.StdEncoding.DecodeString(controlContent.Content)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(
					apimodels.Error{
						Error: err.Error(),
					},
				)
			}
			controlContent.Content = string(contentRaw)
		}
		ret = controlContent
	}

	if dump {
		return fn.JSON(c.Status(http.StatusOK), ret, pretty)
	}

	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: ret,
		},
	)
}

// @Summary New control
// @Tags control
// @Description Send and record new control
// @Security ApiKeyAuth
// @Router /control [post]
// @Param payload body models.ControlPureContent{} false "send control object"
// @Success 200 {object} apimodels.Data{data=apimodels.ID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func postControl(c *fiber.Ctx) error {
	var body models.ControlPureContent
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	if body.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: apimodels.ErrRequiredName.Error(),
			},
		)
	}

	body.Content = base64.StdEncoding.EncodeToString([]byte(body.Content))

	reg := registry.Reg().Get(c.Locals("registry").(string))

	id, err := uuid.NewUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	result := reg.DB.WithContext(c.UserContext()).Model(&models.Control{}).Create(
		&models.Control{
			ControlPureContent: body,
			ModelCU: apimodels.ModelCU{
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

// @Summary New or Update control
// @Tags control
// @Description Send and record control
// @Security ApiKeyAuth
// @Router /control [put]
// @Param payload body models.ControlPureContent{} false "send control object"
// @Success 204 "No Content"
// @failure 400 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func putControl(c *fiber.Ctx) error {
	var body models.ControlPureContent
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	if body.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			apimodels.Error{
				Error: apimodels.ErrRequiredName.Error(),
			},
		)
	}

	body.Content = base64.StdEncoding.EncodeToString([]byte(body.Content))

	reg := registry.Reg().Get(c.Locals("registry").(string))

	id, err := uuid.NewUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			apimodels.Error{
				Error: err.Error(),
			},
		)
	}

	result := reg.DB.WithContext(c.UserContext()).Model(&models.Control{}).Clauses(
		clause.OnConflict{
			UpdateAll: true,
			Columns:   []clause.Column{{Name: "name"}},
		}).Create(
		&models.Control{
			ControlPureContent: body,
			ModelCU: apimodels.ModelCU{
				ID: apimodels.ID{ID: id},
			},
		},
	)

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

// @Summary Replace control
// @Tags control
// @Description Replace with new data, id or name must exist in request
// @Security ApiKeyAuth
// @Router /control [patch]
// @Param payload body ControlPureID{} false "send part of the control object"
// @Success 200 {object} apimodels.Data{data=apimodels.ID{}}
// @failure 400 {object} apimodels.Error{}
// @failure 409 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func patchControl(c *fiber.Ctx) error {
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

	content, _ := body["content"].(string)
	body["content"] = base64.StdEncoding.EncodeToString([]byte(content))

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

	query := reg.DB.WithContext(c.UserContext()).Model(&models.Control{}).Where("id = ?", body["id"])

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
	resultData["id"] = body["id"]

	return c.Status(http.StatusOK).JSON(
		apimodels.Data{
			Data: resultData,
		},
	)
}

// @Summary Delete control
// @Tags control
// @Description Delete with id or name
// @Security ApiKeyAuth
// @Router /control [delete]
// @Param id query string false "get by id"
// @Param name query string false "get by name"
// @Success 204 "No Content"
// @failure 400 {object} apimodels.Error{}
// @failure 404 {object} apimodels.Error{}
// @failure 500 {object} apimodels.Error{}
func deleteControl(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Query("name")

	if id == "" && name == "" {
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

	if name != "" {
		query = query.Where("name = ?", name)
	}

	// delete directly in DB
	result := query.Unscoped().Delete(&models.Control{})

	if result.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(
			apimodels.Error{
				Error: apimodels.ErrNotFound.Error(),
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

func Control(router fiber.Router) {
	router.Get("/controls", middleware.JWTCheck(nil, nil), listControls)
	router.Get("/control", middleware.JWTCheck(nil, nil), getControl)
	router.Post("/control", middleware.JWTCheck(nil, nil), postControl)
	router.Put("/control", middleware.JWTCheck(nil, nil), putControl)
	router.Patch("/control", middleware.JWTCheck(nil, nil), patchControl)
	router.Delete("/control", middleware.JWTCheck(nil, nil), deleteControl)
}
