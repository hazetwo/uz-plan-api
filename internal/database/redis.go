package database

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"uz-plan-api/internal/model"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	rdb *redis.Client
}

type Cached bool

func New(rdb *redis.Client) (*RedisRepository, *redsync.Redsync) {
	pool := goredis.NewPool(rdb)
	return &RedisRepository{rdb: rdb}, redsync.New(pool)
}

func (r RedisRepository) GetFields(ctx context.Context) (model.Fields, Cached, error) {
	data, err := r.rdb.Get(ctx, "fields").Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var fields model.Fields
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, false, err
	}

	return fields, true, nil
}

func (r RedisRepository) StoreFields(ctx context.Context, fields model.Fields) error {
	data, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, "fields", data, 24*time.Hour).Err()
}

func (r RedisRepository) GetGroups(ctx context.Context, fieldID string) (model.Groups, Cached, error) {
	g, err := r.rdb.HGetAll(ctx, "groups:"+fieldID).Result()
	if err != nil {
		return nil, false, err
	}
	if len(g) == 0 {
		return nil, false, nil
	}

	return g, true, nil
}

func (r RedisRepository) StoreGroups(ctx context.Context, fieldID string, groups model.Groups) error {
	p := r.rdb.Pipeline()
	p.HSet(ctx, "groups:"+fieldID, map[string]string(groups))
	p.Expire(ctx, "groups:"+fieldID, 24*time.Hour)
	if _, err := p.Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r RedisRepository) GetSchedule(ctx context.Context, groupID string) ([]model.Entry, Cached, error) {
	data, err := r.rdb.Get(ctx, "schedule:"+groupID).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var entries []model.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, false, err
	}
	return entries, true, nil
}

func (r RedisRepository) StoreSchedule(ctx context.Context, groupID string, entries []model.Entry) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, "schedule:"+groupID, data, 4*time.Hour).Err()
}
