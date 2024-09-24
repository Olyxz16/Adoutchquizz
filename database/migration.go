package database

import (
    "Adoutchquizz/database/entities/anime"
    "Adoutchquizz/database/entities/clip"
    "Adoutchquizz/database/entities/video"
)

func (s *service) Migrate() error {
    db := dbInstance.db
	if err := anime.Migrate(db); err != nil {
		return err
	}
	if err := clip.Migrate(db); err != nil {
		return err
	}
	if err := video.Migrate(db); err != nil {
		return err
	}
	return nil
}

func (s *service) Drop() error {
    db := dbInstance.db
	if err := anime.Drop(db); err != nil {
		return err
	}
	if err := clip.Drop(db); err != nil {
		return err
	}
	if err := video.Drop(db); err != nil {
		return err
	}
	return nil
}
