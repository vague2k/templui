package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	components "github.com/axzilla/goilerplate/internals/ui/pages"
)

func main() {
	mux := http.NewServeMux()
	isDevelopment := os.Getenv("GO_ENV") == "development"
	SetupAssetsRoutes(mux, isDevelopment)

	mux.Handle("GET /", templ.Handler(components.HeaderShowcase()))

	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", mux)
}

func SetupAssetsRoutes(mux *http.ServeMux, isDevelopment bool) {
	// We need this for Templ to work
	disableCacheInDevMode := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isDevelopment {
				w.Header().Set("Cache-Control", "no-store")
			}
			next.ServeHTTP(w, r)
		})
	}

	// Serve static files from the assets directory
	var fs http.Handler
	if isDevelopment {
		fs = http.FileServer(http.Dir("./assets"))
	} else {
		// In production, use embed.FS (you'll need to set this up)
		// fs = http.FileServer(http.FS(assets.Assets))
		// For now, we'll use the same as development
		fs = http.FileServer(http.Dir("./assets"))
	}

	// Handle requests for assets
	mux.Handle("GET /assets/", disableCacheInDevMode(http.StripPrefix("/assets/", fs)))
}
