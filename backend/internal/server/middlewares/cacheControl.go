package middlewares

import (
	"net/http"
)

// PrivateCacheControl adds the `Cache-Control` header to HTTP requests.
// Use in private routes (that require auth) and that should not be cached
// under any circumstances.
func (m *Middlewares) PrivateCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, private")
		next.ServeHTTP(w, r)
	})
}

// PublicCacheCacheControl adds the `Cache-Control` header to HTTP requests
// with the `no-cache` option.
func (m *Middlewares) PublicCacheCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, public")
		next.ServeHTTP(w, r)
	})
}
