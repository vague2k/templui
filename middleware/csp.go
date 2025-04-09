package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	util "github.com/axzilla/templui/internal/util"
)

type CSPConfig struct {
	ScriptSrc []string // External script domains allowed
}

func WithCSP(config CSPConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nonce, err := util.GenerateNonce()
			if err != nil {
				log.Printf("failed to generate nonce: %v", err)
				w.Header().Set("Content-Security-Policy", "script-src 'self'")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Combine all script sources
			scriptSources := append(
				[]string{"'self'", fmt.Sprintf("'nonce-%s'", nonce)},
				config.ScriptSrc...)

			csp := fmt.Sprintf("script-src %s", strings.Join(scriptSources, " "))
			w.Header().Set("Content-Security-Policy", csp)

			next.ServeHTTP(w, r.WithContext(templ.WithNonce(r.Context(), nonce)))
		})
	}
}
