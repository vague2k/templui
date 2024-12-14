package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/axzilla/templui/internals/config"
	"github.com/axzilla/templui/internals/utils"
)

func WithPreviewCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isPreview := strings.HasPrefix(r.Host, "preview.")
		ctx := context.WithValue(r.Context(), config.PreviewContextKey, isPreview)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithNonce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := utils.GenerateNonce()
		csp := fmt.Sprintf(
			"script-src 'self' 'nonce-%s' cdn.jsdelivr.net unpkg.com cdnjs.cloudflare.com; "+
				"style-src 'self' 'unsafe-inline' cdnjs.cloudflare.com;", // highlight.js CSS erlauben
			nonce,
		)
		w.Header().Set("Content-Security-Policy", csp)
		ctx := templ.WithNonce(r.Context(), nonce)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
