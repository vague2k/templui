package helpers

import (
	"github.com/a-h/templ"
	"github.com/axzilla/goilerplate/internals/config"
)

// CSS returns a link tag for the Goilerplate CSS
func CSS() templ.Component {
	if config.AppConfig.GoEnv == "production" {
		return templ.Raw(`<link rel="stylesheet" href="https://github.com/axzilla/goilerplate/blob/main/pkg/styles/goilerplate.css">`)
	}
	return templ.Raw(`<link rel="stylesheet" href="https://github.com/axzilla/goilerplate/blob/dev/pkg/styles/goilerplate.css">`)

}

// Alpine returns script tags for Alpine.js
func Alpine() templ.Component {
	return templ.Raw(`
        <script defer src="https://cdn.jsdelivr.net/npm/@alpinejs/focus@3.x.x/dist/cdn.min.js"></script>
        <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
    `)
}
