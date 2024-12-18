package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/axzilla/templui/internals/utils"
)

type CSPConfig struct {
	ScriptSrc []string // Additional script-src Domains
	StyleSrc  []string // Additional style-src Domains
}

func WithCSP(config CSPConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nonce := utils.GenerateNonce()

			// Base script sources
			scriptSrcs := []string{"'self'", fmt.Sprintf("'nonce-%s'", nonce)}
			scriptSrcs = append(scriptSrcs, config.ScriptSrc...)

			// Base style sources
			styleSrcs := []string{"'self'", "'unsafe-inline'"}
			styleSrcs = append(styleSrcs, config.StyleSrc...)

			csp := fmt.Sprintf(
				"script-src %s; style-src %s;",
				strings.Join(scriptSrcs, " "),
				strings.Join(styleSrcs, " "),
			)

			w.Header().Set("Content-Security-Policy", csp)
			ctx := templ.WithNonce(r.Context(), nonce)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
