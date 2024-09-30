package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"Adoutchquizz/database"
	"Adoutchquizz/views/setviews"
    "Adoutchquizz/views/errors"
)

func Set(c echo.Context) error {
    group := c.Param("value")
    switch group {
        case "anime": return render(c, setviews.Anime()) 
        case "clip": return render(c, setviews.Clip())
        case "video": return render(c, setviews.Video())
    }
    return render(c, errors.Error404())
}

func SetForm(c echo.Context) error {
    group := c.Param("value")
    switch group {
        case "anime": return postAnime(c)
        case "clip": return postClip(c)
        case "video": return postVideo(c)
    }
    return render(c, errors.Error404())
}

func postVideo(c echo.Context) error {
    title := c.FormValue("title")
    typ, err1 := strconv.Atoi(c.FormValue("type"))
    indClip, err2 := strconv.Atoi(c.FormValue("indClip"))
    indVideo, err3 := strconv.Atoi(c.FormValue("indVideo"))
    date := c.FormValue("date")
    
    if err1 != nil {
        return render(c, setviews.VideoResult(false))
    }
    if err2 != nil {
        return render(c, setviews.VideoResult(false))
    }
    if err3 != nil {
        return render(c, setviews.VideoResult(false))
    }

    uid, err := database.ClipIdFromData(title, typ, indClip);
    if err != nil {
        return render(c, setviews.VideoResult(false))
    }
    
    releaseDate, err := time.Parse("2006-01-02", date)
    if err != nil {
        return render(c, setviews.VideoResult(false))
    }
    
    fmt.Println(indVideo)
    video := database.Video {
        ClipRef: uid,
        ReleaseDate: releaseDate,
    }
    err = database.AddVideo(video) 
    if err != nil {
        return render(c, setviews.VideoResult(false))
    }
    
    return render(c, setviews.VideoResult(true))
}

func postAnime(c echo.Context) error {
    title := c.FormValue("title")
    year, err := strconv.Atoi(c.FormValue("year"))
    if err != nil {
        return render(c, setviews.AnimeResult(false))
    }
    typ := c.FormValue("type")
    desc := c.FormValue("description")
    if title == "" || year < 0 || typ == "" || desc == "" {
        return render(c, setviews.AnimeResult(false))
    }
    anime := database.Anime{
        Title: title, Year: year, Type: typ, Description: desc,
    }
    err = database.AddAnime(anime)
    if err != nil {
        fmt.Println(err)
        return render(c, setviews.AnimeResult(false))
    }
    return render(c, setviews.AnimeResult(true))
}

func postClip(c echo.Context) error {
    animetitle := c.FormValue("title")
    typ := c.FormValue("type")
    ind, err1 := strconv.Atoi(c.FormValue("ind"))
    year, err2 := strconv.Atoi(c.FormValue("year"))
    title := c.FormValue("title")
    url := c.FormValue("url")
    path := c.FormValue("path")
    usablestr := c.FormValue("usable")
    difficulty, err3 := strconv.Atoi(c.FormValue("difficulty"))
    
    if err1 != nil || err2 != nil || err3 != nil {
        return render(c, setviews.ClipResult(false))
    }
    animeref, err := database.GetAnimeIDFromTitle(animetitle)
    if err != nil {
        return render(c, setviews.ClipResult(false))
    }
    typeInd, err := clipTypeAtoi(typ)
    if err != nil {
        return render(c, setviews.ClipResult(false))
    }
    usable := false
    if usablestr == "true" {
        usable = true
    }
    if animetitle == "" || title == "" || url == "" || path == "" {
        return render(c, setviews.ClipResult(false))
    }

    clip := database.Clip {
        AnimeRef: animeref, Type: typeInd,
        Ind: ind, Year: year, Title: title,
        Url: url, Path: path,
        Usable: usable, Difficulty: difficulty,
    }
    err = database.AddClip(clip)
    if err != nil {
        fmt.Println(err)
        return render(c, setviews.AnimeResult(false))
    }
    return render(c, setviews.AnimeResult(true))
}

func clipTypeAtoi(typ string) (int, error) {
    switch strings.ToLower(typ) {
    case "opening": return 1, nil
    case "ending": return 2, nil
    case "ost": return 3, nil
    }
    return 0, fmt.Errorf("Wrong parameter")
}
