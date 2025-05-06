package main

import (
	"encoding/json"
	"fmt"
	"io"       // Import für I/O (Lesen von HTTP-Responses)
	"net/http" // Import für HTTP-Anfragen
	"os"
	"path/filepath"
	"regexp"
	"strings"
	// "bufio" // Wird für die Überschreiben-Frage benötigt
	// "log"   // Für Fehlerlogging beim Download
)

const (
	configFileName    = ".templui.json"
	manifestPath      = "internal/manifest.json" // Pfad zum Manifest im Repo
	defaultVersionTag = "main"                   // Standardmäßig von 'main' laden (später latest tag)
	// Basis-URL für Raw-Content. Passe <user>/<repo> an dein GitHub Repo an!
	// Beispiel: "https://raw.githubusercontent.com/axzilla/templui/"
	rawContentBaseURL = "https://raw.githubusercontent.com/axzilla/templui/"
)

// Version des Tools (kann über ldflags gesetzt werden)
var version = "0.0.0-dev" // Standardwert

// Config bleibt gleich
type Config struct {
	ComponentsDir string `json:"componentsDir"`
	ModuleName    string `json:"moduleName"`
}

// Manifest structure (entspricht der JSON-Datei)
type Manifest struct {
	Version    string         `json:"version"`
	Components []ComponentDef `json:"components"`
	Utils      []UtilDef      `json:"utils"`
}

type ComponentDef struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Files         []string `json:"files"`         // Pfade im Repo
	Dependencies  []string `json:"dependencies"`  // Namen anderer Komponenten
	RequiredUtils []string `json:"requiredUtils"` // Pfade zu Utils im Repo
}

type UtilDef struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}

// --- Hauptfunktion ---
func main() {
	// Check if we should show the version
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("TemplUI Component Installer v%s\n", version)
		return
	}

	// Check if we should show help
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		// TODO: showHelp muss Manifest laden, um Komponenten aufzulisten
		showHelpPlaceholder() // Platzhalter, da showHelp noch angepasst werden muss
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
			fmt.Println("Error: No component(s) specified.")
			fmt.Println("Usage: templui add <component1> [component2...] | *")
			return
		}

		// Load the local config
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			fmt.Println("Run 'templui init' to create a config file.")
			return
		}

		// --- Download-Logik ---
		versionTag := defaultVersionTag // TODO: Version aus Arg parsen (z.B. button@v0.1.1)

		fmt.Printf("Fetching component manifest for version '%s'...\n", versionTag)
		manifest, err := fetchManifest(versionTag)
		if err != nil {
			fmt.Printf("Error fetching manifest: %v\n", err)
			return
		}
		fmt.Printf("Using components from templui version %s\n", manifest.Version)

		// Create map for easy lookup
		componentMap := make(map[string]ComponentDef)
		for _, comp := range manifest.Components {
			componentMap[comp.Name] = comp
		}

		// Determine components to install
		componentsToInstall := []string{}
		if len(os.Args) == 3 && os.Args[2] == "*" {
			fmt.Println("Installing all components...")
			for _, comp := range manifest.Components {
				componentsToInstall = append(componentsToInstall, comp.Name)
			}
		} else {
			componentsToInstall = os.Args[2:]
		}

		// Track installed components and required utils to avoid duplicates
		installedComponents := make(map[string]bool)
		requiredUtils := make(map[string]bool) // Map von Util-Pfad zu bool

		// Install each requested component (and dependencies)
		for _, componentName := range componentsToInstall {
			// TODO: Parse version from componentName if format is "name@version"
			// For now, use the global versionTag
			compDef, exists := componentMap[componentName]
			if !exists {
				fmt.Printf("Error: Component '%s' not found in manifest for version %s.\n", componentName, versionTag)
				// TODO: List available components from manifest
				continue // or return error
			}

			err = installComponent(config, compDef, componentMap, versionTag, installedComponents, requiredUtils)
			if err != nil {
				fmt.Printf("Error installing component %s: %v\n", componentName, err)
				// Decide whether to continue or stop
			}
		}

		// Install required utils
		if len(requiredUtils) > 0 {
			fmt.Println("Installing required utils...")
			utilsToInstallPaths := []string{}
			for utilPath := range requiredUtils {
				utilsToInstallPaths = append(utilsToInstallPaths, utilPath)
			}
			err = installUtils(config, utilsToInstallPaths, versionTag)
			if err != nil {
				fmt.Printf("Error installing utils: %v\n", err)
			}
		}

		fmt.Println("\nInstallation finished.")
		// TODO: Ggf. Zusammenfassung ausgeben (was wurde installiert/übersprungen)
		return
	}

	// If no command is specified, show the help
	showHelpPlaceholder() // Platzhalter
}

// --- Platzhalter für Hilfe ---
func showHelpPlaceholder() {
	fmt.Println("TemplUI Component Installer (v" + version + ")")
	fmt.Println("Usage:")
	fmt.Println("  templui init                - Initialize the config file (.templui.json)")
	fmt.Println("  templui add <comp> [<comp>..] - Add component(s) to your project")
	fmt.Println("  templui add *               - Add all available components")
	// fmt.Println("  templui add <comp>@<version> - Add a specific version of a component") // Zukünftig
	fmt.Println("  templui -v, --version       - Show installer version")
	fmt.Println("  templui -h, --help          - Show this help message")
	fmt.Println("\nRun 'templui add <component>' after 'init'.")
	fmt.Println("Available components will be listed once the manifest can be fetched.")
}

// --- Konfigurationsfunktionen (unverändert) ---
func initConfig() {
	if _, err := os.Stat(configFileName); err == nil {
		fmt.Println("Config file already exists.")
		// TODO: Fragen, ob überschrieben werden soll?
		return
	}

	defaultConfig := Config{
		ComponentsDir: "./components", // Standardziel
		ModuleName:    detectModuleName(),
	}

	fmt.Printf("Enter the directory for components [%s]: ", defaultConfig.ComponentsDir)
	var componentsDir string
	fmt.Scanln(&componentsDir)
	if componentsDir == "" {
		componentsDir = defaultConfig.ComponentsDir
	}

	fmt.Printf("Enter your Go module name [%s]: ", defaultConfig.ModuleName)
	var moduleName string
	fmt.Scanln(&moduleName)
	if moduleName == "" {
		moduleName = defaultConfig.ModuleName
	}

	config := Config{
		ComponentsDir: componentsDir,
		ModuleName:    moduleName,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Error creating config data: %v\n", err)
		return
	}

	err = os.WriteFile(configFileName, data, 0644)
	if err != nil {
		fmt.Printf("Error saving config file: %v\n", err)
		return
	}

	fmt.Println("Config file created successfully at", configFileName)
	fmt.Printf("Components will be installed to: %s\n", config.ComponentsDir)
	fmt.Printf("Using module name: %s\n", config.ModuleName)
}

func detectModuleName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Println("Warning: Could not read go.mod to detect module name. Using default.")
		return "your/module/path" // Sicherer Standard
	}
	re := regexp.MustCompile(`(?m)^\s*module\s+(\S+)`)
	matches := re.FindSubmatch(data)
	if len(matches) < 2 {
		fmt.Println("Warning: Could not parse module name from go.mod. Using default.")
		return "your/module/path"
	}
	return string(matches[1])
}

func loadConfig() (Config, error) {
	var config Config
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		return config, fmt.Errorf("config file '%s' not found", configFileName)
	}
	data, err := os.ReadFile(configFileName)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error parsing config file: %w", err)
	}
	// Validate config? (e.g., ensure ModuleName is not empty)
	if config.ComponentsDir == "" || config.ModuleName == "" {
		return config, fmt.Errorf("invalid config: ComponentsDir and ModuleName must be set")
	}
	return config, nil
}

// --- Download-Hilfsfunktionen (Platzhalter) ---

// fetchManifest lädt das Manifest für einen gegebenen Git-Tag/Branch
func fetchManifest(versionTag string) (Manifest, error) {
	manifestURL := rawContentBaseURL + versionTag + "/" + manifestPath
	fmt.Println("Downloading manifest from:", manifestURL) // Debugging

	resp, err := http.Get(manifestURL)
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to start download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Versuch, Fehlermeldung zu lesen
		return Manifest{}, fmt.Errorf("failed to download manifest: status code %d, message: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var manifest Manifest
	err = json.Unmarshal(body, &manifest)
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to parse manifest JSON: %w", err)
	}

	return manifest, nil
}

// downloadFile lädt eine einzelne Datei von einer URL
func downloadFile(url string) ([]byte, error) {
	// fmt.Println("Downloading file:", url) // Debugging
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to start download from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file from %s: status code %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", url, err)
	}
	return data, nil
}

// TODO: getLatestTag() - Funktion zum Ermitteln des neuesten Git-Tags (optional)

// --- Installationslogik (Angepasst für Download) ---

// installComponent installiert eine Komponente und ihre Abhängigkeiten
func installComponent(
	config Config,
	comp ComponentDef,
	componentMap map[string]ComponentDef,
	versionTag string,
	installed map[string]bool,
	requiredUtils map[string]bool,
) error {
	// Verhindere Endlosschleifen und doppelte Installation
	if installed[comp.Name] {
		return nil
	}
	installed[comp.Name] = true // Markiere als (versucht) installiert

	fmt.Printf("Processing component: %s\n", comp.Name)

	// 1. Abhängigkeiten zuerst installieren (rekursiv)
	for _, depName := range comp.Dependencies {
		depComp, exists := componentMap[depName]
		if !exists {
			// Sollte nicht passieren, wenn Manifest korrekt ist
			fmt.Printf("Warning: Dependency '%s' for component '%s' not found in manifest.\n", depName, comp.Name)
			continue
		}
		// Installiere Abhängigkeit (gleiche Version wie Hauptkomponente)
		err := installComponent(config, depComp, componentMap, versionTag, installed, requiredUtils)
		if err != nil {
			// Fehler bei Abhängigkeit weitergeben oder nur loggen?
			return fmt.Errorf("failed to install dependency '%s' for '%s': %w", depName, comp.Name, err)
		}
	}

	// 2. Dateien der aktuellen Komponente herunterladen und schreiben
	for _, repoFilePath := range comp.Files {
		fileName := filepath.Base(repoFilePath) // z.B. button.templ
		destDir := config.ComponentsDir         // z.B. ./components
		// Unterverzeichnis für Komponente erstellen? Nein, aktuell flach.
		// Bsp: Wenn repoFilePath="internal/components/button/button.templ", dann ist fileName="button.templ"
		// Das Ziel ist dann config.ComponentsDir/button.templ

		destPath := filepath.Join(destDir, fileName)

		// Stelle sicher, dass das Zielverzeichnis existiert
		err := os.MkdirAll(destDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create destination directory '%s': %w", destDir, err)
		}

		// --- Überschreiben-Logik ---
		if _, err := os.Stat(destPath); err == nil {
			// Datei existiert
			// TODO: Implementiere Logik zum Lesen der Version aus der Datei
			// TODO: Frage den Benutzer, ob überschrieben werden soll (bufio.NewReader)
			fmt.Printf("Warning: File '%s' already exists. Skipping overwrite for now.\n", destPath)
			continue // Überspringe diese Datei
		}

		// --- Download ---
		fileURL := rawContentBaseURL + versionTag + "/" + repoFilePath
		fmt.Printf("  Downloading %s...\n", fileName)
		data, err := downloadFile(fileURL)
		if err != nil {
			return fmt.Errorf("failed to download file '%s' for component '%s': %w", fileName, comp.Name, err)
		}

		// --- Modifikationen ---
		// a) Versionskommentar hinzufügen
		versionComment := fmt.Sprintf("// templui component version: %s installed by templui v%s\n", versionTag, version)
		modifiedData := append([]byte(versionComment), data...)

		// b) Importpfade ersetzen (nur für .templ und .go Dateien relevant?)
		if strings.HasSuffix(fileName, ".templ") || strings.HasSuffix(fileName, ".go") {
			modifiedData = replaceImports(modifiedData, config.ModuleName)
		}

		// --- Schreiben ---
		err = os.WriteFile(destPath, modifiedData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file '%s': %w", destPath, err)
		}
		fmt.Printf("  Installed %s\n", destPath)
	}

	// 3. Benötigte Utils für diese Komponente sammeln
	for _, utilPath := range comp.RequiredUtils {
		requiredUtils[utilPath] = true
	}

	return nil
}

// installUtils installiert die benötigten Util-Dateien
func installUtils(config Config, utilPaths []string, versionTag string) error {
	if len(utilPaths) == 0 {
		return nil
	}

	// Bestimme das Zielverzeichnis für Utils (z.B. ./utils neben ./components)
	utilsDestDir := filepath.Join(filepath.Dir(config.ComponentsDir), "utils")
	fmt.Printf("Installing utils to: %s\n", utilsDestDir)

	err := os.MkdirAll(utilsDestDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create utils directory '%s': %w", utilsDestDir, err)
	}

	for _, repoUtilPath := range utilPaths {
		fileName := filepath.Base(repoUtilPath) // z.B. classname.go
		destPath := filepath.Join(utilsDestDir, fileName)

		// --- Überschreiben-Logik (ähnlich wie bei Komponenten) ---
		if _, err := os.Stat(destPath); err == nil {
			// Datei existiert
			// TODO: Implementiere Logik zum Lesen der Version aus der Datei
			// TODO: Frage den Benutzer, ob überschrieben werden soll
			fmt.Printf("Warning: Util file '%s' already exists. Skipping overwrite for now.\n", destPath)
			continue // Überspringe diese Datei
		}

		// --- Download ---
		fileURL := rawContentBaseURL + versionTag + "/" + repoUtilPath
		fmt.Printf("  Downloading util %s...\n", fileName)
		data, err := downloadFile(fileURL)
		if err != nil {
			return fmt.Errorf("failed to download util '%s': %w", fileName, err)
		}

		// --- Modifikationen ---
		// a) Versionskommentar
		versionComment := fmt.Sprintf("// templui util version: %s installed by templui v%s\n", versionTag, version)
		modifiedData := append([]byte(versionComment), data...)

		// b) Importpfade ersetzen (nur für .go Dateien relevant)
		if strings.HasSuffix(fileName, ".go") {
			// Ersetze interne Pfade (z.B. "github.com/axzilla/templui/internal/utilss")
			// mit dem Modulnamen des Nutzers (z.B. "your/module/path/utils")
			// Wichtig: Der Pfad im replaceImports muss dem internen Pfad entsprechen!
			modifiedData = replaceImports(modifiedData, config.ModuleName)
		}

		// --- Schreiben ---
		err = os.WriteFile(destPath, modifiedData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write util file '%s': %w", destPath, err)
		}
		fmt.Printf("  Installed %s\n", destPath)
	}

	return nil
}

// replaceImports ersetzt die internen TemplUI-Importpfade durch den Modulnamen des Benutzers
func replaceImports(data []byte, userModuleName string) []byte {
	content := string(data)
	// Das Muster muss exakt den Importpfad matchen, der in deinen internen Komponenten/Utils verwendet wird!
	// Annahme: Deine internen Imports sehen aus wie "github.com/axzilla/templui/internal/..."
	// Passe <user>/<repo> an dein Repo an!
	internalImportPattern := `"github.com/axzilla/templui/internal/([^"]+)"`
	targetImportFormat := fmt.Sprintf(`"%s/$1"`, userModuleName) // z.B. "your/module/path/utils" oder "your/module/path/components/icon"

	re := regexp.MustCompile(internalImportPattern)
	newContent := re.ReplaceAllString(content, targetImportFormat)

	// Logge, ob Ersetzungen stattgefunden haben (optional für Debugging)
	// if content != newContent {
	//  fmt.Println("  Replaced import paths.")
	// }

	return []byte(newContent)
}

// --- Alte, nicht mehr verwendete Funktionen ---
/*
func findTemplUIComponents() string { ... }
func findTemplUIUtils() string { ... }
func loadComponents() []Component { ... }
func extractDependencies(content string, componentMap map[string]*Component) []string { ... }
func extractDescription(content string) string { ... }
func removeDuplicates(slice []string) []string { ... }
func getDefaultComponents() []Component { ... }
*/
