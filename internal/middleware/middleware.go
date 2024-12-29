package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/axzilla/templui/internal/config"
)

func WithPreviewCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isPreview := strings.HasPrefix(r.Host, "preview.")
		ctx := context.WithValue(r.Context(), config.PreviewContextKey, isPreview)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0, must-revalidate, no-cache, no-store, private")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		next.ServeHTTP(w, r)
	})
}
