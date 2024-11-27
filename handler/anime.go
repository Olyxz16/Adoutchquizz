package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"Adoutchquizz/database"
)


func Anime(c echo.Context) error {
    action := c.Param("action")
    switch action {
    case "add": return addAnime(c)
    }
    return Error404(c)
}

func addAnime(c echo.Context) error {
    form, err := c.FormParams()
    if err != nil {
        log.Print(err)
        return c.String(http.StatusOK, "Erreur formulaire")
    }
    title := form.Get("title")
    if title == "" {
        err = fmt.Errorf("Title is empty")
        log.Print(err)
        return c.String(http.StatusOK, "Titre vide")
    }
    typ, err := strconv.Atoi(form.Get("type"))
    if err != nil {
        log.Print(err)
        return c.String(http.StatusOK, "Type vide")
    }
    year, err := strconv.Atoi(form.Get("year"))
    if err != nil {
        log.Print(err)
        return c.String(http.StatusOK, "Annee vide")
    }
    description := form.Get("description")
    if description == "" {
        err = fmt.Errorf("Description is empty")
        return c.String(http.StatusOK, "Description vide")
    }
    err = database.AddAnime(title, description, typ, year)
    if err != nil {
        log.Print(err)
        return c.String(http.StatusOK, "Erreur")
    }
    return c.String(http.StatusOK, "Ajouté !")
}
