package handler

import (
	"fmt"
	"path"

	"github.com/gofiber/fiber/v2"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/store/inf"
)

func RouterKV(f fiber.Router, v string) {
	newRouter := f.Group(v)
	newRouter.Use(func(c *fiber.Ctx) error {
		c.Locals("type", v[1:])

		return c.Next() //nolint:wrapcheck // not need here
	})

	// list and get sameplace but query different
	newRouter.Get("/", apiList)
	newRouter.Get("/", apiGet)

	newRouter.Put("/", apiPut)
	newRouter.Post("/", apiPost)
	newRouter.Delete("/", apiDelete)
}

// @Summary Get List
// @Description Get list of keys
// @Param type path string true "type" Enums(templates, auths, binds)
// @Param key  query string false "key of the file"
// @Param list query string true "is it for listing" default(true)
// @Router /kv/{type}/ [get]
// @Success 200 {array} string "folder/ or file"
// @Success 418 {string} string "not a list"
// @failure 500 {string} string
func apiList(c *fiber.Ctx) error {
	if c.Query("list", "false") == "false" {
		// this is not a list go to next handler
		return c.Next()
	}

	search := c.Locals("type").(string)
	if key := c.Query("key"); key != "" {
		search = path.Join(search, key)
	}

	crud := c.Locals("storeHandler").(*inf.CRUD)

	list, errCrud := (*crud).List(search)
	if errCrud != nil {
		return fiber.NewError(errCrud.GetCode(), errCrud.Error())
	}

	err := c.JSON(list)
	if err != nil {
		return fmt.Errorf("apiList to json; %w", err)
	}

	return nil
}

// @Summary Get Specific key
// @Description Get the value of the key
// @Param type path string true "type" Enums(templates, auths, binds)
// @Param key query string false "key of the file"
// @Router /kv/{type} [get]
// @Success 200 {string} string
// @failure 404 {string} string
// @failure 500 {string} string
func apiGet(c *fiber.Ctx) error {
	crud := c.Locals("storeHandler").(*inf.CRUD)

	search := c.Locals("type").(string)
	if key := c.Query("key"); key != "" {
		search = path.Join(search, key)
	}

	data, errCrud := (*crud).Get(search)
	if errCrud != nil {
		return fiber.NewError(errCrud.GetCode(), errCrud.Error())
	}

	err := c.JSON(data)
	if err != nil {
		return fmt.Errorf("apiGet to json; %w", err)
	}

	return nil
}

// @Summary New key or replace key
// @Description Set a key with value
// @Param type path string true "type" Enums(templates, auths, binds)
// @Param key query string false "key of the file"
// @Accept text/plain
// @Param payload body string false "any value to store"
// @Router /kv/{type} [put]
// @Success 200 {object} string
// @failure 500 {string} string
func apiPut(c *fiber.Ctx) error {
	crud := c.Locals("storeHandler").(*inf.CRUD)

	search := c.Locals("type").(string)
	if key := c.Query("key"); key != "" {
		search = path.Join(search, key)
	}

	errCrud := (*crud).Put(search, c.Body())
	if errCrud != nil {
		return fiber.NewError(errCrud.GetCode(), errCrud.Error())
	}

	return c.SendStatus(fiber.StatusOK) //nolint:wrapcheck // not need
}

// @Summary New key
// @Description Set a key with value
// @Param type path string true "type" Enums(templates, auths, binds)
// @Param key query string false "key of the file"
// @Accept text/plain
// @Param payload body string false "any value to store"
// @Router /kv/{type} [post]
// @Success 200 {string} string
// @Success 406 {string} string
// @Success 409 {string} string
// @failure 500 {string} string
func apiPost(c *fiber.Ctx) error {
	crud := c.Locals("storeHandler").(*inf.CRUD)

	search := c.Locals("type").(string)
	if key := c.Query("key"); key != "" {
		search = path.Join(search, key)
	}

	errCrud := (*crud).Post(search, c.Body())
	if errCrud != nil {
		return fiber.NewError(errCrud.GetCode(), errCrud.Error())
	}

	return c.SendStatus(fiber.StatusOK) //nolint:wrapcheck // not need
}

// @Summary Delete key
// @Description Delete key
// @Param type path string true "type" Enums(templates, auths, binds)
// @Param key query string false "key of the file"
// @Router /kv/{type} [delete]
// @Success 204 {object} map[string]interface{}
// @failure 404 {string} string
// @failure 500 {string} string
func apiDelete(c *fiber.Ctx) error {
	crud := c.Locals("storeHandler").(*inf.CRUD)

	search := c.Locals("type").(string)
	if key := c.Query("key"); key != "" {
		search = path.Join(search, key)
	}

	errCrud := (*crud).Delete(search)
	if errCrud != nil {
		return fiber.NewError(errCrud.GetCode(), errCrud.Error())
	}

	return c.SendStatus(fiber.StatusNoContent) //nolint:wrapcheck // not need
}
