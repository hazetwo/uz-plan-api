package schedule

import (
	"context"
	"log/slog"
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

func (s Service) GetFields(ctx context.Context) (model.Fields, error) {
	f, ok, err := s.repo.GetFields(ctx)
	if err != nil {
		return nil, errs.FetchFailed(ctx, err)
	}
	if ok {
		return f, nil
	}

	var result map[string]string
	if err := s.withLock(ctx, "lock:fields", func() error {
		cached, ok, err := s.repo.GetFields(ctx)
		if err != nil {
			return errs.FetchFailed(ctx, err)
		}
		if ok {
			result = cached
			return nil
		}
		result, err = s.scraper.GetFields(scraper.FieldsURL)
		if err != nil {
			return errs.FetchFailed(ctx, err)
		}
		if len(result) == 0 {
			return errs.ErrNotFound
		}
		return s.repo.StoreFields(ctx, result)
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s Service) GetGroups(ctx context.Context, fieldsID string) (model.Groups, error) {
	g, ok, err := s.repo.GetGroups(ctx, fieldsID)
	if err != nil {
		return nil, errs.FetchFailed(ctx, err)
	}
	if ok {
		return g, nil
	}

	var result map[string]string
	if err := s.withLock(ctx, "lock:group:"+fieldsID, func() error {
		cached, ok, err := s.repo.GetGroups(ctx, fieldsID)
		if err != nil {
			return errs.FetchFailed(ctx, err)
		}
		if ok {
			result = cached
			return nil
		}
		ids, err := s.GetFields(ctx)
		if err != nil {
			return err
		}
		result, err = s.scraper.GetGroupsFromID(scraper.GroupsURL, fieldsID, ids)
		if err != nil {
			return errs.FetchFailed(ctx, err)
		}
		if len(result) == 0 {
			return errs.ErrNotFound
		}
		return s.repo.StoreGroups(ctx, fieldsID, result)
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s Service) getSchedule(ctx context.Context, groupID string) ([]model.Entry, error) {
	sh, ok, err := s.repo.GetSchedule(ctx, groupID)
	if err != nil {
		return nil, errs.FetchFailed(ctx, err)
	}
	if ok {
		return sh, nil
	}

	var result []model.Entry
	if err := s.withLock(ctx, "lock:schedule:"+groupID, func() error {
		cached, ok, err := s.repo.GetSchedule(ctx, groupID)
		if err != nil {
			return errs.FetchFailed(ctx, err)
		}
		if ok {
			result = cached
			return nil
		}
		result, err = s.scraper.GetScheduleForID(scraper.ScheduleURL, groupID)
		if err != nil {
			return errs.FetchFailed(ctx, err)
		}
		if len(result) == 0 {
			return errs.ErrNotFound
		}
		return s.repo.StoreSchedule(ctx, groupID, result)
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s Service) GetFilteredSchedule(ctx context.Context, groupID string, f model.Filter) ([]model.Entry, error) {
	entries, err := s.getSchedule(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var filtered = model.FilterEntries(
		entries, model.DayPredicate(f.Day),
		model.WeekPredicate(f.Week),
		model.SubgroupPredicate(f.Subgroup))

	return filtered, nil
}

func (s Service) withLock(ctx context.Context, key string, fn func() error) error {
	mu := s.rs.NewMutex(key)
	if err := mu.LockContext(ctx); err != nil {
		return errs.FetchFailed(ctx, err)
	}
	defer func() {
		if _, err := mu.UnlockContext(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to release lock", "key", key, "err", err)
		}
	}()

	return fn()
}
