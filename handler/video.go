package handler

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"Adoutchquizz/database"
	"Adoutchquizz/views/errors"
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
    default: return render(c, errors.Error404())
    }
}
func VideoPost(c echo.Context) error {
    par := c.Param("action")
    switch par {
    case "setok": postClipState(c); break;
    case "addclip": addClip(c); break;
    default: return render(c, errors.Error404())
    }
    return nil
}
func VideoDelete(c echo.Context) error {
    url := c.Path()
    switch url {
    case "/video/clip/:uid": return deleteClip(c)
    default: return render(c, errors.Error404())
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
    vid, err := strconv.Atoi(c.FormValue("videoID"))
    if err != nil {
        log.Error(err)
        return err
    }
    cord, err := strconv.Atoi(c.FormValue("clipOrder"))
    if err != nil {
        log.Error(err)
        return err
    }
    data := video.ClipData {
        VideoId: vid,
        Order: cord,
        State: ok,
    } 
    err = database.SetClipOKInVideo(vid, cord, ok)
    if err != nil {
        log.Error(err)
        return err
    }
    return render(c, video.StateCheckBox(data))
}

func addClip(c echo.Context) error {
    form, err := c.FormParams()
    if err != nil {
        log.Error(err)
        return err
    }
    log.Print(form)
    videoId, err := strconv.Atoi(form.Get("videoId"))
    if err != nil {
        log.Error(err)
        return err
    }
    title := form.Get("title")
    if title == "" {
        err = fmt.Errorf("Title is empty")
        log.Error(err)
        return err
    }
    typ, err := strconv.Atoi(form.Get("type"))
    if err != nil {
        log.Error(err)
        return err
    }
    ind, err := strconv.Atoi(form.Get("ind"))
    if err != nil {
        log.Error(err)
        return err
    }
    err = database.AddClipToVideo(videoId, title, typ, ind)
    if err != nil {
        log.Error(err)
        return err
    }
    return nil
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
