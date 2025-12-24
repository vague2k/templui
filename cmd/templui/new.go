package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/templui/templui/internal/templates"
)

// TemplateConfig defines the structure of template.json
type TemplateConfig struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Components    []string `json:"components"`
	DefaultConfig struct {
		ComponentsDir string `json:"componentsDir"`
		UtilsDir      string `json:"utilsDir"`
		JSDir         string `json:"jsDir"`
		JSPublicPath  string `json:"jsPublicPath"`
	} `json:"defaultConfig"`
}

// TemplateData holds data for Go template processing
type TemplateData struct {
	ModuleName string
}

// runNew handles the 'new' command logic.
func runNew(args []string, commandArg string, force bool, moduleFlag string) {
	targetRef := getDefaultRef()

	// Parse optional @ref from the command argument.
	baseCommand := commandArg
	if strings.Contains(commandArg, "@") {
		parts := strings.SplitN(commandArg, "@", 2)
		if len(parts) == 2 && parts[0] == "new" && parts[1] != "" {
			targetRef = parts[1]
			baseCommand = "new"
			fmt.Printf("Using specified ref: %s\n", targetRef)
		} else {
			fmt.Printf("Error: Invalid format '%s'. Use 'new' or 'new@<ref>'.\n", commandArg)
			return
		}
	} else if commandArg != "new" {
		fmt.Printf("Error: Unknown command '%s'. Did you mean 'new'?\n", commandArg)
		showHelp(nil, getDefaultRef())
		return
	}
	_ = baseCommand // Avoid unused variable warning

	// Ensure project name is provided
	if len(args) < 2 {
		fmt.Println("Error: No project name specified.")
		fmt.Println("Usage: templui new <project-name> [--module <module-name>]")
		return
	}

	projectName := args[1]

	// Check if project directory already exists
	if _, err := os.Stat(projectName); err == nil {
		if !force {
			fmt.Printf("Error: Directory '%s' already exists. Use -f to overwrite.\n", projectName)
			return
		}
		fmt.Printf("Warning: Directory '%s' exists. Overwriting...\n", projectName)
	}

	// Determine module name
	moduleName := moduleFlag
	if moduleName == "" {
		// Prompt for module name
		reader := bufio.NewReader(os.Stdin)
		defaultModule := "github.com/username/" + projectName
		fmt.Printf("? Go module name [%s]: ", defaultModule)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			moduleName = defaultModule
		} else {
			moduleName = input
		}
	}

	fmt.Println()
	fmt.Println("🚀 Creating new templUI project...")
	fmt.Printf("   Project: %s\n", projectName)
	fmt.Printf("   Module: %s\n", moduleName)
	fmt.Printf("   Version: %s\n", targetRef)
	fmt.Println()

	// Load template config
	templateConfig, err := loadTemplateConfig()
	if err != nil {
		fmt.Printf("Error loading template config: %v\n", err)
		return
	}

	// Create project directory
	err = os.MkdirAll(projectName, 0755)
	if err != nil {
		fmt.Printf("Error creating project directory: %v\n", err)
		return
	}

	// Process and copy template files
	templateData := TemplateData{
		ModuleName: moduleName,
	}

	err = copyTemplateFiles(projectName, templateData)
	if err != nil {
		fmt.Printf("Error copying template files: %v\n", err)
		return
	}
	fmt.Println("✅ Created project structure")

	// Create .templui.json config
	config := Config{
		ComponentsDir: templateConfig.DefaultConfig.ComponentsDir,
		UtilsDir:      templateConfig.DefaultConfig.UtilsDir,
		ModuleName:    moduleName,
		JSDir:         templateConfig.DefaultConfig.JSDir,
		JSPublicPath:  templateConfig.DefaultConfig.JSPublicPath,
	}

	configPath := filepath.Join(projectName, configFileName)
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Error creating config: %v\n", err)
		return
	}
	err = os.WriteFile(configPath, configData, 0644)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}
	fmt.Println("✅ Created .templui.json")

	// Change to project directory and install components
	originalDir, _ := os.Getwd()
	err = os.Chdir(projectName)
	if err != nil {
		fmt.Printf("Error changing to project directory: %v\n", err)
		return
	}
	defer os.Chdir(originalDir)

	// Install utils first
	fmt.Println("\n📦 Installing utils...")
	registry, err := fetchRegistry(targetRef)
	if err != nil {
		fmt.Printf("Warning: Could not fetch registry: %v\n", err)
	} else {
		if len(registry.Utils) > 0 {
			allUtilPaths := []string{}
			for _, utilDef := range registry.Utils {
				allUtilPaths = append(allUtilPaths, utilDef.Path)
			}
			err = installUtils(config, allUtilPaths, targetRef, true)
			if err != nil {
				fmt.Printf("Warning: Error installing utils: %v\n", err)
			}
		}
	}

	// Install components
	if len(templateConfig.Components) > 0 {
		fmt.Println("\n📦 Installing components...")

		// Build component map
		componentMap := make(map[string]ComponentDef)
		for _, comp := range registry.Components {
			componentMap[comp.Name] = comp
		}

		installedComponents := make(map[string]bool)
		requiredUtils := make(map[string]bool)

		for _, compName := range templateConfig.Components {
			compDef, exists := componentMap[compName]
			if !exists {
				fmt.Printf("   ⚠️  Component '%s' not found in registry\n", compName)
				continue
			}

			err = installComponent(config, compDef, componentMap, targetRef, installedComponents, requiredUtils, true)
			if err != nil {
				fmt.Printf("   ⚠️  Error installing %s: %v\n", compName, err)
			} else {
				fmt.Printf("   ✅ %s\n", compName)
			}
		}

		// Install any additional required utils from components
		if len(requiredUtils) > 0 {
			utilsToInstall := []string{}
			for utilPath := range requiredUtils {
				utilsToInstall = append(utilsToInstall, utilPath)
			}
			installUtils(config, utilsToInstall, targetRef, true)
		}
	}

	// Initialize go.mod
	fmt.Println("\n📦 Initializing Go modules...")

	// Print success message
	fmt.Println()
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("🎉 Project created successfully!")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go mod tidy")
	fmt.Println("  task dev")
	fmt.Println()
	fmt.Println("Happy coding! 🚀")
}

// loadTemplateConfig loads the template.json configuration
func loadTemplateConfig() (TemplateConfig, error) {
	var config TemplateConfig

	data, err := templates.QuickstartFS.ReadFile("quickstart/template.json")
	if err != nil {
		return config, fmt.Errorf("failed to read template.json: %w", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse template.json: %w", err)
	}

	return config, nil
}

// copyTemplateFiles copies and processes template files to the destination
func copyTemplateFiles(destDir string, data TemplateData) error {
	return fs.WalkDir(templates.QuickstartFS, "quickstart", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory and template.json
		if path == "quickstart" || path == "quickstart/template.json" {
			return nil
		}

		// Calculate destination path
		relPath := strings.TrimPrefix(path, "quickstart/")
		destPath := filepath.Join(destDir, relPath)

		// Handle .tmpl files - strip the .tmpl extension
		if strings.HasSuffix(destPath, ".tmpl") {
			destPath = strings.TrimSuffix(destPath, ".tmpl")
		}

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// Read file content
		content, err := templates.QuickstartFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// Process .tmpl files as Go templates
		if strings.HasSuffix(path, ".tmpl") {
			tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
			if err != nil {
				return fmt.Errorf("failed to parse template %s: %w", path, err)
			}

			var processed strings.Builder
			err = tmpl.Execute(&processed, data)
			if err != nil {
				return fmt.Errorf("failed to execute template %s: %w", path, err)
			}
			content = []byte(processed.String())
		}

		// Ensure parent directory exists
		parentDir := filepath.Dir(destPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", parentDir, err)
		}

		// Write file
		err = os.WriteFile(destPath, content, 0644)
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", destPath, err)
		}

		return nil
	})
}
