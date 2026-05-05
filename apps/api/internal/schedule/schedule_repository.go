package schedule

import "context"

type Repository interface {
	GetFields(ctx context.Context) (map[string]string, bool, error)
	StoreFields(ctx context.Context, fields map[string]string) error

	GetGroups(ctx context.Context, fieldID string) (map[string]string, bool, error)
	StoreGroups(ctx context.Context, fieldID string, groups map[string]string) error

	GetSchedule(ctx context.Context, groupID string) ([]Entry, bool, error)
	StoreSchedule(ctx context.Context, groupID string, entries []Entry) error
}

var _ Repository = RedisRepository{}
