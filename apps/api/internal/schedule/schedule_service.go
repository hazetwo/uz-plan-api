package schedule

type Service struct {
	scraper Scraper
	repo    Repository
}

func NewService(scraper Scraper, repo Repository) Service {
	return Service{scraper: scraper, repo: repo}
}

func (s Service) GetFields(site string) (map[string]string, error) {
	if f, ok := s.repo.GetFields(); ok {
		return f, nil
	}

	f := s.scraper.GetFields(site)
	if err := s.repo.StoreFields(f); err != nil {
		return nil, err
	}

	return f, nil
}

func (s Service) GetGroups(site string, fieldsID string) (map[string]string, error) {
	if g, ok := s.repo.GetGroups(fieldsID); ok {
		return g, nil
	}

	g := s.scraper.GetGroupsFromID(site, fieldsID)
	if err := s.repo.StoreGroups(fieldsID, g); err != nil {
		return nil, err
	}

	return g, nil

}

func (s Service) GetSchedule(site string, groupID string) ([]Entry, error) {
	if sh, ok := s.repo.GetSchedule(groupID); ok {
		return sh, nil
	}

	sh, err := s.scraper.GetScheduleForID(site, groupID)
	if err != nil {
		return nil, err
	}

	err = s.repo.StoreSchedule(groupID, sh)
	if err != nil {
		return nil, err
	}

	return sh, nil

}
