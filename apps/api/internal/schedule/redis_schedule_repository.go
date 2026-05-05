package schedule

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) (*RedisRepository, *redsync.Redsync) {
	pool := goredis.NewPool(rdb)
	return &RedisRepository{rdb: rdb}, redsync.New(pool)
}

func (r RedisRepository) GetFields(ctx context.Context) (map[string]string, bool, error) {
	fields, err := r.rdb.HGetAll(ctx, "fields").Result()
	if err != nil {
		return nil, false, err
	}
	if len(fields) == 0 {
		return nil, false, nil
	}

	return fields, true, nil
}

func (r RedisRepository) StoreFields(ctx context.Context, fields map[string]string) error {
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, "fields", fields)
	pipe.Expire(ctx, "fields", 24*time.Hour)
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r RedisRepository) GetGroups(ctx context.Context, fieldID string) (map[string]string, bool, error) {
	g, err := r.rdb.HGetAll(ctx, "groups:"+fieldID).Result()
	if err != nil {
		return nil, false, err
	}
	if len(g) == 0 {
		return nil, false, nil
	}

	return g, true, nil
}

func (r RedisRepository) StoreGroups(ctx context.Context, fieldID string, groups map[string]string) error {
	p := r.rdb.Pipeline()
	p.HSet(ctx, "groups:"+fieldID, groups)
	p.Expire(ctx, "groups:"+fieldID, 24*time.Hour)
	if _, err := p.Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r RedisRepository) GetSchedule(ctx context.Context, groupID string) ([]Entry, bool, error) {
	data, err := r.rdb.Get(ctx, "schedule:"+groupID).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, false, err
	}
	return entries, true, nil
}

func (r RedisRepository) StoreSchedule(ctx context.Context, groupID string, entries []Entry) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, "schedule:"+groupID, data, 4*time.Hour).Err()
}
