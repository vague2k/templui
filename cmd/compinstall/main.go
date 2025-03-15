package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	configFileName = ".templui.json"
	// Version of the tool
	version = "0.1.0"
)

// Config represents the configuration for the component installer
type Config struct {
	ComponentsDir string `json:"componentsDir"`
	ModuleName    string `json:"moduleName"`
}

// Component represents a component that can be installed
type Component struct {
	Name         string
	Description  string
	Files        []string
	Dependencies []string // Names of components this component depends on
}

func main() {
	// Check if we should show the version
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("TemplUI Component Installer v%s\n", version)
		return
	}

	// Check if we should show help
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		showHelp()
		return
	}

	// Check if we should initialize the config
	if len(os.Args) > 1 && os.Args[1] == "init" {
		initConfig()
		return
	}

	// Check if we should add a component
	if len(os.Args) > 1 && os.Args[1] == "add" {
		if len(os.Args) < 3 {
			fmt.Println("Error: No component specified")
			fmt.Println("Usage: templui-install add <component>")
			return
		}

		// Get the component name
		componentName := os.Args[2]

		// Load the config
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			fmt.Println("Run 'templui-install init' to create a config file")
			return
		}

		// Load the components
		components := loadComponents()

		// Find the component
		var component Component
		found := false
		for _, comp := range components {
			if comp.Name == componentName {
				component = comp
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Error: Component '%s' not found\n", componentName)
			fmt.Println("Available components:")
			for _, comp := range components {
				fmt.Printf("  - %s: %s\n", comp.Name, comp.Description)
			}
			return
		}

		// Install the component
		err = installComponent(config.ComponentsDir, component, components, config.ModuleName)
		if err != nil {
			fmt.Printf("Error installing component: %v\n", err)
			return
		}

		// Install utils if needed
		err = installUtils(config.ComponentsDir, config.ModuleName)
		if err != nil {
			fmt.Printf("Error installing utils: %v\n", err)
			return
		}

		fmt.Printf("Component '%s' installed successfully\n", componentName)
		return
	}

	// If no command is specified, show the help
	showHelp()
}

// showHelp shows the help message
func showHelp() {
	fmt.Println("TemplUI Component Installer")
	fmt.Println("Usage:")
	fmt.Println("  templui-install init                - Initialize the config file")
	fmt.Println("  templui-install add <component>     - Add a component to your project")
	fmt.Println("  templui-install -v, --version       - Show version")
	fmt.Println("  templui-install -h, --help          - Show help")

	// Show available components
	components := loadComponents()
	if len(components) > 0 {
		fmt.Println("\nAvailable components:")
		for _, comp := range components {
			fmt.Printf("  - %s: %s\n", comp.Name, comp.Description)
		}
	}
}

// initConfig initializes the config file
func initConfig() {
	// Check if the config file already exists
	if _, err := os.Stat(configFileName); err == nil {
		fmt.Println("Config file already exists")
		return
	}

	// Create a default config
	config := Config{
		ComponentsDir: "./components",
		ModuleName:    detectModuleName(),
	}

	// Ask the user for the components directory
	fmt.Print("Enter the directory where you want to install the components [./components]: ")
	var componentsDir string
	fmt.Scanln(&componentsDir)

	if componentsDir != "" {
		config.ComponentsDir = componentsDir
	}

	// Ask the user for the module name
	fmt.Printf("Enter your Go module name [%s]: ", config.ModuleName)
	var moduleName string
	fmt.Scanln(&moduleName)

	if moduleName != "" {
		config.ModuleName = moduleName
	}

	// Save the config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Error creating config: %v\n", err)
		return
	}

	err = os.WriteFile(configFileName, data, 0644)
	if err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		return
	}

	fmt.Println("Config file created successfully")
	fmt.Printf("Components will be installed to: %s\n", config.ComponentsDir)
	fmt.Printf("Using module name: %s\n", config.ModuleName)
}

// detectModuleName tries to detect the Go module name from go.mod
func detectModuleName() string {
	// Try to read go.mod
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "github.com/username/project"
	}

	// Extract the module name
	re := regexp.MustCompile(`module\s+(.+)`)
	matches := re.FindSubmatch(data)
	if len(matches) < 2 {
		return "github.com/username/project"
	}

	return strings.TrimSpace(string(matches[1]))
}

// loadConfig loads the configuration from the config file
func loadConfig() (Config, error) {
	var config Config

	// Try to find config file in current directory
	configPath := configFileName
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, fmt.Errorf("config file not found")
	}

	// Read and parse config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}

// findTemplUIComponents tries to find the TemplUI components directory
func findTemplUIComponents() string {
	// First, check if we're in a TemplUI project
	// Try to find the components directory in the current directory or parent directories
	currentDir, err := os.Getwd()
	if err == nil {
		// Try the current directory first
		if _, err := os.Stat(filepath.Join(currentDir, "components")); err == nil {
			return filepath.Join(currentDir, "components")
		}

		// Try going up directories until we find the components directory
		for {
			// Check if we're at the root directory
			if currentDir == "/" || currentDir == "." {
				break
			}

			// Check if the components directory exists in this directory
			if _, err := os.Stat(filepath.Join(currentDir, "components")); err == nil {
				return filepath.Join(currentDir, "components")
			}

			// Move up one directory
			parentDir := filepath.Dir(currentDir)
			if parentDir == currentDir {
				// We've reached the root
				break
			}
			currentDir = parentDir
		}
	}

	// If we couldn't find it in the project, try to find it in the global TemplUI installation
	// This assumes TemplUI is installed globally
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		templUIPath := filepath.Join(gopath, "src", "github.com", "axzilla", "templui", "components")
		if _, err := os.Stat(templUIPath); err == nil {
			return templUIPath
		}
	}

	// If we still couldn't find it, return a default path
	return "./components"
}

// findTemplUIUtils tries to find the TemplUI utils directory
func findTemplUIUtils() string {
	// First, check if we're in a TemplUI project
	// Try to find the utils directory in the current directory or parent directories
	currentDir, err := os.Getwd()
	if err == nil {
		// Try the current directory first
		if _, err := os.Stat(filepath.Join(currentDir, "utils")); err == nil {
			return filepath.Join(currentDir, "utils")
		}

		// Try going up directories until we find the utils directory
		for {
			// Check if we're at the root directory
			if currentDir == "/" || currentDir == "." {
				break
			}

			// Check if the utils directory exists in this directory
			if _, err := os.Stat(filepath.Join(currentDir, "utils")); err == nil {
				return filepath.Join(currentDir, "utils")
			}

			// Move up one directory
			parentDir := filepath.Dir(currentDir)
			if parentDir == currentDir {
				// We've reached the root
				break
			}
			currentDir = parentDir
		}
	}

	// If we couldn't find it in the project, try to find it in the global TemplUI installation
	// This assumes TemplUI is installed globally
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		templUIPath := filepath.Join(gopath, "src", "github.com", "axzilla", "templui", "utils")
		if _, err := os.Stat(templUIPath); err == nil {
			return templUIPath
		}
	}

	// If we still couldn't find it, return a default path
	return "./utils"
}

// loadComponents loads the available components
func loadComponents() []Component {
	// Get the source directory (where the components are stored)
	srcDir := findTemplUIComponents()

	// Read the directory
	files, err := os.ReadDir(srcDir)
	if err != nil {
		fmt.Printf("Error reading components directory: %v\n", err)
		// Return some default components if we can't read the directory
		return getDefaultComponents()
	}

	// Map to track components by name
	componentMap := make(map[string]*Component)

	// Process each file
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()

		// Extract component name from file name
		// Handle different file patterns:
		// 1. component.templ
		// 2. component_templ.go
		// 3. component_templ.txt

		var componentName string

		if strings.HasSuffix(name, ".templ") {
			componentName = strings.TrimSuffix(name, ".templ")
		} else if strings.HasSuffix(name, "_templ.go") {
			componentName = strings.TrimSuffix(name, "_templ.go")
		} else if strings.HasSuffix(name, "_templ.txt") {
			componentName = strings.TrimSuffix(name, "_templ.txt")
		} else {
			// Skip files that don't match our patterns
			continue
		}

		// Get or create component
		comp, exists := componentMap[componentName]
		if !exists {
			comp = &Component{
				Name:         componentName,
				Description:  fmt.Sprintf("The %s component", componentName),
				Files:        []string{},
				Dependencies: []string{},
			}
			componentMap[componentName] = comp
		}

		// Add file to component
		comp.Files = append(comp.Files, name)

		// If this is a .templ file, try to extract dependencies and description
		if strings.HasSuffix(name, ".templ") {
			// Read the file to extract dependencies
			filePath := filepath.Join(srcDir, name)
			content, err := os.ReadFile(filePath)
			if err == nil {
				// Extract dependencies
				deps := extractDependencies(string(content), componentMap)
				comp.Dependencies = append(comp.Dependencies, deps...)

				// Extract description
				description := extractDescription(string(content))
				if description != "" {
					comp.Description = description
				}
			}
		}
	}

	// Convert map to slice
	var components []Component
	for _, comp := range componentMap {
		// Remove duplicate dependencies
		comp.Dependencies = removeDuplicates(comp.Dependencies)
		components = append(components, *comp)
	}

	// If no components found, return defaults
	if len(components) == 0 {
		return getDefaultComponents()
	}

	return components
}

// extractDependencies extracts dependencies from component content
func extractDependencies(content string, componentMap map[string]*Component) []string {
	var deps []string

	// This is a simple implementation that looks for component names in the content
	// A more sophisticated implementation would parse the template and look for actual component usage
	for name := range componentMap {
		// Skip self-references
		if strings.Contains(content, name+"(") || strings.Contains(content, name+".") {
			deps = append(deps, name)
		}
	}

	return deps
}

// extractDescription extracts a description from component content
func extractDescription(content string) string {
	// Look for a comment that describes the component
	// This is a simple implementation that looks for specific comment patterns

	// Try to find a comment block at the beginning of the file
	lines := strings.Split(content, "\n")
	var commentLines []string
	inComment := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for comment start
		if strings.HasPrefix(trimmed, "//") {
			// Add to comment lines
			commentLine := strings.TrimPrefix(trimmed, "//")
			commentLine = strings.TrimSpace(commentLine)

			// Skip empty comments
			if commentLine == "" {
				continue
			}

			// Skip comments that are likely not descriptions
			if strings.HasPrefix(commentLine, "TODO") ||
				strings.HasPrefix(commentLine, "FIXME") ||
				strings.HasPrefix(commentLine, "NOTE") {
				continue
			}

			commentLines = append(commentLines, commentLine)
			inComment = true
		} else if inComment {
			// End of comment block
			break
		}
	}

	// If we found a comment block, use it as the description
	if len(commentLines) > 0 {
		// Use the first non-empty line as the description
		for _, line := range commentLines {
			if line != "" {
				return line
			}
		}
	}

	return ""
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

// getDefaultComponents returns a default set of components
func getDefaultComponents() []Component {
	return []Component{
		{
			Name:         "button",
			Description:  "A button component with various styles",
			Files:        []string{"button.templ", "button_templ.go", "button_templ.txt"},
			Dependencies: []string{},
		},
		{
			Name:         "card",
			Description:  "A card component for displaying content",
			Files:        []string{"card.templ", "card_templ.go", "card_templ.txt"},
			Dependencies: []string{},
		},
		{
			Name:         "input",
			Description:  "An input component for forms",
			Files:        []string{"input.templ", "input_templ.go", "input_templ.txt"},
			Dependencies: []string{},
		},
		{
			Name:         "modal",
			Description:  "A modal dialog component",
			Files:        []string{"modal.templ", "modal_templ.go", "modal_templ.txt"},
			Dependencies: []string{},
		},
		{
			Name:         "toast",
			Description:  "A toast notification component",
			Files:        []string{"toast.templ", "toast_templ.go", "toast_templ.txt"},
			Dependencies: []string{},
		},
	}
}

// installComponent installs a component and its dependencies
func installComponent(destDir string, component Component, allComponents []Component, moduleName string) error {
	// Create destination directory if it doesn't exist
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	// Get the source directory (where the components are stored)
	srcDir := findTemplUIComponents()

	// Create a map for quick component lookup
	compMap := make(map[string]Component)
	for _, comp := range allComponents {
		compMap[comp.Name] = comp
	}

	// Install dependencies first
	for _, depName := range component.Dependencies {
		if dep, ok := compMap[depName]; ok {
			// Skip circular dependencies
			if dep.Name == component.Name {
				continue
			}

			// Install the dependency
			err := installComponent(destDir, dep, allComponents, moduleName)
			if err != nil {
				return err
			}

			fmt.Printf("Installed dependency: %s\n", dep.Name)
		}
	}

	// Copy each file
	for _, file := range component.Files {
		srcPath := filepath.Join(srcDir, file)
		destPath := filepath.Join(destDir, file)

		// Read source file
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return err
		}

		// Replace imports if it's a .templ file
		if strings.HasSuffix(file, ".templ") {
			data = replaceImports(data, moduleName)
		}

		// Write to destination
		err = os.WriteFile(destPath, data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// installUtils installs the utils directory
func installUtils(componentsDir string, moduleName string) error {
	// Create the utils directory
	utilsDestDir := filepath.Join(filepath.Dir(componentsDir), "utils")
	err := os.MkdirAll(utilsDestDir, 0755)
	if err != nil {
		return err
	}

	// Get the source directory (where the utils are stored)
	srcDir := findTemplUIUtils()

	// Copy the utils
	files, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("error reading utils directory: %v", err)
	}

	// Copy each file
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		srcPath := filepath.Join(srcDir, file.Name())
		destPath := filepath.Join(utilsDestDir, file.Name())

		// Read source file
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return err
		}

		// Replace imports if it's a .go file
		if strings.HasSuffix(file.Name(), ".go") {
			content := string(data)
			content = regexp.MustCompile(`"github\.com/axzilla/templui/([^"]+)"`).ReplaceAllString(content, fmt.Sprintf(`"%s/$1"`, moduleName))
			data = []byte(content)
		}

		// Write to destination
		err = os.WriteFile(destPath, data, 0644)
		if err != nil {
			return err
		}
	}

	fmt.Println("Installed utils")
	return nil
}

// replaceImports replaces the imports in a .templ file
func replaceImports(data []byte, moduleName string) []byte {
	content := string(data)

	// Replace imports
	// Example: "github.com/axzilla/templui/utils" -> "<moduleName>/utils"
	content = regexp.MustCompile(`"github\.com/axzilla/templui/([^"]+)"`).ReplaceAllString(content, fmt.Sprintf(`"%s/$1"`, moduleName))

	return []byte(content)
}
