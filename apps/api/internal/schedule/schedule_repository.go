package schedule

import "context"

type Repository interface {
	GetFields(ctx context.Context) (map[string]string, bool)
	StoreFields(ctx context.Context, fields map[string]string) error

	GetGroups(ctx context.Context, fieldID string) (map[string]string, bool)
	StoreGroups(ctx context.Context, fieldID string, groups map[string]string) error

	GetSchedule(ctx context.Context, groupID string) ([]Entry, bool)
	StoreSchedule(ctx context.Context, groupID string, entries []Entry) error
}
