package database

import (
	"database/sql"
	"strings"
	"time"
)

type Video struct {
    VideoID             int
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

func GetAllClipsFromVideo(uid int) ([]Video, []Clip, []Anime, error) {
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
    rows, err := db.Query(qvideos, uid)
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

func (s *service) migrateVideo() error {
    q := `CREATE TABLE IF NOT EXISTS Video (
        uid         SERIAL PRIMARY KEY,
        videoID     INT,
        releaseDate DATE,
        clipID      INT,
        clipInd     INT,
        ok          BOOL
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
