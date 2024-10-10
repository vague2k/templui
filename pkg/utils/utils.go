package utils

import (
	"fmt"

	"github.com/a-h/templ"
)

// CSS returns a link tag for the Goilerplate CSS
// If no branch is specified, it defaults to 'main'
func CSS(branch string) templ.Component {
	if branch == "" {
		branch = "main"
	}

	cssURL := fmt.Sprintf("https://cdn.jsdelivr.net/gh/axzilla/goilerplate@%s/pgk/styles/goilerplate.css", branch)
	return templ.Raw(fmt.Sprintf(`<link rel="stylesheet" href="%s">`, cssURL))
}

// Alpine returns script tags for Alpine.js
func AlpineJS() templ.Component {
	return templ.Raw(`
        <script defer src="https://cdn.jsdelivr.net/npm/@alpinejs/focus@3.x.x/dist/cdn.min.js"></script>
        <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
    `)
}
