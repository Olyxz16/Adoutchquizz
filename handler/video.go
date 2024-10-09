package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"Adoutchquizz/database"
	"Adoutchquizz/views/errors"
	"Adoutchquizz/views/video"
)

/*
*    /video?uid={int}   |   /video/new
 */
func Video(c echo.Context) error {
    url := c.Request().URL.Path
    switch url {
        case "/video": return videoPage(c)
        case "/video/new": return newVideoPage(c)
        default: return render(c, errors.Error404())
    }
}

func newVideoPage(c echo.Context) error {
    id, err := database.GetNextVideoID()
    if err != nil {
        return err
    }
    return render(c, video.Layout(id, []video.ClipData{}, []video.ClipData{}, []video.ClipData{}))
}


func videoPage(c echo.Context) error {
    id, err := strconv.Atoi(c.QueryParam("uid"))
    if err != nil {
        // Si l'id n'est pas bon, on retourne une page vierge
        return newVideoPage(c)
    }
    videos, clips, animes, err := database.GetAllClipsFromVideo(id)
    if err != nil {
        return err
    }
    opening, ending, ost := sortVideos(videos, clips, animes)
    return render(c, video.Layout(id, opening, ending, ost))
}

func sortVideos(videos []database.Video, clips []database.Clip, animes []database.Anime) ([]video.ClipData, []video.ClipData, []video.ClipData) {
    opening := []video.ClipData{}
    ending := []video.ClipData{}
    ost := []video.ClipData{}
    for i, v := range videos {
        c := clips[i]
        a := animes[i]
        data := video.ClipData {
            VideoInd: v.ClipInd,
            AnimeTitle: a.Title,
            ClipInd: c.Ind,
            State: v.Ok, 
        }
        switch clips[i].Type {
            case 1: opening = append(opening, data); break;
            case 2: ending = append(ending, data); break;
            case 3: ost = append(ost, data); break;
        }
    }
    return opening, ending, ost
}
