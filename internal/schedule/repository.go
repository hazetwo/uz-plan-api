package schedule

import (
	"context"
	"uz-plan-api/internal/database"
	"uz-plan-api/internal/model"
)

type Repository interface {
	GetFields(ctx context.Context) (map[string]string, bool, error)
	StoreFields(ctx context.Context, fields map[string]string) error

	GetGroups(ctx context.Context, fieldID string) (map[string]string, bool, error)
	StoreGroups(ctx context.Context, fieldID string, groups map[string]string) error

	GetSchedule(ctx context.Context, groupID string) ([]model.Entry, bool, error)
	StoreSchedule(ctx context.Context, groupID string, entries []model.Entry) error
}

var _ Repository = database.RedisRepository{}
