package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	iconDir        = "./lucide/icons" // Path to the Lucide SVG files
	outputFile     = "./icons/icon_defs.go"
	iconContentDir = "./icons/content" // Directory for individual icon contents
)

func main() {
	// Read all files from the icon directory
	files, err := os.ReadDir(iconDir)
	if err != nil {
		panic(err)
	}

	// Initialize slice for icon definitions
	var iconDefs []string
	iconDefs = append(iconDefs, "package icons\n")
	iconDefs = append(iconDefs, "// This file is auto generated\n")

	// Create the content directory if it doesn't exist
	err = os.MkdirAll(iconContentDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Process each SVG file
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".svg" {
			name := strings.TrimSuffix(file.Name(), ".svg")
			funcName := toPascalCase(name)

			// Add icon definition
			iconDefs = append(iconDefs, fmt.Sprintf("var %s = Icon(%q)\n", funcName, name))

			// Save icon content to a separate file
			content, err := os.ReadFile(filepath.Join(iconDir, file.Name()))
			if err != nil {
				panic(err)
			}

			err = os.WriteFile(filepath.Join(iconContentDir, name+".svg"), content, 0644)
			if err != nil {
				panic(err)
			}
		}
	}

	// Write all icon definitions to the output file
	err = os.WriteFile(outputFile, []byte(strings.Join(iconDefs, "")), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Icon definitions and contents generated successfully!")
}

// toPascalCase converts a kebab-case string to PascalCase
func toPascalCase(s string) string {
	words := strings.Split(s, "-")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}
