package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"Adoutchquizz/database"
	"Adoutchquizz/views/getviews"
	"Adoutchquizz/views/errors"
)

func Get(c echo.Context) error {
    group := c.Param("group")
    switch group {
        case "video": return getVideo(c)
        case "anime": return getAnime(c) 
    }
    return render(c, errors.Error404())
}

/*********************************/
/*            SEARCH             */
/*********************************/

func getSearchView(c echo.Context, desc string) error {
    req := c.Request().URL.RequestURI()
    return render(c, getviews.Search(req, desc))
}


/*********************************/
/*            VIDEO              */
/*********************************/

func getVideo(c echo.Context) error {
    value := c.Param("value")
    switch value {
        case "clip": return getVideoClip(c)
    }
    return render(c, errors.Error404())
}


func getVideoClip(c echo.Context) error {
    url := c.QueryParam("search")
    if url == "" {
        return getSearchView(c, "Tous les clips d'une vidéo")
    }
    clips, err := database.GetAllClipsFromVideo(url)
    if err != nil {
        return render(c, getviews.VideoTableError())
    }
    return render(c, getviews.VideoTable(url, clips)) 
}


/*********************************/
/*            ANIME              */
/*********************************/

func getAnime(c echo.Context) error {
    value := c.Param("value")
    switch value {
        case "clip": return getAnimeClip(c)
        case "usable": return getAnimeUsableClip(c)
    }
    return render(c, errors.Error404())
}


func getAnimeClip(c echo.Context) error {
    title := c.QueryParam("search")
    if title == "" {
        return getSearchView(c, "Tous les clips d'un anime")
    }
    clips, err := database.GetAllClipsFromAnime(title)
    if err != nil {
        fmt.Println(err)
        return render(c, getviews.AnimeTableError())
    }
    return render(c, getviews.AnimeTable(title, clips))
}
func getAnimeUsableClip(c echo.Context) error {
    title := c.QueryParam("search")
    if title == "" {
        return getSearchView(c, "Tous les clips utilisables d'un anime")
    }
    clips, err := database.GetAllUsableClipsFromAnime(title)
    if err != nil {
        fmt.Println(err)
        return render(c, getviews.AnimeTableError())
    }
    return render(c, getviews.AnimeTable(title, clips))
}
