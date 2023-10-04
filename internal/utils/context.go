package utils

import (
	"context"

	"github.com/labstack/echo/v4"
)

func Context(c echo.Context) context.Context {
	ctx, _ := c.Get("context").(context.Context)
	if ctx == nil {
		ctx = context.Background()
	}

	return ctx
}
