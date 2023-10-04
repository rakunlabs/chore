package server

import (
	"github.com/labstack/echo/v4"

	"github.com/worldline-go/chore/docs"
	"github.com/worldline-go/chore/internal/config"
	echoSwagger "github.com/worldline-go/echo-swagger"
)

func routerSwagger(apiGroup *echo.Group) error {
	if err := docs.SetInfo(
		config.AppName,
		config.AppVersion,
	); err != nil {
		return err
	}

	apiGroup.GET("/swagger/*", echoSwagger.WrapHandler)

	return nil
}
