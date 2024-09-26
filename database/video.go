package database

import (
	"database/sql"
	"strings"
	"time"
)

type Video struct {
	Url             string
	ClipRef         int
	Ind             int
    ReleaseDate     time.Time
}


func AddVideo(video Video) error {
    db := dbInstance.db
    tx, err := db.Begin()
	if err != nil {
		return err
	}
	q := `INSERT INTO 
	Video (url, clipID, ind, releaseDate)
	VALUES ($1, $2, $3, $4);`
    time := sql.NullTime { Time: video.ReleaseDate, Valid: true }
    _, err = tx.Exec(q, video.Url, video.ClipRef, video.Ind, time)
    if err != nil {
        return err
    }
	tx.Commit()
	return nil
}

func GetAllClipsFromVideo(url string) ([]Clip, error) {
    db := dbInstance.db
    result := []Clip{}
    q := `
        SELECT * FROM Video
        WHERE clipId IN 
            (SELECT clipID FROM Video
            WHERE url=$1)
    `
    rows, err := db.Query(q, url)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
        var uid int
		var animeRef int 
		var typ int
        var ind int
        var year int
        var title string
        var url string
        var path string
        var usable bool
        var difficulty int
		err := rows.Scan(&uid, &animeRef, &typ, &ind, &year, &title, &url, &path, &usable, &difficulty)
		if err != nil {
			return nil, err
		}
        clip := Clip{ 
                    AnimeRef: animeRef, Type: typ,
                    Ind: ind, Year: year,
                    Title: title, Url: url,
                    Path: path, Usable: usable,
                    Difficulty: difficulty,
                }
	    result = append(result, clip)
	}

	return result, nil
}

func (s *service) migrateVideo() error {
    q := `CREATE TABLE IF NOT EXISTS Video (
        url         VARCHAR(255),
        clipID      INT,
        ind         INT,
        releaseDate DATE,
        PRIMARY KEY (url, clipID)
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
