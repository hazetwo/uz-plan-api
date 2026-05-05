package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context) (*redis.Client, error) {
	u := os.Getenv("REDIS_URL")
	if u == "" {
		u = "redis://localhost:6379"
	}
	opts, err := redis.ParseURL(u)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)

	const maxAttempts = 5
	const retryDelay = 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := rdb.Ping(ctx).Err(); err == nil {
			return rdb, nil
		}
		if attempt == maxAttempts {
			_ = rdb.Close()
			return nil, fmt.Errorf("could not connect to Redis at %s after %d attempts", opts.Addr, maxAttempts)
		}
		slog.Warn("Redis not ready, retrying...", "attempt", attempt, "addr", opts.Addr, "retry_in", retryDelay)
		select {
		case <-ctx.Done():
			_ = rdb.Close()
			return nil, ctx.Err()
		case <-time.After(retryDelay):
		}
	}

	return rdb, nil
}
