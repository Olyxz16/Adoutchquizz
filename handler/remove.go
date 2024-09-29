package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"Adoutchquizz/database"
)

func RemoveClip(c echo.Context) error {
    idstr := c.Param("id")
    id, err := strconv.Atoi(idstr)
    if err != nil {
        return err
    }
    err = database.DeleteClip(id)
    return err
}
