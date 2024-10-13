package handler

import (
	"log"

	"github.com/labstack/echo/v4"

	"Adoutchquizz/database"
	"Adoutchquizz/views"
)

func Index(c echo.Context) error {
    lastVideos, err := database.GetLatestVideos(10)
    if err != nil {
        log.Fatal(err)
        return err
    }
    return render(c, views.Index(lastVideos))
}
