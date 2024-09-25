package database

import (
    "strings"
)

type Clip struct {
    Uid             string
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



func AddClip(clip Clip) error {
    db := dbInstance.db
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

func (s *service) migrateClip() error {
    q := `CREATE TABLE IF NOT EXISTS Clip (
        uid             SERIAL PRIMARY KEY,
        animeID         VARCHAR(255),
        type            VARCHAR(255),
        ind             VARCHAR(255),
        year            INT,
        title           VARCHAR(255),
        url             VARCHAR(255),
        path            VARCHAR(255),
        usable          BOOLEAN,
        difficulty      INT
    )`
    _, err := s.db.Exec(q)
    return err
}

func (s *service) dropClip() error {
	q := `DROP TABLE Clip`
	_, err := s.db.Exec(q)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 42P01") {
			return nil
		}
		return err
	}
	return nil
}
