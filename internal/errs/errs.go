package errs

import (
	"errors"
	"net/http"
)

var (
	ErrFetchFailed = errors.New("failed to fetch data")
	ErrNotFound    = errors.New("not found")
	ErrTooManyReq  = errors.New("too many requests")
)

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
