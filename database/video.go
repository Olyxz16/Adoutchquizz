package database

import (
	"database/sql"
	"strings"
	"time"
)

type Video struct {
    Uid             int
    VideoID         int
    ReleaseDate     time.Time
	ClipRef         int
	ClipInd         int
    Ok              bool
}


func AddVideo(video Video) error {
    db := dbInstance.db
    tx, err := db.Begin()
	if err != nil {
		return err
	}
	q := `INSERT INTO 
	Video (videoID, releaseDate, clipID, clipInd, ok)
	VALUES ($1, $2, $3, $4, $5);`
    date := sql.NullTime { Time: video.ReleaseDate, Valid: true }
    _, err = tx.Exec(q, video.VideoID, date, video.ClipRef, video.ClipInd, video.Ok)
    if err != nil {
        return err
    }
	tx.Commit()
	return nil
}

func AddClipToVideo(videoId int, title string, typ, ind int) error {
    db := dbInstance.db
    id, err := ClipIdFromData(title, typ, ind)
    if err != nil {
        return err
    }
    var order int
    var releaseDate time.Time
    q := `SELECT releaseDate, MAX(clipInd) FROM Video
            WHERE videoID=$1
            GROUP BY releaseDate
            LIMIT 1`
    rows, err := db.Query(q, videoId)
    if err != nil {
        return err
    }
    contains := rows.Next()
    if !contains {
        order = 1
        // a changer
        releaseDate = time.Now()
    } else {
        err = rows.Scan(&releaseDate, &order)
        if err != nil {
            return err
        }
        order += 1
    }

    video := Video { VideoID: videoId, ReleaseDate: releaseDate, ClipRef: id, ClipInd: order, Ok: false }
    err = AddVideo(video)
    if err != nil {
        return err
    }
    return nil
}

func GetLatestVideos(n int) ([]Video, error) {
    db := dbInstance.db
    q := `SELECT DISTINCT videoID, releaseDate FROM Video
            ORDER BY releaseDate DESC
            LIMIT $1`
    rows, err := db.Query(q, n)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    videos := []Video{}
    for rows.Next() {
        var videoID int
        var releaseDate time.Time
        err := rows.Scan(&videoID, &releaseDate)
        if err != nil {
            return nil, err
        }
        vid := Video{
            VideoID: videoID,
            ReleaseDate: releaseDate,
        }
        videos = append(videos, vid)
    }
    return videos, nil
}

func GetAllClipsFromVideo(videoId int) ([]Video, []Clip, []Anime, error) {
    db := dbInstance.db
    qvideos := `SELECT * FROM Video
                FULL JOIN (
                    SELECT uid as clID, animeID, type, ind FROM Clip
                ) ON Video.clipID = clID
                FULL JOIN (
                    SELECT uid as aniID, title as aniTitle FROM Anime            
                ) ON animeID = aniID
                WHERE Video.videoID=$1
                ORDER BY Video.clipInd ASC`
    rows, err := db.Query(qvideos, videoId)
	if err != nil {
		return nil, nil, nil, err
	}
    defer rows.Close()
    
    videos := []Video{}
    clips := []Clip{}
    animes := []Anime{}

	for rows.Next() {
        var vuid int
        var videoId int
        var releaseDate time.Time
        var clipRef int
        var clipInd int
        var ok bool

        var clipId int
		var animeRef int 
		var typ int
        var ind int

        var aniId int
        var animeTitle string
		err := rows.Scan(
            &vuid, &videoId, &releaseDate, &clipRef, &clipInd, &ok,
            &clipId, &animeRef, &typ, &ind,
            &aniId, &animeTitle,
            )
		if err != nil {
			return nil, nil, nil, err
		}
        video := Video {
            Uid: vuid,
            VideoID: videoId,
            ReleaseDate: releaseDate,
            ClipRef: clipRef,
            ClipInd: clipInd,
            Ok: ok,
        }
        clip := Clip { 
            AnimeRef: animeRef,
            Type: typ,
            Ind: ind,
        }
        anime := Anime {
            Uid: aniId,
            Title: animeTitle,
        }
	    videos = append(videos, video)
        clips = append(clips, clip)
        animes = append(animes, anime)
	}

	return videos, clips, animes, nil
}

func GetNextVideoID() (int, error) {
    db := dbInstance.db
    q := `SELECT MAX(videoID) FROM Video`
    rows, err := db.Query(q)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
    
    rows.Next()
    var uid int
    err = rows.Scan(&uid)
    if err != nil {
        return 0, err
    }
    return uid+1, nil
}

func SetClipOKInVideo(uid int, ok bool) (int, error) {
    db := dbInstance.db
    q := `UPDATE Video SET ok=$2 
            WHERE uid=$1
            RETURNING videoID`
    var videoId int
    rows, err := db.Query(q, uid, ok)
    rows.Next()
    rows.Scan(&videoId)
    return videoId, err
}

func RemoveClipFromVideo(uid int) (error) {
    db := dbInstance.db
    q := `DELETE FROM Video WHERE uid=$1`
    _, err := db.Exec(q, uid)
    return err
}

func (s *service) migrateVideo() error {
    q := `CREATE TABLE IF NOT EXISTS Video (
        uid         SERIAL PRIMARY KEY,
        videoID     INT,
        releaseDate DATE,
        clipID      INT,
        clipInd     INT,
        ok          BOOL,
        UNIQUE (clipID)
    )`
    _, err := s.db.Exec(q)
    return err
}

func (s *service) dropVideo() error {
	q := `DROP TABLE Video`
	_, err := s.db.Exec(q)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 42P01") {
			return nil
		}
		return err
	}
	return nil
}
