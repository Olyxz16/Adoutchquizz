package database

import (
	"database/sql"
	"strings"
	"time"
)

type Video struct {
    Uid             int
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
	Video (uid, releaseDate, clipID, clipInd, ok)
	VALUES ($1, $2, $3, $4, $5);`
    date := sql.NullTime { Time: video.ReleaseDate, Valid: true }
    _, err = tx.Exec(q, video.Uid, date, video.ClipRef, video.ClipInd, video.Ok)
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
        uid         SERIAL PRIMARY KEY,
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
