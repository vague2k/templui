package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/axzilla/templui/assets"
	"github.com/axzilla/templui/components"
	"github.com/axzilla/templui/internal/config"
	"github.com/axzilla/templui/internal/middleware"
	"github.com/axzilla/templui/internal/ui/pages"
	mw "github.com/axzilla/templui/middleware"
	"github.com/axzilla/templui/static"
)

func toastDemoHandler(w http.ResponseWriter, r *http.Request) {
	duration, err := strconv.Atoi(r.FormValue("duration"))
	if err != nil {
		duration = 0
	}

	toastProps := components.ToastProps{
		Title:         r.FormValue("title"),
		Description:   r.FormValue("description"),
		Variant:       components.ToastVariant(r.FormValue("type")),
		Position:      components.ToastPosition(r.FormValue("position")),
		Duration:      duration,
		Dismissible:   r.FormValue("dismissible") == "on",
		ShowIndicator: r.FormValue("indicator") == "on",
		Icon:          r.FormValue("icon") == "on",
	}

	components.Toast(toastProps).Render(r.Context(), w)
}

func buttonHtmxLoadingHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	mux := http.NewServeMux()
	config.LoadConfig()
	SetupAssetsRoutes(mux)

	cspConfig := mw.CSPConfig{
		ScriptSrc: []string{
			"cdn.jsdelivr.net",     // HTMX
			"cdnjs.cloudflare.com", // highlight.js
		},
	}

	wrappedMux := middleware.WithURLPathValue(
		middleware.CacheControlMiddleware(
			mw.WithCSP(cspConfig)(mux),
		),
	)

	// SEO
	// Sitemap.xml Handler mit eingebetteter Datei
	mux.HandleFunc("GET /sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		content, err := static.Files.ReadFile("sitemap.xml")
		if err != nil {
			http.Error(w, "Sitemap not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Write(content)
	})

	// Robots.txt Handler mit eingebetteter Datei
	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		content, err := static.Files.ReadFile("robots.txt")
		if err != nil {
			http.Error(w, "Robots.txt not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write(content)
	})
	//
	mux.Handle("GET /", templ.Handler(pages.Landing()))
	mux.Handle("GET /docs/getting-started", http.RedirectHandler("/docs/introduction", http.StatusSeeOther))
	mux.Handle("GET /docs/components", http.RedirectHandler("/docs/components/accordion", http.StatusSeeOther))
	mux.Handle("GET /docs/introduction", templ.Handler(pages.Introduction()))
	mux.Handle("GET /docs/how-to-use", templ.Handler(pages.HowToUse()))
	mux.Handle("GET /docs/themes", templ.Handler(pages.Themes()))
	// Components
	mux.Handle("GET /docs/components/accordion", templ.Handler(pages.Accordion()))
	mux.Handle("GET /docs/components/alert", templ.Handler(pages.Alert()))
	mux.Handle("GET /docs/components/aspect-ratio", templ.Handler(pages.AspectRatio()))
	mux.Handle("GET /docs/components/avatar", templ.Handler(pages.Avatar()))
	mux.Handle("GET /docs/components/badge", templ.Handler(pages.Badge()))
	mux.Handle("GET /docs/components/breadcrumb", templ.Handler(pages.Breadcrumb()))
	mux.Handle("GET /docs/components/button", templ.Handler(pages.Button()))
	mux.Handle("GET /docs/components/card", templ.Handler(pages.Card()))
	mux.Handle("GET /docs/components/carousel", templ.Handler(pages.Carousel()))
	mux.Handle("GET /docs/components/charts", templ.Handler(pages.Chart()))
	mux.Handle("GET /docs/components/code", templ.Handler(pages.Code()))
	mux.Handle("GET /docs/components/checkbox", templ.Handler(pages.Checkbox()))
	mux.Handle("GET /docs/components/checkbox-card", templ.Handler(pages.CheckboxCard()))
	mux.Handle("GET /docs/components/date-picker", templ.Handler(pages.DatePicker()))
	mux.Handle("GET /docs/components/drawer", templ.Handler(pages.Drawer()))
	mux.Handle("GET /docs/components/dropdown-menu", templ.Handler(pages.DropdownMenu()))
	mux.Handle("GET /docs/components/form", templ.Handler(pages.Form()))
	mux.Handle("GET /docs/components/icon", templ.Handler(pages.Icon()))
	mux.Handle("GET /docs/components/input", templ.Handler(pages.Input()))
	mux.Handle("GET /docs/components/input-otp", templ.Handler(pages.InputOtp()))
	mux.Handle("GET /docs/components/label", templ.Handler(pages.Label()))
	mux.Handle("GET /docs/components/modal", templ.Handler(pages.Modal()))
	mux.Handle("GET /docs/components/pagination", templ.Handler(pages.Pagination()))
	mux.Handle("GET /docs/components/progress", templ.Handler(pages.Progress()))
	mux.Handle("GET /docs/components/radio", templ.Handler(pages.Radio()))
	mux.Handle("GET /docs/components/radio-card", templ.Handler(pages.RadioCard()))
	mux.Handle("GET /docs/components/rating", templ.Handler(pages.Rating()))
	mux.Handle("GET /docs/components/select", templ.Handler(pages.Select()))
	mux.Handle("GET /docs/components/separator", templ.Handler(pages.Separator()))
	mux.Handle("GET /docs/components/skeleton", templ.Handler(pages.Skeleton()))
	mux.Handle("GET /docs/components/slider", templ.Handler(pages.Slider()))
	mux.Handle("GET /docs/components/spinner", templ.Handler(pages.Spinner()))
	mux.Handle("GET /docs/components/table", templ.Handler(pages.Table()))
	mux.Handle("GET /docs/components/tabs", templ.Handler(pages.Tabs()))
	mux.Handle("GET /docs/components/textarea", templ.Handler(pages.Textarea()))
	// mux.Handle("GET /docs/components/time-picker", templ.Handler(pages.TimePicker()))
	mux.Handle("GET /docs/components/toast", templ.Handler(pages.Toast()))
	mux.Handle("GET /docs/components/toggle", templ.Handler(pages.Toggle()))
	mux.Handle("GET /docs/components/tooltip", templ.Handler(pages.Tooltip()))
	// Showcase API
	mux.Handle("POST /docs/toast/demo", http.HandlerFunc(toastDemoHandler))
	mux.Handle("POST /docs/button/htmx-loading", http.HandlerFunc(buttonHtmxLoadingHandler))

	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", wrappedMux)
}

func SetupAssetsRoutes(mux *http.ServeMux) {
	var isDevelopment = config.AppConfig.GoEnv != "production"

	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
