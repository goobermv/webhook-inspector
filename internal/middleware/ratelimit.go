package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/goobermv/webhook-inspector/pkg/httputil"
)

type bucket struct {
	count int
	reset time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	perMin  int
}

func NewRateLimiter(perMinute int) *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*bucket),
		perMin:  perMinute,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)

		rl.mu.Lock()
		b, ok := rl.buckets[ip]
		now := time.Now()
		if !ok || now.After(b.reset) {
			b = &bucket{reset: now.Add(time.Minute)}
			rl.buckets[ip] = b
		}
		b.count++
		exceeded := b.count > rl.perMin
		rl.mu.Unlock()

		if exceeded {
			w.Header().Set("Retry-After", "60")
			httputil.Error(w, http.StatusTooManyRequests, "too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
