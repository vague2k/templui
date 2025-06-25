package middleware

import (
	"context"
	"net/http"
	"os"
	"os/exec"

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

func LatestVersion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("version.txt")
		if err != nil {
			// If version.txt doesn't exist, fall back to git commit hash
			commit, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
			if err != nil {
				http.Error(w, "Unable to determine version", http.StatusInternalServerError)
				return
			}
			data = commit
		}
		ctx := context.WithValue(r.Context(), ctxkeys.Version, string(data))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
