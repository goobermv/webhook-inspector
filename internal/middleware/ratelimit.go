package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/goobermv/webhook-inspector/pkg/httputil"
)

type bucket struct {
	count int
	reset time.Time
}

func RateLimit(perMinute int) func(http.Handler) http.Handler {
	var mu sync.Mutex
	buckets := map[string]*bucket{}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			mu.Lock()
			b, ok := buckets[ip]
			if !ok || time.Now().After(b.reset) {
				b = &bucket{reset: time.Now().Add(time.Minute)}
				buckets[ip] = b
			}
			b.count++
			exceeded := b.count > perMinute
			mu.Unlock()
			if exceeded {
				httputil.Error(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
