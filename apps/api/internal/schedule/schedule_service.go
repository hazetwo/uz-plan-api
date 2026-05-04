package schedule

import "context"

type Service struct {
	scraper Scraper
	repo    Repository
}

func NewService(scraper Scraper, repo Repository) Service {
	return Service{scraper: scraper, repo: repo}
}

func (s Service) GetFields(ctx context.Context, site string) (map[string]string, error) {
	if f, ok := s.repo.GetFields(ctx); ok {
		return f, nil
	}

	f := s.scraper.GetFields(site)
	if err := s.repo.StoreFields(ctx, f); err != nil {
		return nil, err
	}

	return f, nil
}

func (s Service) GetGroups(ctx context.Context, site string, fieldsID string) (map[string]string, error) {
	if g, ok := s.repo.GetGroups(ctx, fieldsID); ok {
		return g, nil
	}

	g := s.scraper.GetGroupsFromID(site, fieldsID)
	if err := s.repo.StoreGroups(ctx, fieldsID, g); err != nil {
		return nil, err
	}

	return g, nil

}

func (s Service) GetSchedule(ctx context.Context, site string, groupID string) ([]Entry, error) {
	if sh, ok := s.repo.GetSchedule(ctx, groupID); ok {
		return sh, nil
	}

	sh, err := s.scraper.GetScheduleForID(site, groupID)
	if err != nil {
		return nil, err
	}

	err = s.repo.StoreSchedule(ctx, groupID, sh)
	if err != nil {
		return nil, err
	}

	return sh, nil

}
