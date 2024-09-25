package database

func (s *service) Migrate() error {
	if err := s.migrateAnime(); err != nil {
		return err
	}
	if err := s.migrateClip(); err != nil {
		return err
	}
	if err := s.migrateVideo(); err != nil {
		return err
	}
	return nil
}

func (s *service) Drop() error {
	if err := s.dropAnime(); err != nil {
		return err
	}
	if err := s.dropClip(); err != nil {
		return err
	}
	if err := s.dropVideo(); err != nil {
		return err
	}
	return nil
}
