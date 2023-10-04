package server

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

//go:embed dist/*
var embedWeb embed.FS

func setFileHandler(e *echo.Group) {
	embedWebFolder, err := fs.Sub(embedWeb, "dist")
	if err != nil {
		log.Error().Err(err).Msg("cannot go to sub folder [dist]")
	}

	handlerFunc := http.FileServer(http.FS(embedWebFolder)).ServeHTTP

	e.Use(func(_ echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-cache")
			handlerFunc(c.Response().Writer, c.Request())

			return nil
		}
	})
}
