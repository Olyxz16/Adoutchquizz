package handler

import (
	"github.com/labstack/echo/v4"

    "Adoutchquizz/views"
)

func Index(c echo.Context) error {
    return render(c, views.Index())
}
