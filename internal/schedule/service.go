package schedule

import (
	"context"
	"uz-plan-api/internal/errs"
	"uz-plan-api/internal/model"
	"uz-plan-api/internal/scraper"

	"github.com/go-redsync/redsync/v4"
)

type Service struct {
	scraper *scraper.Scraper
	repo    Repository
	rs      *redsync.Redsync
}

func NewService(scraper *scraper.Scraper, repo Repository, rs *redsync.Redsync) *Service {
	return &Service{scraper: scraper, repo: repo, rs: rs}
}

func (s Service) GetFields(ctx context.Context) (map[string]string, error) {
	f, ok, err := s.repo.GetFields(ctx)
	if err != nil {
		return nil, errs.ErrFetchFailed
	}
	if ok {
		return f, nil
	}

	mu := s.rs.NewMutex("lock:fields")
	if err := mu.LockContext(ctx); err != nil {
		return nil, errs.ErrFetchFailed
	}
	defer func() {
		_, err := mu.UnlockContext(ctx)
		if err != nil {
			return
		}
	}()

	f, err = s.scraper.GetFields(scraper.FieldsURL)
	if err != nil {
		return nil, errs.ErrFetchFailed
	}
	if len(f) == 0 {
		return nil, errs.ErrNotFound
	}

	if err := s.repo.StoreFields(ctx, f); err != nil {
		return nil, errs.ErrFetchFailed
	}

	return f, nil
}

func (s Service) GetGroups(ctx context.Context, fieldsID string) (map[string]string, error) {
	g, ok, err := s.repo.GetGroups(ctx, fieldsID)
	if err != nil {
		return nil, errs.ErrFetchFailed
	}
	if ok {
		return g, nil
	}

	mu := s.rs.NewMutex("lock:group:" + fieldsID)
	if err := mu.LockContext(ctx); err != nil {
		return nil, errs.ErrFetchFailed
	}
	defer func() {
		_, err := mu.UnlockContext(ctx)
		if err != nil {
			return
		}
	}()

	g, err = s.scraper.GetGroupsFromID(scraper.GroupsURL, fieldsID)
	if err != nil {
		return nil, errs.ErrFetchFailed
	}
	if len(g) == 0 {
		return nil, errs.ErrNotFound
	}

	if err := s.repo.StoreGroups(ctx, fieldsID, g); err != nil {
		return nil, errs.ErrFetchFailed
	}

	return g, nil

}

func (s Service) getSchedule(ctx context.Context, groupID string) ([]model.Entry, error) {
	sh, ok, err := s.repo.GetSchedule(ctx, groupID)
	if err != nil {
		return nil, errs.ErrFetchFailed
	}
	if ok {
		return sh, nil
	}

	mu := s.rs.NewMutex("lock:schedule:" + groupID)
	if err := mu.LockContext(ctx); err != nil {
		return nil, errs.ErrFetchFailed
	}
	defer func() {
		_, err := mu.UnlockContext(ctx)
		if err != nil {
			return
		}
	}()

	sh, err = s.scraper.GetScheduleForID(scraper.ScheduleURL, groupID)
	if err != nil {
		return nil, errs.ErrFetchFailed
	}
	if len(sh) == 0 {
		return nil, errs.ErrNotFound
	}

	err = s.repo.StoreSchedule(ctx, groupID, sh)
	if err != nil {
		return nil, errs.ErrFetchFailed
	}

	return sh, nil

}

func (s Service) GetFilteredSchedule(ctx context.Context, groupID string, f model.Filter) ([]model.Entry, error) {
	entries, err := s.getSchedule(ctx, groupID)
	if err != nil {
		return nil, errs.ErrFetchFailed
	}

	var filtered = model.FilterEntries(
		entries, model.DayPredicate(f.Day),
		model.WeekPredicate(f.Week),
		model.SubgroupPredicate(f.Subgroup))

	return filtered, nil
}
