package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/axzilla/goilerplate/assets"
	"github.com/axzilla/goilerplate/internals/config"
	"github.com/axzilla/goilerplate/internals/middleware"
	"github.com/axzilla/goilerplate/internals/ui/pages"
)

func main() {
	mux := http.NewServeMux()
	config.LoadConfig()
	SetupAssetsRoutes(mux)

	wrappedMux := middleware.WithPreviewCheck(mux)

	mux.Handle("GET /", templ.Handler(pages.Landing()))
	mux.Handle("GET /docs/components", http.RedirectHandler("/docs/components/accordion", http.StatusSeeOther))
	mux.Handle("GET /docs/getting-started", http.RedirectHandler("/docs/introduction", http.StatusSeeOther))
	mux.Handle("GET /docs/introduction", templ.Handler(pages.Introduction()))
	mux.Handle("GET /docs/how-to-use", templ.Handler(pages.HowToUse()))
	mux.Handle("GET /docs/themes", templ.Handler(pages.Themes()))
	// Components
	mux.Handle("GET /docs/components/accordion", templ.Handler(pages.Accordion()))
	mux.Handle("GET /docs/components/alert", templ.Handler(pages.Alert()))
	mux.Handle("GET /docs/components/avatar", templ.Handler(pages.Avatar()))
	mux.Handle("GET /docs/components/button", templ.Handler(pages.Button()))
	mux.Handle("GET /docs/components/card", templ.Handler(pages.Card()))
	mux.Handle("GET /docs/components/checkbox", templ.Handler(pages.Checkbox()))
	mux.Handle("GET /docs/components/datepicker", templ.Handler(pages.Datepicker()))
	mux.Handle("GET /docs/components/dropdown-menu", templ.Handler(pages.DropdownMenu()))
	mux.Handle("GET /docs/components/icon", templ.Handler(pages.Icon()))
	mux.Handle("GET /docs/components/input", templ.Handler(pages.Input()))
	mux.Handle("GET /docs/components/modal", templ.Handler(pages.Modal()))
	mux.Handle("GET /docs/components/radio", templ.Handler(pages.Radio()))
	mux.Handle("GET /docs/components/select", templ.Handler(pages.Select()))
	mux.Handle("GET /docs/components/sheet", templ.Handler(pages.Sheet()))
	mux.Handle("GET /docs/components/slider", templ.Handler(pages.Slider()))
	mux.Handle("GET /docs/components/tabs", templ.Handler(pages.Tabs()))
	mux.Handle("GET /docs/components/textarea", templ.Handler(pages.Textarea()))
	mux.Handle("GET /docs/components/toggle", templ.Handler(pages.Toggle()))

	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", wrappedMux)
}

func SetupAssetsRoutes(mux *http.ServeMux) {
	var isDevelopment = config.AppConfig.GoEnv != "production"

	mimeTypes := map[string]string{
		".css":   "text/css; charset=utf-8",
		".js":    "application/javascript; charset=utf-8",
		".svg":   "image/svg+xml",
		".html":  "text/html; charset=utf-8",
		".jpg":   "image/jpeg",
		".jpeg":  "image/jpeg",
		".png":   "image/png",
		".gif":   "image/gif",
		".woff":  "font/woff",
		".woff2": "font/woff2",
		".ttf":   "font/ttf",
		".ico":   "image/x-icon",
	}

	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := filepath.Ext(r.URL.Path)

		if mimeType, ok := mimeTypes[ext]; ok {
			w.Header().Set("Content-Type", mimeType)
		}

		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
		}

		var fs http.Handler
		if isDevelopment {
			fs = http.FileServer(http.Dir("./assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}

		fs.ServeHTTP(w, r)
	})

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assetHandler))
}
