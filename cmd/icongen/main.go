package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Icon struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func main() {
	svgDir := "./cmd/icongen/icons"
	icons := []Icon{}

	err := filepath.Walk(svgDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".svg" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			name := strings.TrimSuffix(filepath.Base(path), ".svg")
			svgContent := extractSVGContent(string(content))
			if svgContent != "" {
				icons = append(icons, Icon{Name: name, Content: svgContent})
			} else {
				fmt.Printf("Warning: Empty content for icon %s\n", name)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", svgDir, err)
		return
	}

	generateGoCode(icons)
}
func extractSVGContent(svgContent string) string {
	start := strings.Index(svgContent, "<svg")
	if start == -1 {
		return ""
	}
	contentStart := strings.Index(svgContent[start:], ">") + start + 1
	end := strings.LastIndex(svgContent, "</svg>")
	if end == -1 || contentStart >= end {
		return ""
	}
	return strings.TrimSpace(svgContent[contentStart:end])
}

func generateGoCode(icons []Icon) {
	file, err := os.Create("./internals/ui/components/icon_contents.go")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	fmt.Fprintln(file, "package components")
	fmt.Fprintln(file, "\nvar iconContents = map[string]string{")
	for _, icon := range icons {
		if icon.Content != "" {
			fmt.Fprintf(file, "\t\"%s\": `%s`,\n", icon.Name, icon.Content)
		}
	}
	fmt.Fprintln(file, "}")
}
