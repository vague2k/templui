package middleware

import (
	"context"
	"net/http"

	"github.com/axzilla/templui/internal/ctxkeys"
)

func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0, must-revalidate, no-cache, no-store, private")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// WithURLPathValue adds the current URL's path to the context.
func WithURLPathValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(
			r.Context(),
			ctxkeys.URLPathValue,
			r.URL.Path,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
