package database

import (
	"fmt"
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
	Anime (title, year, type, description)
	VALUES ($1, $2, $3, $4);`
    _, err = tx.Exec(q, anime.Title, anime.Year, anime.Type, anime.Description)
    if err != nil {
        return err
    }
	tx.Commit()
	return nil
}

func GetAnimeIDFromTitle(title string) (int, error) {
    db := dbInstance.db
    q := `SELECT uid FROM Anime
        WHERE title=$1`
    rows, err := db.Query(q, title)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
    
    ok := rows.Next()
    if !ok {
        return 0, fmt.Errorf("Not enough result")
    }
    var uid int
    err = rows.Scan(&uid)
    if err != nil {
        return 0, err
    }
    return uid, nil
}

func GetAllClipsFromAnime(title string) ([]Clip, error) {
    db := dbInstance.db
    result := []Clip{}
    q := `SELECT * FROM Clip
          WHERE clip.animeID IN 
            (SELECT DISTINCT uid FROM Anime
            WHERE title=$1)`
    rows, err := db.Query(q, title)
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
                    Uid: uid, Ind: ind, Year: year,
                    Title: title, Url: url,
                    Path: path, Usable: usable,
                    Difficulty: difficulty,
                }
	    result = append(result, clip)
	}

	return result, nil
}

func GetAllUsableClipsFromAnime(title string) ([]Clip, error) {
    db := dbInstance.db
    result := []Clip{}
    q := `SELECT * FROM Clip
          WHERE usable=true
          AND clip.animeID IN 
            (SELECT DISTINCT uid FROM Anime
            WHERE title=$1)`

    rows, err := db.Query(q, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var animeRef int 
		var typ int
        var ind int
        var year int
        var title string
        var url string
        var path string
        var usable bool
        var difficulty int
		err := rows.Scan(&animeRef, &typ, &ind, &year, &title, &url, &path, &usable, &difficulty)
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
