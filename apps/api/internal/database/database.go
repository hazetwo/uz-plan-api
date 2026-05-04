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
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	const maxAttempts = 5
	const retryDelay = 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := rdb.Ping(ctx).Err(); err == nil {
			return rdb, nil
		}
		if attempt == maxAttempts {
			_ = rdb.Close()
			return nil, fmt.Errorf("could not connect to Redis at %s after %d attempts", addr, maxAttempts)
		}
		slog.Warn("Redis not ready, retrying...", "attempt", attempt, "addr", addr, "retry_in", retryDelay)
		select {
		case <-ctx.Done():
			_ = rdb.Close()
			return nil, ctx.Err()
		case <-time.After(retryDelay):
		}
	}

	return rdb, nil
}
