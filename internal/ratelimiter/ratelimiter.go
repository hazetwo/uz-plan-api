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

const (
	violationThreshold = 10
	violationWindow    = time.Minute
	banDuration        = 5 * time.Minute
)

type visitor struct {
	limiter     *rate.Limiter
	lastSeen    time.Time
	violations  int
	windowStart time.Time
	bannedUntil time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     rate.Limit
	bucket   int
}

func New(r rate.Limit, b int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		bucket:   b,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[ip]
	if !ok {
		v = &visitor{
			limiter:     rate.NewLimiter(rl.rate, rl.bucket),
			windowStart: time.Now(),
		}
		rl.visitors[ip] = v
	}
	v.lastSeen = time.Now()

	if time.Now().Before(v.bannedUntil) {
		return false
	}

	if time.Since(v.windowStart) > violationWindow {
		v.violations = 0
		v.windowStart = time.Now()
	}

	if !v.limiter.Allow() {
		v.violations++
		if v.violations >= violationThreshold {
			v.bannedUntil = time.Now().Add(banDuration)
			v.violations = 0
		}
		return false
	}

	return true
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func realIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
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
