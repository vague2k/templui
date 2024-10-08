package icons

import (
	"fmt"
	"strings"
)

// Get returns the SVG content for a given icon name
func Get(name string) (string, bool) {
	content, ok := IconContent[name]
	return content, ok
}

// List returns a slice of all available icon names
func List() []string {
	var names []string
	for name := range IconContent {
		names = append(names, name)
	}
	return names
}

// GetClasses generates the class string for the icon
func GetClasses(name, class string) string {
	classes := []string{"lucide", "lucide-" + strings.ToLower(name)}
	if class != "" {
		classes = append(classes, class)
	}
	return strings.Join(classes, " ")
}

// In pkg/icons/icons.go
func GenerateSVG(name, size, fill, stroke, class string) (string, bool) {
	content, ok := Get(name)
	if !ok {
		return "", false
	}
	svg := fmt.Sprintf(`<svg 
		xmlns="http://www.w3.org/2000/svg"
		width="%s"
		height="%s"
		viewBox="0 0 24 24"
		fill="%s"
		stroke="%s"
		stroke-width="2"
		stroke-linecap="round"
		stroke-linejoin="round"
		class="%s">%s</svg>`,
		size,
		size,
		fill,
		stroke,
		GetClasses(name, class),
		content)
	return svg, true
}
