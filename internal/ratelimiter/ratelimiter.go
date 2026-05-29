package ratelimiter

import (
	"net/http"
	"time"
	"uz-plan-api/internal/errs"

	"github.com/go-chi/render"
	gorl "github.com/hazefyro/go-ratelimiter"
)

type errorResponse struct {
	Error string `json:"error"`
}

// Options mirrors the previous configuration surface so callers don't need to
// know about the underlying go-ratelimiter library.
type Options struct {
	RateLimit          int
	Bucket             int
	ViolationThreshold int
	ViolationWindow    time.Duration
	BanDuration        time.Duration
}

// RateLimiter is a thin adapter around github.com/hazefyro/go-ratelimiter that
// exposes a chi middleware preserving this API's JSON error contract.
type RateLimiter struct {
	limiter *gorl.RateLimiter
}

func New(options Options) (*RateLimiter, error) {
	limiter, err := gorl.New(&gorl.Options{
		RateLimit: options.RateLimit,
		Bucket:    options.Bucket,
		KeyFunc:   gorl.RealIPKey,
		Banning: &gorl.BanOptions{
			Threshold: options.ViolationThreshold,
			Window:    options.ViolationWindow,
			Duration:  options.BanDuration,
		},
	})
	if err != nil {
		return nil, err
	}
	return &RateLimiter{limiter: limiter}, nil
}

// Stop shuts down the underlying limiter's background cleanup goroutine.
func (rl *RateLimiter) Stop() {
	rl.limiter.Stop()
}

func (rl *RateLimiter) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// RealIPKey reads X-Real-IP (set directly by Railway) with a
		// RemoteAddr fallback.
		if err := rl.limiter.Allow(gorl.RealIPKey(r)); err != nil {
			render.Status(r, http.StatusTooManyRequests)
			render.JSON(w, r, errorResponse{Error: errs.ErrTooManyReq.Error()})
			return
		}
		next.ServeHTTP(w, r)
	})
}
