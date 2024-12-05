package utils

import (
	"fmt"

	"github.com/a-h/templ"
	"math/rand"

	twmerge "github.com/Oudwins/tailwind-merge-go"
)

// TwMerge combines Tailwind classes and handles conflicts.
// Later classes override earlier ones with the same base.
// Example: "bg-red-500 hover:bg-blue-500", "bg-green-500" → "hover:bg-blue-500 bg-green-500"
func TwMerge(classes ...string) string {
	return twmerge.Merge(classes...)
}

// TwIf returns a class if a condition is true, otherwise an empty string
// Example: "bg-red-500", true → "bg-red-500", false → ""
func TwIf(class string, condition bool) string {
	result := templ.KV(class, condition)
	if result.Value == true {
		return result.Key
	}
	return ""
}

// RandomID returns a random ID string
// Example: "id-123456"
func RandomID() string {
	return fmt.Sprintf("id-%d", rand.Intn(1000000))
}

// CSS returns a link tag for the Goilerplate CSS
// If no branch is specified, it defaults to 'main'
func CSS(branch string) templ.Component {
	if branch == "" {
		branch = "main"
	}

	cssURL := fmt.Sprintf("https://cdn.jsdelivr.net/gh/axzilla/goilerplate@%s/pkg/styles/goilerplate.css", branch)
	return templ.Raw(fmt.Sprintf(`<link rel="stylesheet" href="%s">`, cssURL))
}

// Alpine returns script tags for Alpine.js
func AlpineJS() templ.Component {
	return templ.Raw(`
        <script defer src="https://cdn.jsdelivr.net/npm/@alpinejs/focus@3.x.x/dist/cdn.min.js"></script>
        <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
    `)
}
