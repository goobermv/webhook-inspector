package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/goobermv/webhook-inspector/pkg/httputil"
)

func Auth(apiKey string) func(http.Handler) http.Handler {
	expected := []byte(apiKey)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Header.Get("Authorization")
			got = strings.TrimPrefix(got, "Bearer ")
			if got == "" {
				got = r.Header.Get("X-API-Key")
			}
			if subtle.ConstantTimeCompare([]byte(got), expected) != 1 {
				httputil.Error(w, http.StatusUnauthorized, "missing or invalid api key")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
