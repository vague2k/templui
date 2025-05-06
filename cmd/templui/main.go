package main

import (
	"bufio" // Import für Nutzereingabe
	"encoding/json"
	"flag" // Import für Command-Line Flags
	"fmt"
	"io"       // Import für I/O (Lesen von HTTP-Responses)
	"net/http" // Import für HTTP-Anfragen
	"os"
	"path/filepath"
	"regexp"
	"strings"
	// "log"   // Für Fehlerlogging beim Download
)

const (
	configFileName = ".templui.json"
	manifestPath   = "internal/manifest.json" // Pfad zum Manifest im Repo
	// Standard-Ref (branch/tag/commit) zum Laden, wenn keiner angegeben wird.
	// TODO: Später durch getLatestTag() ersetzen?
	defaultRef = "main"
	// Basis-URL für Raw-Content. Passe <user>/<repo> an dein GitHub Repo an!
	// Beispiel: "https://raw.githubusercontent.com/axzilla/templui/"
	rawContentBaseURL = "https://raw.githubusercontent.com/axzilla/templui/"
)

// Version des Tools (kann über ldflags gesetzt werden)
var version = "0.0.0-dev" // Standardwert

// Regex zum Extrahieren der Version aus dem Kommentar (einmal kompilieren)
var versionRegex = regexp.MustCompile(`(?m)^\s*//\s*templui\s+(?:component|util)\s+.*\s+-\s+version:\s+(\S+)`)

// --- Flags ---
var forceOverwrite = flag.Bool("force", false, "Force overwrite existing files without asking")

// Config um UtilsDir erweitern
type Config struct {
	ComponentsDir string `json:"componentsDir"`
	UtilsDir      string `json:"utilsDir"` // Hinzugefügt
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
	// Alias -f für --force hinzufügen
	flag.BoolVar(forceOverwrite, "f", false, "Force overwrite existing files without asking (shorthand)")
	flag.Usage = func() { // Custom Usage für bessere Hilfe
		showHelp(nil, defaultRef) // Zeige unsere Hilfe statt der Standard-Flag-Hilfe
	}
	flag.Parse() // Parse Flags zuerst

	args := flag.Args() // Get non-flag arguments

	if len(args) == 0 { // No command provided
		fmt.Println("No command specified.")
		showHelp(nil, defaultRef)
		return
	}

	commandArg := args[0] // The first non-flag argument is the command

	// Check if we should show the version (bleibt speziell, kein @ref)
	if commandArg == "-v" || commandArg == "--version" {
		fmt.Printf("templUI Component Installer v%s\n", version)
		return
	}

	// Check if we should show help (bleibt speziell)
	if commandArg == "-h" || commandArg == "--help" {
		fmt.Println("Fetching manifest for help...")
		manifest, err := fetchManifest(defaultRef)
		if err != nil {
			fmt.Println("Could not fetch component list for help:", err)
			showHelp(nil, defaultRef)
		} else {
			showHelp(&manifest, defaultRef)
		}
		return
	}

	// --- Check for 'init' command ---
	if strings.HasPrefix(commandArg, "init") {
		initRef := defaultRef // Standard-Ref

		// Parse @ref from the command argument itself (commandArg = args[0])
		if strings.Contains(commandArg, "@") {
			parts := strings.SplitN(commandArg, "@", 2)
			if len(parts) == 2 && parts[0] == "init" && parts[1] != "" {
				initRef = parts[1]
				fmt.Printf("Initializing using specified ref: %s\n", initRef)
			} else {
				fmt.Printf("Error: Invalid format '%s'. Use 'init' or 'init@<ref>'.\n", commandArg)
				return
			}
		} else if commandArg != "init" { // Verhindere z.B. "initialise"
			fmt.Printf("Error: Unknown command '%s'. Did you mean 'init'?\n", commandArg)
			showHelp(nil, defaultRef)
			return
		}

		// Prüfen, ob zusätzliche unerwartete Argumente vorhanden sind
		if len(args) > 1 { // Check length of non-flag args
			fmt.Printf("Warning: Extra arguments found after '%s'. Ignoring: %v\n", commandArg, args[1:])
		}

		initConfig(initRef, *forceOverwrite) // Übergebe den ermittelten Ref und Force-Status
		return
	}

	// --- Check for 'add' command ---
	if strings.HasPrefix(commandArg, "add") {
		targetRef := defaultRef     // Standard-Ref
		commandRefProvided := false // Flag, ob @ref mit dem add-Befehl kam

		// Parse @ref from the command argument itself (commandArg = args[0])
		if strings.Contains(commandArg, "@") {
			parts := strings.SplitN(commandArg, "@", 2)
			if len(parts) == 2 && parts[0] == "add" && parts[1] != "" {
				targetRef = parts[1]
				commandRefProvided = true
				fmt.Printf("Using specified ref from command: %s\n", targetRef)
			} else {
				fmt.Printf("Error: Invalid format '%s'. Use 'add' or 'add@<ref>'.\n", commandArg)
				return
			}
		} else if commandArg != "add" { // Verhindere z.B. "addition"
			fmt.Printf("Error: Unknown command '%s'. Did you mean 'add'?\n", commandArg)
			showHelp(nil, defaultRef)
			return
		}

		// Komponenten-Argumente beginnen bei args[1]
		if len(args) < 2 { // Need at least command + component/star
			fmt.Println("Error: No component(s) specified after 'add'.")
			fmt.Println("Usage: templui add[@<ref>] <component>... | *")
			return
		}

		// Load the local config
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			fmt.Println("Run 'templui init' to create a config file.")
			return
		}

		// --- Argument Parsing für Komponenten (ab args[1]) ---
		componentsToInstallNames := []string{} // Nur die Namen der Komponenten
		isInstallAll := false                  // Flag für '*'

		// Prüfen, ob '*' das erste Komponenten-Argument ist
		firstCompArg := args[1]
		if firstCompArg == "*" {
			if len(args) > 2 { // Darf nur '*' sein
				fmt.Println("Error: '*' must be the only component argument after 'add'.")
				fmt.Println("Usage: templui add[@<ref>] *")
				return
			}
			isInstallAll = true
		} else {
			// Parse component names (ab args[1])
			for _, arg := range args[1:] {
				// @ref bei einzelnen Komponenten wird ignoriert/führt zu Fehler, da der Ref am Befehl hängt
				if strings.Contains(arg, "@") {
					compName := strings.SplitN(arg, "@", 2)[0]
					if commandRefProvided {
						fmt.Printf("Warning: Ignoring '@...' for component '%s' because ref '%s' was specified with the 'add' command.\n", compName, targetRef)
						componentsToInstallNames = append(componentsToInstallNames, compName)
					} else {
						fmt.Printf("Error: Specify the ref with the 'add' command (e.g., 'add@%s %s'), not on individual components like '%s'.\n", targetRef, compName, arg)
						return
					}
				} else {
					componentsToInstallNames = append(componentsToInstallNames, arg)
				}
			}
		}

		// --- Download-Logik ---
		fmt.Printf("Using ref: %s\n", targetRef)

		fmt.Printf("Fetching component manifest from ref '%s'...\n", targetRef)
		manifest, err := fetchManifest(targetRef)
		if err != nil {
			if strings.Contains(err.Error(), "status code 404") {
				fmt.Printf("Error fetching manifest: %v\n", err)
				fmt.Printf("  Check if the ref '%s' exists and contains the file '%s'.\n", targetRef, manifestPath)
				fmt.Printf("  Manifest URL attempted: %s%s/%s\n", rawContentBaseURL, targetRef, manifestPath)
			} else {
				fmt.Printf("Error fetching manifest: %v\n", err)
			}
			return
		}
		fmt.Printf("Using components from templui manifest version %s (fetched from ref %s)\n", manifest.Version, targetRef)

		componentMap := make(map[string]ComponentDef)
		for _, comp := range manifest.Components {
			componentMap[comp.Name] = comp
		}

		if isInstallAll {
			fmt.Println("Preparing to install all components...")
			componentsToInstallNames = []string{}
			for _, comp := range manifest.Components {
				componentsToInstallNames = append(componentsToInstallNames, comp.Name)
			}
		}

		installedComponents := make(map[string]bool)
		requiredUtils := make(map[string]bool)

		// Install each requested component (and dependencies)
		for _, componentName := range componentsToInstallNames {
			compDef, exists := componentMap[componentName]
			if !exists {
				fmt.Printf("Error: Component '%s' not found in manifest for ref '%s'.\n", componentName, targetRef)
				fmt.Println("Available components in this manifest:")
				for _, availableComp := range manifest.Components {
					fmt.Printf(" - %s\n", availableComp.Name)
				}
				continue
			}

			// Übergebe Force-Status an installComponent
			err = installComponent(config, compDef, componentMap, targetRef, installedComponents, requiredUtils, *forceOverwrite)
			if err != nil {
				fmt.Printf("Error installing component %s: %v\n", componentName, err)
				// TODO: Entscheiden, ob hier abgebrochen werden soll
			}
		}

		// Install required utils
		if len(requiredUtils) > 0 {
			fmt.Println("Installing required utils...")
			utilsToInstallPaths := []string{}
			for utilPath := range requiredUtils {
				utilsToInstallPaths = append(utilsToInstallPaths, utilPath)
			}
			// Übergebe Force-Status an installUtils
			err = installUtils(config, utilsToInstallPaths, targetRef, *forceOverwrite)
			if err != nil {
				fmt.Printf("Error installing utils: %v\n", err)
			}
		}

		fmt.Println("\nInstallation finished.")
		return
	}

	// --- Check for 'list' command ---
	if strings.HasPrefix(commandArg, "list") {
		listRef := defaultRef // Standard-Ref

		// Parse @ref from the command argument itself
		if strings.Contains(commandArg, "@") {
			parts := strings.SplitN(commandArg, "@", 2)
			if len(parts) == 2 && parts[0] == "list" && parts[1] != "" {
				listRef = parts[1]
				fmt.Printf("Listing components using specified ref: %s\n", listRef)
			} else {
				fmt.Printf("Error: Invalid format '%s'. Use 'list' or 'list@<ref>'.\n", commandArg)
				return
			}
		} else if commandArg != "list" {
			fmt.Printf("Error: Unknown command '%s'. Did you mean 'list'?\n", commandArg)
			showHelp(nil, defaultRef)
			return
		}

		// Prüfen auf zusätzliche Argumente
		if len(args) > 1 {
			fmt.Printf("Warning: Extra arguments found after '%s'. Ignoring: %v\n", commandArg, args[1:])
		}

		err := listComponents(listRef)
		if err != nil {
			fmt.Printf("Error listing components: %v\n", err)
		}
		return
	}

	// If no known command was matched
	fmt.Printf("Error: Unknown command '%s'\n", commandArg)
	showHelp(nil, defaultRef)
}

// --- Hilfefunktion ---
func showHelp(manifest *Manifest, refUsedForHelp string) {
	// Usage angepasst
	fmt.Println("templUI Component Installer (v" + version + ")")
	fmt.Println("Usage:")
	fmt.Println("  templui [flags] init[@<ref>]         - Initialize config and install utils from <ref>")
	fmt.Println("  templui [flags] add[@<ref>] <comp>... - Add component(s) from specified <ref>")
	fmt.Println("  templui [flags] add[@<ref>] *         - Add all components from specified <ref>")
	fmt.Println("  templui [flags] list[@<ref>]        - List available components and utils from <ref>")
	fmt.Println("  templui -v, --version               - Show installer version")
	fmt.Println("  templui -h, --help                  - Show this help message")
	fmt.Println("\n<ref> can be a branch name, tag name, or commit hash.")
	fmt.Printf("If no <ref> is specified, components are fetched from the default ref (currently '%s').\n", refUsedForHelp)
	fmt.Println("\nFlags:")
	flag.PrintDefaults() // Zeigt die definierten Flags (-f, --force)

	if manifest != nil && len(manifest.Components) > 0 {
		fmt.Printf("\nAvailable components in manifest (fetched from ref '%s'):\n", refUsedForHelp)
		for _, comp := range manifest.Components {
			desc := comp.Description
			if len(desc) > 60 {
				desc = desc[:57] + "..."
			}
			fmt.Printf("  - %-15s: %s\n", comp.Name, desc)
		}
	} else {
		fmt.Println("\nCould not list available components (maybe run 'templui init' first?).")
	}
}

// --- Konfigurationsfunktionen ---
func initConfig(ref string, force bool) { // Nimmt jetzt force entgegen
	if _, err := os.Stat(configFileName); err == nil {
		// Wenn --force angegeben, überschreibe bestehende Config nicht, warne nur.
		if !force {
			fmt.Println("Config file already exists. Use --force with init to overwrite *utils* if needed, but the config file itself won't be changed.")
			// TODO: Bessere Logik? Config ggf. aktualisieren?
			// return // Hier nicht returnen, damit Utils trotzdem installiert/geprüft werden
		} else {
			fmt.Println("Config file already exists, proceeding with utils installation (--force specified).")
		}
	} else {
		// Config existiert nicht, erstelle sie
		// --- Standardwerte vorschlagen ---
		defaultComponentsDir := "./components"
		defaultUtilsDir := filepath.Join(filepath.Dir(defaultComponentsDir), "utils")
		defaultModuleName := detectModuleName()

		// --- Benutzerabfragen ---
		fmt.Printf("Enter the directory for components [%s]: ", defaultComponentsDir)
		var componentsDir string
		fmt.Scanln(&componentsDir)
		if componentsDir == "" {
			componentsDir = defaultComponentsDir
		}
		if componentsDir != defaultComponentsDir {
			defaultUtilsDir = filepath.Join(filepath.Dir(componentsDir), "utils")
		}
		fmt.Printf("Enter the directory for utils [%s]: ", defaultUtilsDir)
		var utilsDir string
		fmt.Scanln(&utilsDir)
		if utilsDir == "" {
			utilsDir = defaultUtilsDir
		}
		fmt.Printf("Enter your Go module name [%s]: ", defaultModuleName)
		var moduleName string
		fmt.Scanln(&moduleName)
		if moduleName == "" {
			moduleName = defaultModuleName
		}

		config := Config{
			ComponentsDir: componentsDir,
			UtilsDir:      utilsDir,
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
		fmt.Printf("Utils will be installed to: %s\n", config.UtilsDir)
		fmt.Printf("Using module name: %s\n", config.ModuleName)
	}

	// Lade die (eventuell neu erstellte oder vorhandene) Config, um die Pfade sicher zu haben
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config for initial utils installation: %v\n", err)
		return // Abbruch, wenn Config nicht geladen werden kann
	}

	// --- Utils direkt nach init installieren ---
	fmt.Printf("\nAttempting to install initial utils from ref '%s'...\n", ref)
	manifest, err := fetchManifest(ref)
	if err != nil {
		if strings.Contains(err.Error(), "status code 404") {
			fmt.Printf("Warning: Could not fetch manifest from ref '%s': %v\n", ref, err)
			fmt.Printf("  Check if the ref '%s' exists and contains the file '%s'.\n", ref, manifestPath)
		} else {
			fmt.Printf("Warning: Could not fetch manifest to install initial utils: %v\n", err)
		}
		return
	}

	if len(manifest.Utils) == 0 {
		fmt.Println("No utils defined in the manifest. Skipping initial utils installation.")
		return
	}

	allUtilPaths := []string{}
	fmt.Println("Found utils in manifest:")
	for _, utilDef := range manifest.Utils {
		allUtilPaths = append(allUtilPaths, utilDef.Path)
		fmt.Printf(" - %s\n", utilDef.Path)
	}

	// Verwende den force-Status von initConfig
	err = installUtils(config, allUtilPaths, ref, *forceOverwrite)
	if err != nil {
		fmt.Printf("Error during initial utils installation: %v\n", err)
	} else {
		fmt.Println("Initial utils installation completed.")
	}
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
	// Validate config
	if config.ComponentsDir == "" || config.ModuleName == "" || config.UtilsDir == "" { // Prüfung für UtilsDir hinzugefügt
		return config, fmt.Errorf("invalid config: ComponentsDir, UtilsDir, and ModuleName must be set")
	}
	return config, nil
}

// --- Download-Hilfsfunktionen ---

// fetchManifest lädt das Manifest für einen gegebenen Git ref (branch/tag/commit)
func fetchManifest(ref string) (Manifest, error) {
	manifestURL := rawContentBaseURL + ref + "/" + manifestPath
	resp, err := http.Get(manifestURL)
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to start download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return Manifest{}, fmt.Errorf("failed to download manifest from %s: status code %d, message: %s", manifestURL, resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var manifest Manifest
	err = json.Unmarshal(body, &manifest)
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to parse manifest JSON (from %s): %w", manifestURL, err)
	}

	return manifest, nil
}

// downloadFile lädt eine einzelne Datei von einer URL
func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to start download from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to download file from %s: status code %d, message: %s", url, resp.StatusCode, string(bodyBytes))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", url, err)
	}
	return data, nil
}

// --- Installationslogik ---

// installComponent installiert eine Komponente und ihre Abhängigkeiten
func installComponent(
	config Config,
	comp ComponentDef,
	componentMap map[string]ComponentDef,
	ref string,
	installed map[string]bool,
	requiredUtils map[string]bool,
	force bool, // force parameter hinzugefügt
) error {
	if installed[comp.Name] {
		return nil
	}
	installed[comp.Name] = true

	fmt.Printf("Processing component: %s (from ref: %s)\n", comp.Name, ref)

	// 1. Abhängigkeiten zuerst installieren (rekursiv)
	for _, depName := range comp.Dependencies {
		if installed[depName] {
			continue
		}
		depComp, exists := componentMap[depName]
		if !exists {
			fmt.Printf("Warning: Dependency '%s' for component '%s' not found in manifest for ref '%s'. Skipping dependency.\n", depName, comp.Name, ref)
			continue
		}
		// Übergebe force an rekursive Aufrufe
		err := installComponent(config, depComp, componentMap, ref, installed, requiredUtils, force)
		if err != nil {
			return fmt.Errorf("failed to install dependency '%s' for '%s': %w", depName, comp.Name, err)
		}
		fmt.Printf(" -> Installed dependency: %s\n", depName)
	}

	// 2. Dateien der aktuellen Komponente herunterladen und schreiben
	fmt.Printf(" Installing files for: %s\n", comp.Name)
	repoComponentBasePath := "internal/components/"

	for _, repoFilePath := range comp.Files {
		var destPath string
		if strings.HasPrefix(repoFilePath, repoComponentBasePath) {
			relativePath := repoFilePath[len(repoComponentBasePath):]
			destPath = filepath.Join(config.ComponentsDir, relativePath)
		} else {
			fmt.Printf("  Warning: File path '%s' does not start with '%s'. Placing it directly in '%s'.\n", repoFilePath, repoComponentBasePath, config.ComponentsDir)
			fileName := filepath.Base(repoFilePath)
			destPath = filepath.Join(config.ComponentsDir, fileName)
		}

		compDestDir := filepath.Dir(destPath)
		err := os.MkdirAll(compDestDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create destination directory '%s': %w", compDestDir, err)
		}

		// --- Überschreiben-Logik ---
		fileExists := false
		if _, err := os.Stat(destPath); err == nil {
			fileExists = true
		}

		shouldWriteFile := true // Standardmäßig schreiben (wenn Datei nicht existiert oder überschrieben wird)

		if fileExists {
			existingRef, _ := readFileVersion(destPath) // Fehler hier ignorieren? Oder loggen?
			if existingRef == ref {
				fmt.Printf("  Info: File '%s' already up-to-date (ref: %s). Skipping.\n", destPath, ref)
				shouldWriteFile = false // Nicht schreiben
			} else {
				// Versionen unterschiedlich oder alte Version nicht lesbar
				if force {
					fmt.Printf("  Info: File '%s' exists (Version: '%s'). Forcing overwrite with ref '%s'.\n", destPath, existingRef, ref)
					// shouldWriteFile bleibt true
				} else {
					shouldOverwrite := askForOverwrite(destPath, existingRef, ref)
					if !shouldOverwrite {
						fmt.Printf("  Info: Skipping overwrite for '%s'.\n", destPath)
						shouldWriteFile = false // Nicht schreiben
					}
					// Ansonsten: shouldWriteFile bleibt true
				}
			}
		}

		// Nur weitermachen (Download + Schreiben), wenn nötig
		if shouldWriteFile {
			// --- Download ---
			fileURL := rawContentBaseURL + ref + "/" + repoFilePath
			fmt.Printf("   Downloading %s...\n", fileURL)
			data, err := downloadFile(fileURL)
			if err != nil {
				fileNameForError := filepath.Base(repoFilePath)
				return fmt.Errorf("failed to download file '%s' for component '%s' from %s: %w", fileNameForError, comp.Name, fileURL, err)
			}

			// --- Modifikationen ---
			versionComment := fmt.Sprintf("// templui component %s - version: %s installed by templui v%s\n", comp.Name, ref, version)
			modifiedData := append([]byte(versionComment), data...)
			if strings.HasSuffix(repoFilePath, ".templ") || strings.HasSuffix(repoFilePath, ".go") {
				modifiedData = replaceImports(modifiedData, config.ModuleName, comp.Name)
			}

			// --- Schreiben ---
			err = os.WriteFile(destPath, modifiedData, 0644)
			if err != nil {
				return fmt.Errorf("failed to write file '%s': %w", destPath, err)
			}
			if fileExists { // Nur melden, wenn tatsächlich überschrieben wurde
				fmt.Printf("   Overwritten %s\n", destPath)
			} else {
				fmt.Printf("   Installed %s\n", destPath)
			}
		}
	}

	// 3. Benötigte Utils sammeln
	for _, repoUtilPath := range comp.RequiredUtils {
		requiredUtils[repoUtilPath] = true
	}

	return nil
}

// installUtils installiert die benötigten Util-Dateien
func installUtils(config Config, utilPaths []string, ref string, force bool) error { // force parameter hinzugefügt
	if len(utilPaths) == 0 {
		return nil
	}

	utilsBaseDestDir := config.UtilsDir
	fmt.Printf("Ensuring utils are installed in: %s (from ref: %s)\n", utilsBaseDestDir, ref)
	repoUtilBasePath := "internal/utils/"

	err := os.MkdirAll(utilsBaseDestDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create base utils directory '%s': %w", utilsBaseDestDir, err)
	}

	for _, repoUtilPath := range utilPaths {
		var destPath string
		if strings.HasPrefix(repoUtilPath, repoUtilBasePath) {
			relativePath := repoUtilPath[len(repoUtilBasePath):]
			destPath = filepath.Join(utilsBaseDestDir, relativePath)
		} else {
			fmt.Printf("  Warning: Util path '%s' does not start with '%s'. Placing it directly in '%s'.\n", repoUtilPath, repoUtilBasePath, utilsBaseDestDir)
			fileName := filepath.Base(repoUtilPath)
			destPath = filepath.Join(utilsBaseDestDir, fileName)
		}

		utilDestDir := filepath.Dir(destPath)
		err := os.MkdirAll(utilDestDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create destination utils directory '%s': %w", utilDestDir, err)
		}

		// --- Überschreiben-Logik ---
		fileExists := false
		if _, err := os.Stat(destPath); err == nil {
			fileExists = true
		}

		shouldWriteFile := true // Standardmäßig schreiben

		if fileExists {
			existingRef, _ := readFileVersion(destPath)
			if existingRef == ref {
				fmt.Printf("  Info: Util file '%s' already up-to-date (ref: %s). Skipping.\n", destPath, ref)
				shouldWriteFile = false
			} else {
				if force {
					fmt.Printf("  Info: Util file '%s' exists (Version: '%s'). Forcing overwrite with ref '%s'.\n", destPath, existingRef, ref)
				} else {
					shouldOverwrite := askForOverwrite(destPath, existingRef, ref)
					if !shouldOverwrite {
						fmt.Printf("  Info: Skipping overwrite for '%s'.\n", destPath)
						shouldWriteFile = false
					}
				}
			}
		}

		if shouldWriteFile {
			// --- Download ---
			fileURL := rawContentBaseURL + ref + "/" + repoUtilPath
			fmt.Printf("   Downloading util %s...\n", fileURL)
			data, err := downloadFile(fileURL)
			if err != nil {
				fileNameForError := filepath.Base(repoUtilPath)
				return fmt.Errorf("failed to download util '%s' from %s: %w", fileNameForError, fileURL, err)
			}

			// --- Modifikationen ---
			utilNameForComment := filepath.Base(repoUtilPath)
			versionComment := fmt.Sprintf("// templui util %s - version: %s installed by templui v%s\n", utilNameForComment, ref, version)
			modifiedData := append([]byte(versionComment), data...)
			if strings.HasSuffix(repoUtilPath, ".go") {
				modifiedData = replaceImports(modifiedData, config.ModuleName, "")
			}

			// --- Schreiben ---
			err = os.WriteFile(destPath, modifiedData, 0644)
			if err != nil {
				return fmt.Errorf("failed to write util file '%s': %w", destPath, err)
			}
			if fileExists {
				fmt.Printf("   Overwritten %s\n", destPath)
			} else {
				fmt.Printf("   Installed %s\n", destPath)
			}
		}
	}

	return nil
}

// --- Neue Funktion für den list-Befehl ---

// listComponents fetches the manifest for a given ref and lists available components.
func listComponents(ref string) error {
	fmt.Printf("Fetching component manifest from ref '%s'...\n", ref)
	manifest, err := fetchManifest(ref)
	if err != nil {
		// Spezifische Fehlermeldung für 404
		if strings.Contains(err.Error(), "status code 404") {
			return fmt.Errorf("could not fetch manifest: ref '%s' not found or does not contain '%s'", ref, manifestPath)
		}
		return fmt.Errorf("could not fetch manifest: %w", err)
	}

	fmt.Printf("\nAvailable components in ref '%s' (Manifest Version: %s):\n", ref, manifest.Version)
	if len(manifest.Components) == 0 {
		fmt.Println("  No components found in this manifest.")
		return nil
	}

	// Komponenten ausgeben (optional sortieren)
	// sort.Slice(manifest.Components, func(i, j int) bool {
	// 	return manifest.Components[i].Name < manifest.Components[j].Name
	// })

	for _, comp := range manifest.Components {
		desc := comp.Description
		if len(desc) > 60 { // Beschreibung kürzen
			desc = desc[:57] + "..."
		}
		// Formatierte Ausgabe
		fmt.Printf("  - %-20s : %s\n", comp.Name, desc)
	}

	// Optional: Auch Utils auflisten?
	if len(manifest.Utils) > 0 {
		fmt.Printf("\nAvailable utils in ref '%s':\n", ref)
		for _, util := range manifest.Utils {
			// Nur den Dateinamen oder den relativen Pfad anzeigen
			utilName := filepath.Base(util.Path)
			if util.Description != "" {
				desc := util.Description
				if len(desc) > 50 {
					desc = desc[:47] + "..."
				}
				fmt.Printf("  - %-20s : %s\n", utilName, desc)
			} else {
				fmt.Printf("  - %s\n", utilName)
			}
		}
	}

	return nil
}

// --- Neue Hilfsfunktionen für Überschreiben ---

// readFileVersion liest die Version aus dem Kommentar in der ersten Zeile
func readFileVersion(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		// Wenn Datei nicht existiert, ist das kein Fehler für diese Funktion
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Nur die erste Zeile lesen sollte reichen
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		matches := versionRegex.FindStringSubmatch(line)
		// Index 1 enthält die erste Capturing Group (den Ref)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	// Kein Kommentar gefunden oder Fehler beim Scannen
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning file %s: %w", filePath, err)
	}

	return "", nil // Kein Fehler, aber auch keine Version gefunden
}

// askForOverwrite fragt den Benutzer, ob eine Datei überschrieben werden soll
func askForOverwrite(filePath, oldRef, newRef string) bool {
	reader := bufio.NewReader(os.Stdin)

	oldVersionStr := oldRef
	if oldVersionStr == "" {
		oldVersionStr = "<unknown or no comment>"
	}

	fmt.Printf("  Confirm: File '%s' (Existing Version: %s) differs from requested ref '%s'. Overwrite? [y/N]: ",
		filePath, oldVersionStr, newRef)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nError reading input: %v. Assuming No.\n", err)
		return false
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y"
}

// replaceImports ersetzt die internen templUI-Importpfade durch den Modulnamen des Benutzers
func replaceImports(data []byte, userModuleName string, context string) []byte { // context hinzugefügt
	content := string(data)
	// Das Muster muss exakt den Importpfad matchen, der in deinen internen Komponenten/Utils verwendet wird!
	// Annahme: Deine internen Imports sehen aus wie "github.com/axzilla/templui/internal/..."
	// Passe <user>/<repo> an dein Repo an!
	internalImportPattern := `"github.com/axzilla/templui/internal/([^"]+)"`
	// targetImportFormat bildet den neuen Pfad basierend auf dem Modulnamen des Users
	// z.B. wenn $1="utils", wird es zu "your/module/path/utils"
	// z.B. wenn $1="components/icon", wird es zu "your/module/path/components/icon" (falls benötigt)
	targetImportFormat := fmt.Sprintf(`"%s/$1"`, userModuleName)

	re := regexp.MustCompile(internalImportPattern)
	newContent := re.ReplaceAllString(content, targetImportFormat)

	// Logge, ob Ersetzungen stattgefunden haben (optional für Debugging)
	if content != newContent {
		logPrefix := "    ->"
		if context != "" {
			logPrefix = fmt.Sprintf("    -> [%s]", context)
		}
		fmt.Printf("%s Replaced import paths.\n", logPrefix)
	} else {
		//fmt.Println("    -> No import paths needed replacement.") // Weniger verbose
	}

	return []byte(newContent)
}

// --- Alte, nicht mehr verwendete Funktionen ---
// Sind nicht mehr nötig, da wir das Manifest verwenden.
