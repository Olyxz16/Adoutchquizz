package database

import (
    "strings"
)
type Anime struct {
    Uid             int
    Title           string
    Year            int
    Type            string
    Description     string
}

func AddAnime(anime Anime) error {
    db := dbInstance.db
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

func (s *service) migrateAnime() error {
    q := `CREATE TABLE IF NOT EXISTS Anime (
        uid         SERIAL PRIMARY KEY,
        title       VARCHAR(255),
        year        INT,
        type        VARCHAR(255),
        description VARCHAR(255)
    )`
    _, err := s.db.Exec(q)
    return err
}

func (s *service) dropAnime() error {
	q := `DROP TABLE Anime`
	_, err := s.db.Exec(q)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 42P01") {
			return nil
		}
		return err
	}
	return nil
}
