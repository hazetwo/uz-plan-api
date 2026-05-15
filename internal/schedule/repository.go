package schedule

import (
	"context"
	"uz-plan-api/internal/database"
	"uz-plan-api/internal/model"
)

type Repository interface {
	GetFields(ctx context.Context) (model.Fields, database.Cached, error)
	StoreFields(ctx context.Context, fields model.Fields) error

	GetGroups(ctx context.Context, fieldID string) (model.Groups, database.Cached, error)
	StoreGroups(ctx context.Context, fieldID string, groups model.Groups) error

	GetSchedule(ctx context.Context, groupID string) ([]model.Entry, database.Cached, error)
	StoreSchedule(ctx context.Context, groupID string, entries []model.Entry) error
}

var _ Repository = database.RedisRepository{}
