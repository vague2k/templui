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
