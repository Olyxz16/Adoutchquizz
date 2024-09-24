package anime

import (
    "strings"
	"database/sql"
)
type Anime struct {
    Uid             string
    Title           string
    Year            int
    Type            string
    Description     string
}

func Add(db *sql.DB, anime Anime) error {
    tx, err := db.Begin()
	if err != nil {
		return err
	}
	q := `INSERT INTO 
	Anime (uid, title, year, type, description)
	VALUES ($1, $2, $3, $4, $5);`
    _, err = tx.Exec(q, anime.Uid, anime.Title, anime.Year, anime.Type, anime.Description)
    if err != nil {
        return err
    }
	tx.Commit()
	return nil
}

func Migrate(db *sql.DB) error {
    q := `CREATE TABLE IF NOT EXISTS Anime (
        uid         VARCHAR(255) PRIMARY KEY,
        title       VARCHAR(255),
        year        INT,
        type        VARCHAR(255),
        description VARCHAR(255)
    )`
    _, err := db.Exec(q)
    return err
}

func Drop(db *sql.DB) error {
	q := `DROP TABLE Anime`
	_, err := db.Exec(q)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 42P01") {
			return nil
		}
		return err
	}
	return nil
}
