package server

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/rs/zerolog/log"
)

//go:embed dist/*
var embedWeb embed.FS

func setFileHandler(app *fiber.App) {
	embedWebFolder, err := fs.Sub(embedWeb, "dist")
	if err != nil {
		log.Logger.Error().Err(err).Msg("cannot go to sub folder [dist]")
	}

	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(embedWebFolder),
	}))
}
