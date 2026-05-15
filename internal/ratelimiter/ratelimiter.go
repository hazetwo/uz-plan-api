package ratelimiter

import (
	"net"
	"net/http"
	"sync"
	"time"
	"uz-plan-api/internal/errs"

	"github.com/go-chi/render"
	"golang.org/x/time/rate"
)

type errorResponse struct {
	Error string `json:"error"`
}

type visitor struct {
	limiter     *rate.Limiter
	lastSeen    time.Time
	violations  int
	windowStart time.Time
	bannedUntil time.Time
}

type Options struct {
	RateLimit          int
	Bucket             int
	ViolationThreshold int
	ViolationWindow    time.Duration
	BanDuration        time.Duration
}

type RateLimiter struct {
	mu        sync.Mutex
	rateLimit rate.Limit
	visitors  map[string]*visitor
	bucket    int
	options   Options
}

func New(options Options) *RateLimiter {
	rl := &RateLimiter{
		visitors:  make(map[string]*visitor),
		rateLimit: rate.Limit(options.RateLimit),
		options:   options,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) allow(ip string) bool {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[ip]
	if !ok {
		v = &visitor{
			limiter:     rate.NewLimiter(rl.rateLimit, rl.options.Bucket),
			windowStart: now,
		}
		rl.visitors[ip] = v
	}
	v.lastSeen = now

	if now.Before(v.bannedUntil) {
		return false
	}

	if time.Since(v.windowStart) > rl.options.ViolationWindow {
		v.violations = 0
		v.windowStart = now
	}

	if !v.limiter.Allow() {
		v.violations++
		if v.violations >= rl.options.ViolationThreshold {
			v.bannedUntil = now.Add(rl.options.BanDuration)
			v.violations = 0
		}
		return false
	}

	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.options.BanDuration+time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func realIP(r *http.Request) string {
	// Railway sets X-Real-IP directly so we can use it
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (rl *RateLimiter) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(realIP(r)) {
			render.Status(r, http.StatusTooManyRequests)
			render.JSON(w, r, errorResponse{Error: errs.ErrTooManyReq.Error()})
			return
		}
		next.ServeHTTP(w, r)
	})
}
