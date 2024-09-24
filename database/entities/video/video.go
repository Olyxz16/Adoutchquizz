package video

import (
    "strings"
	"database/sql"
)

type Video struct {
	Url             string
	ClipRef         string
	Ind             int
    ReleaseDate     sql.NullTime
}

func Add(db *sql.DB, video Video) error {
    tx, err := db.Begin()
	if err != nil {
		return err
	}
	q := `INSERT INTO 
	Video (url, clipID, ind, releaseDate)
	VALUES ($1, $2, $3, $4);`
    _, err = tx.Exec(q, video.Url, video.ClipRef, video.Ind, video.ReleaseDate)
    if err != nil {
        return err
    }
	tx.Commit()
	return nil
}

func Migrate(db *sql.DB) error {
    q := `CREATE TABLE IF NOT EXISTS Video (
        url         VARCHAR(255),
        clipID      VARCHAR(255),
        ind         INT,
        releaseDate DATE,
        PRIMARY KEY (url, clipID)
    )`
    _, err := db.Exec(q)
    return err
}

func Drop(db *sql.DB) error {
	q := `DROP TABLE Video`
	_, err := db.Exec(q)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 42P01") {
			return nil
		}
		return err
	}
	return nil
}
