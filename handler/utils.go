package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

    "Adoutchquizz/views/errors"
)

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func Error404(c echo.Context) error {
    c.Response().Status = 404
    return errors.Error404().Render(c.Request().Context(), c.Response())
}
