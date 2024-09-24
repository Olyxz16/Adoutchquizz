package clip

import (
    "strings"
	"database/sql"
)

type Clip struct {
    AnimeRef        string
    Type            int
    Ind             int
    Year            int
    Title           string
    Url             string
    Path            string
    Usable          bool
    Difficulty      int
}

func Add(db *sql.DB, clip Clip) error {
    tx, err := db.Begin()
	if err != nil {
		return err
	}
	q := `INSERT INTO 
	Clip (animeID, type, ind, year, title, url, path, usable, difficulty)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
    _, err = tx.Exec(q, clip.AnimeRef, clip.Type, clip.Ind, clip.Year, clip.Title, clip.Url, clip.Path, clip.Usable, clip.Difficulty)
    if err != nil {
        return err
    }
	tx.Commit()
	return nil
}

func Migrate(db *sql.DB) error {
    q := `CREATE TABLE IF NOT EXISTS Clip (
        animeID         VARCHAR(255),
        type            VARCHAR(255),
        ind             VARCHAR(255),
        year            INT,
        title           VARCHAR(255),
        url             VARCHAR(255),
        path            VARCHAR(255),
        usable          BOOLEAN,
        difficulty      INT
        PRIMARY KEY (animeID, type, ind)
    )`
    _, err := db.Exec(q)
    return err
}

func Drop(db *sql.DB) error {
	q := `DROP TABLE Clip`
	_, err := db.Exec(q)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 42P01") {
			return nil
		}
		return err
	}
	return nil
}
