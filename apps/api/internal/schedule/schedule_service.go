package schedule

import (
	"context"

	"github.com/go-redsync/redsync/v4"
)

type Service struct {
	scraper *Scraper
	repo    Repository
	rs      *redsync.Redsync
}

func NewService(scraper *Scraper, repo Repository, rs *redsync.Redsync) *Service {
	return &Service{scraper: scraper, repo: repo, rs: rs}
}

func (s Service) GetFields(ctx context.Context) (map[string]string, error) {
	f, ok, err := s.repo.GetFields(ctx)
	if err != nil {
		return nil, err
	}
	if ok {
		return f, nil
	}

	mu := s.rs.NewMutex("lock:fields")
	if err := mu.LockContext(ctx); err != nil {
		return nil, err
	}
	defer func() {
		_, err := mu.UnlockContext(ctx)
		if err != nil {
			return
		}
	}()

	f, err = s.scraper.GetFields(fieldsURL)
	if err != nil {
		return nil, err
	}

	if err := s.repo.StoreFields(ctx, f); err != nil {
		return nil, err
	}

	return f, nil
}

func (s Service) GetGroups(ctx context.Context, fieldsID string) (map[string]string, error) {
	//if g, ok := s.repo.GetGroups(ctx, fieldsID); ok {
	//	return g, nil
	//}

	g, err := s.scraper.GetGroupsFromID(groupsURL, fieldsID)
	if err != nil {
		return nil, err
	}
	//if err := s.repo.StoreGroups(ctx, fieldsID, g); err != nil {
	//	return nil, err
	//}

	return g, nil

}

func (s Service) getSchedule(ctx context.Context, groupID string) ([]Entry, error) {
	//if sh, ok := s.repo.GetSchedule(ctx, groupID); ok {
	//	return sh, nil
	//}

	sh, err := s.scraper.GetScheduleForID(scheduleURL, groupID)
	if err != nil {
		return nil, err
	}

	//err = s.repo.StoreSchedule(ctx, groupID, sh)
	//if err != nil {
	//	return nil, err
	//}

	return sh, nil

}

func (s Service) GetFilteredSchedule(ctx context.Context, groupID string, f Filter) ([]Entry, error) {
	entries, err := s.getSchedule(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var filtered = filterEntries(entries, dayPredicate(f.Day), weekPredicate(f.Week), subgroupPredicate(f.Subgroup))

	return filtered, nil
}
