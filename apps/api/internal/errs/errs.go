package errs

import "errors"

var (
	ErrFetchFailed = errors.New("failed to fetch data")
	ErrTooManyReq  = errors.New("too many requests")
)
