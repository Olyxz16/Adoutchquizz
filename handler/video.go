package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"Adoutchquizz/database"
	"Adoutchquizz/views/video"
)

/*
*    /video?uid={int}   |   /video/new
 */
func Video(c echo.Context) error {
    url := c.Path()
    switch url {
    case "/video/:uid": return videoPage(c)
    case "/video/new": return newVideoPage(c)
    default: return Error404(c)
    }
}
func VideoPost(c echo.Context) error {
    par := c.Param("action")
    switch par {
    case "setok": postClipState(c); break;
    case "addclip": addClip(c); break;
    default: return Error404(c)
    }
    return nil
}
func VideoDelete(c echo.Context) error {
    url := c.Path()
    switch url {
    case "/video/clip/:uid": return deleteClip(c)
    default: return Error404(c)
    }
}


/************************************/
/*               GET                */
/************************************/

func newVideoPage(c echo.Context) error {
    id, err := database.GetNextVideoID()
    if err != nil {
        log.Error(err)
        return err
    }
    return render(c, video.Layout(id, nil, []video.ClipData{}, []video.ClipData{}, []video.ClipData{}))
}


func videoPage(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("uid"))
    if err != nil {
        log.Warn(err)
        // Si l'id n'est pas bon, on retourne une page vierge
        return newVideoPage(c)
    }
    videos, clips, animes, err := database.GetAllClipsFromVideo(id)
    if err != nil {
        log.Error(err)
        return err
    }
    time := videos[0].ReleaseDate
    opening, ending, ost := sortVideos(videos, clips, animes)
    return render(c, video.Layout(id, &time, opening, ending, ost))
}
func videoColumns(c echo.Context, videoId int) error {
    videos, clips, animes, err := database.GetAllClipsFromVideo(videoId)
    if err != nil {
        log.Error(err)
        return err
    }
    opening, ending, ost := sortVideos(videos, clips, animes)
    return render(c, video.Cols(videoId, opening, ending, ost))
}

func sortVideos(videos []database.Video, clips []database.Clip, animes []database.Anime) ([]video.ClipData, []video.ClipData, []video.ClipData) {
    opening := []video.ClipData{}
    ending := []video.ClipData{}
    ost := []video.ClipData{}
    for i, v := range videos {
        c := clips[i]
        a := animes[i]
        data := video.ClipData {
            Uid: v.Uid,
            VideoId: videos[0].VideoID,
            Order: v.ClipInd,
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


/************************************/
/*              POST                */
/************************************/
func postClipState(c echo.Context) error {
    var ok bool
    if c.FormValue("ok") == "true" {
        ok = true
    } else {
        ok = false
    }
    uid, err := strconv.Atoi(c.FormValue("uid"))
    if err != nil {
        log.Error(err)
        return err
    }
    videoId, err := database.SetClipOKInVideo(uid, ok)
    if err != nil {
        log.Error(err)
        return err
    }
    return videoColumns(c, videoId)
}


/* Returns the updated columns */
func addClip(c echo.Context) error {
    var errfun func(err error, displayText string) error
    errfun = func(err error, displayText string) error {
        log.Error(err)
        c.Response().Header().Set("HX-Retarget", "#adderror")
        c.Response().Header().Set("HX-Reswap", "innerHTML")
        c.String(http.StatusOK, displayText)
        return err
    }
    form, err := c.FormParams()
    if err != nil {
        return errfun(err, "Error reading form")
    }
    videoId, err := strconv.Atoi(form.Get("videoID"))
    if err != nil {
        return errfun(err, "Id manquant")
    }
    title := form.Get("title")
    if title == "" {
        err = fmt.Errorf("Title is empty")
        return errfun(err, "Titre manquant")
    }
    typ, err := strconv.Atoi(form.Get("type"))
    if err != nil {
        return errfun(err, "Type manquant")
    }
    ind, err := strconv.Atoi(form.Get("ind"))
    if err != nil {
        return errfun(err, "Indice manquant")
    }
    err = database.AddClipToVideo(videoId, title, typ, ind)
    if err != nil {
        return errfun(err, "Clip inexistant")
    }
    return videoColumns(c, videoId)
}

/************************************/
/*             DELETE               */
/************************************/

func deleteClip(c echo.Context) error {
    uid, err := strconv.Atoi(c.Param("uid")) 
    if err != nil {
        log.Error(err)
        return err
    }
    err = database.RemoveClipFromVideo(uid)
    if err != nil {
        log.Error(err)
        return err
    }
    return nil
}
