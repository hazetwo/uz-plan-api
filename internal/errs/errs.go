package errs

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrFetchFailed = errors.New("failed to fetch data")
	ErrNotFound    = errors.New("not found")
	ErrTooManyReq  = errors.New("too many requests")
)

func FetchFailed(ctx context.Context, err error) error {
	slog.ErrorContext(ctx, "fetch failed", "err", err)
	return ErrFetchFailed
}

func StatusFromErr(err error) int {
	switch {
	case errors.Is(err, ErrFetchFailed):
		return http.StatusInternalServerError
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}
