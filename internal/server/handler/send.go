package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/store/inf"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/request"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/translate"
)

// Send request
// @Summary Send request
// @Description Send request to api.
// @Accept */*
// @Param payload body string false "values"
// @Param key query string false "key of the binds entry"
// @Router /send [post]
// @Success 201 {object} map[string]interface{}
func Send(c *fiber.Ctx) error {
	crud := c.Locals("storeHandler").(inf.CRUD)

	search := "binds"
	if key := c.Query("key"); key != "" {
		search = path.Join(search, key)
	}

	var values map[string]interface{}

	if err := yaml.Unmarshal(c.Body(), &values); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	datas, errCrud := crud.Get(search)
	if errCrud != nil {
		return fiber.NewError(errCrud.GetCode(), errCrud.Error())
	}

	templateEngine := c.Locals("templateEngine").(*translate.Template)
	client := c.Locals("client").(*request.Client)

	for dataIndex := range datas {
		var bind map[string]interface{}
		if err := json.Unmarshal(datas[dataIndex], &bind); err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}

		templateName := bind["template"].(string)

		templates, errCrud := crud.Get(path.Join("templates", templateName))
		if errCrud != nil {
			return fiber.NewError(errCrud.GetCode(), errCrud.Error())
		}

		authenticationName := bind["authentication"].(string)

		authenticationsRaw, errCrud := crud.Get(path.Join("auths", authenticationName))
		if errCrud != nil {
			return fiber.NewError(errCrud.GetCode(), errCrud.Error())
		}

		authentications := make([]map[string]interface{}, len(authenticationsRaw))
		for authIndex := range authentications {
			if err := json.Unmarshal(authenticationsRaw[authIndex], &authentications[authIndex]); err != nil {
				return fiber.NewError(http.StatusInternalServerError, err.Error())
			}
		}

		for templateIndex := range templates {
			payload, err := templateEngine.Ext(values, string(templates[templateIndex]))
			if err != nil {
				return fiber.NewError(http.StatusInternalServerError, err.Error())
			}

			for authIndex := range authentications {
				var headers map[string]string

				headersString := authentications[authIndex]["headers"].(string)

				if err := yaml.Unmarshal([]byte(headersString), &headers); err != nil {
					return fiber.NewError(http.StatusInternalServerError, err.Error())
				}

				if data, err := client.Send(
					c.Context(),
					authentications[authIndex]["URL"].(string),
					authentications[authIndex]["method"].(string),
					headers,
					payload,
				); err != nil {
					return fiber.NewError(http.StatusInternalServerError, err.Error())
				} else {
					_ = c.SendStream(bytes.NewReader(data))
				}
			}
		}
	}

	return c.SendStatus(fiber.StatusOK) //nolint:wrapcheck // not need
}
