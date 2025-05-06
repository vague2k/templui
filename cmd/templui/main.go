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
	if len(os.Args) < 2 { // Mindestens ein Argument (der Befehl) wird benötigt
		fmt.Println("No command specified.")
		showHelp(nil, defaultRef)
		return
	}

	commandArg := os.Args[1] // Das erste Argument (z.B. "init", "init@ref", "add", "add@ref", "-v")

	// Check if we should show the version
	// Behält die alte Logik bei, da -v und --version keine @ref haben
	if commandArg == "-v" || commandArg == "--version" {
		fmt.Printf("templUI Component Installer v%s\n", version)
		return
	}

	// Check if we should show help
	// Behält die alte Logik bei
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

		// Parse @ref from the command argument itself (commandArg = os.Args[1])
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
		if len(os.Args) > 2 {
			fmt.Printf("Warning: Extra arguments found after '%s'. Ignoring: %v\n", commandArg, os.Args[2:])
		}

		initConfig(initRef) // Übergebe den ermittelten Ref
		return
	}

	// --- Check for 'add' command ---
	if strings.HasPrefix(commandArg, "add") {
		targetRef := defaultRef     // Standard-Ref
		commandRefProvided := false // Flag, ob @ref mit dem add-Befehl kam

		// Parse @ref from the command argument itself (commandArg = os.Args[1])
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

		// Komponenten-Argumente beginnen bei os.Args[2]
		if len(os.Args) < 3 {
			fmt.Println("Error: No component(s) specified after 'add'.")
			fmt.Println("Usage: templui add[@<ref>] <component>[@<ref>]... | *")
			return
		}

		// Load the local config
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			fmt.Println("Run 'templui init' to create a config file.")
			return
		}

		// --- Argument Parsing für Komponenten (ab os.Args[2]) ---
		componentsToInstallNames := []string{} // Nur die Namen der Komponenten
		isInstallAll := false                  // Flag für '*'

		// Prüfen, ob '*' das erste Komponenten-Argument ist
		firstCompArg := os.Args[2]
		if firstCompArg == "*" {
			if len(os.Args) > 3 { // Darf nur '*' sein
				fmt.Println("Error: '*' must be the only component argument after 'add'.")
				fmt.Println("Usage: templui add[@<ref>] *")
				return
			}
			isInstallAll = true
			// Ref wurde bereits aus commandArg geparst (targetRef)
		} else {
			// Parse 'component' (ohne @ref, da der Ref vom Befehl kommt)
			// Die Argumente sind jetzt os.Args[2:]
			for _, arg := range os.Args[2:] {
				// Wenn ein Ref mit dem *Befehl* kam, ignoriere @ref bei Komponenten
				if commandRefProvided && strings.Contains(arg, "@") {
					parts := strings.SplitN(arg, "@", 2)
					compName := parts[0]
					fmt.Printf("Warning: Ignoring '@%s' for component '%s' because ref '%s' was specified with the 'add' command.\n", parts[1], compName, targetRef)
					componentsToInstallNames = append(componentsToInstallNames, compName)
				} else if strings.Contains(arg, "@") {
					// Wenn KEIN Ref mit dem Befehl kam, ABER hier eins ist -> Fehler (oder wir erlauben es doch?)
					// Aktuell: Fehler, da wir oben entschieden haben, dass der Befehls-Ref gilt.
					// Alternative: targetRef hier überschreiben, aber nur wenn es das erste ist. Komplex.
					// Einfacher: @ref nur am Befehl oder gar nicht erlauben.
					// Für jetzt lassen wir es wie oben: Wenn commandRefProvided, wird ignoriert.
					// Wenn !commandRefProvided und @ hier -> Fehler (oder targetRef setzen?)
					// Wir bleiben dabei: @Ref gehört an den Befehl `add` oder `init`.
					fmt.Printf("Error: Specify the ref with the 'add' command (e.g., 'add@%s %s'), not on individual components.\n", defaultRef, arg)
					return

				} else {
					// Normaler Komponentenname ohne @ref
					componentsToInstallNames = append(componentsToInstallNames, arg)
				}
			}
		}

		// --- Download-Logik ---
		fmt.Printf("Using ref: %s\n", targetRef)

		fmt.Printf("Fetching component manifest from ref '%s'...\n", targetRef)
		manifest, err := fetchManifest(targetRef)
		if err != nil {
			// Gib eine spezifischere Fehlermeldung aus, wenn 404 aufgetreten ist
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

		// Create map for easy lookup
		componentMap := make(map[string]ComponentDef)
		for _, comp := range manifest.Components {
			componentMap[comp.Name] = comp
		}

		// Wenn '*' angefordert wurde, jetzt die Namen aus dem Manifest holen
		if isInstallAll {
			fmt.Println("Preparing to install all components...")
			componentsToInstallNames = []string{} // Leeren und neu befüllen
			for _, comp := range manifest.Components {
				componentsToInstallNames = append(componentsToInstallNames, comp.Name)
			}
		}

		// Track installed components and required utils to avoid duplicates
		installedComponents := make(map[string]bool)
		requiredUtils := make(map[string]bool) // Map von Util-Pfad zu bool

		// Install each requested component (and dependencies)
		for _, componentName := range componentsToInstallNames {
			compDef, exists := componentMap[componentName]
			if !exists {
				fmt.Printf("Error: Component '%s' not found in manifest for ref '%s'.\n", componentName, targetRef)
				fmt.Println("Available components in this manifest:")
				for _, availableComp := range manifest.Components {
					fmt.Printf(" - %s\n", availableComp.Name)
				}
				continue // oder return error
			}

			err = installComponent(config, compDef, componentMap, targetRef, installedComponents, requiredUtils)
			if err != nil {
				fmt.Printf("Error installing component %s: %v\n", componentName, err)
			}
		}

		// Install required utils
		if len(requiredUtils) > 0 {
			fmt.Println("Installing required utils...")
			utilsToInstallPaths := []string{}
			for utilPath := range requiredUtils {
				utilsToInstallPaths = append(utilsToInstallPaths, utilPath)
			}
			err = installUtils(config, utilsToInstallPaths, targetRef)
			if err != nil {
				fmt.Printf("Error installing utils: %v\n", err)
			}
		}

		fmt.Println("\nInstallation finished.")
		return // Ende des add-Blocks
	} // Ende 'if strings.HasPrefix(commandArg, "add")'

	// If no known command was matched
	fmt.Printf("Error: Unknown command '%s'\n", commandArg)
	showHelp(nil, defaultRef)
}

// --- Hilfefunktion ---
func showHelp(manifest *Manifest, refUsedForHelp string) {
	fmt.Println("templUI Component Installer (v" + version + ")")
	fmt.Println("Usage:")
	fmt.Println("  templui init[@<ref>]         - Initialize the config file and install utils from <ref>")
	fmt.Println("  templui add[@<ref>] <comp>... - Add component(s) from specified <ref>")
	fmt.Println("  templui add[@<ref>] *         - Add all available components from a specific ref")
	fmt.Println("  templui -v, --version       - Show installer version")
	fmt.Println("  templui -h, --help          - Show this help message")
	fmt.Println("\n<ref> can be a branch name, tag name, or commit hash.")
	fmt.Printf("If no <ref> is specified, components are fetched from the default ref (currently '%s').\n", refUsedForHelp)

	// Komponenten auflisten, wenn Manifest vorhanden ist
	if manifest != nil && len(manifest.Components) > 0 {
		fmt.Printf("\nAvailable components in manifest (fetched from ref '%s'):\n", refUsedForHelp)
		// Sortieren für bessere Lesbarkeit? (Optional)
		// sort.Slice(manifest.Components, func(i, j int) bool {
		// 	return manifest.Components[i].Name < manifest.Components[j].Name
		// })
		for _, comp := range manifest.Components {
			// Beschreibung ggf. kürzen oder nur erste Zeile anzeigen?
			desc := comp.Description
			if len(desc) > 60 {
				desc = desc[:57] + "..."
			}
			fmt.Printf("  - %-15s: %s\n", comp.Name, desc) // Formatierung für Ausrichtung
		}
	} else {
		fmt.Println("\nCould not list available components (maybe run 'templui init' first?).")
	}
}

// --- Konfigurationsfunktionen ---
func initConfig(ref string) {
	if _, err := os.Stat(configFileName); err == nil {
		fmt.Println("Config file already exists.")
		// TODO: Fragen, ob überschrieben werden soll?
		return
	}

	// --- Standardwerte vorschlagen ---
	defaultComponentsDir := "./components"
	// Schlage Utils-Verzeichnis relativ zum Komponenten-Verzeichnis vor
	defaultUtilsDir := filepath.Join(filepath.Dir(defaultComponentsDir), "utils") // Ergibt "./utils"
	defaultModuleName := detectModuleName()

	// --- Benutzerabfragen ---
	fmt.Printf("Enter the directory for components [%s]: ", defaultComponentsDir)
	var componentsDir string
	fmt.Scanln(&componentsDir)
	if componentsDir == "" {
		componentsDir = defaultComponentsDir
	}
	// Passe den Vorschlag für UtilsDir an, falls ComponentsDir geändert wurde
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

	// --- Config erstellen und speichern ---
	config := Config{
		ComponentsDir: componentsDir,
		UtilsDir:      utilsDir, // Hinzugefügt
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
	fmt.Printf("Utils will be installed to: %s\n", config.UtilsDir) // Info hinzugefügt
	fmt.Printf("Using module name: %s\n", config.ModuleName)

	// --- Utils direkt nach init installieren ---
	fmt.Printf("\nAttempting to install initial utils from ref '%s'...\n", ref)
	manifest, err := fetchManifest(ref)
	if err != nil {
		// Gib spezifischere Fehlermeldung aus
		if strings.Contains(err.Error(), "status code 404") {
			fmt.Printf("Warning: Could not fetch manifest from ref '%s': %v\n", ref, err)
			fmt.Printf("  Check if the ref '%s' exists and contains the file '%s'.\n", ref, manifestPath)
		} else {
			fmt.Printf("Warning: Could not fetch manifest to install initial utils: %v\n", err)
		}
		return // Beende hier, da Utils nicht installiert werden können
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

	err = installUtils(config, allUtilPaths, ref) // Verwende übergebenen ref
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
	// Bereinige den Ref (optional, entfernt unsichere Zeichen)
	// ref = url.PathEscape(ref) // Kann Probleme machen, wenn Ref '/' enthält, z.B. release/v1
	manifestURL := rawContentBaseURL + ref + "/" + manifestPath
	// fmt.Println("Downloading manifest from:", manifestURL) // Debugging

	resp, err := http.Get(manifestURL)
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to start download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Versuch, Fehlermeldung zu lesen
		return Manifest{}, fmt.Errorf("failed to download manifest from %s: status code %d, message: %s", manifestURL, resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var manifest Manifest
	err = json.Unmarshal(body, &manifest)
	if err != nil {
		// Versuche, eine bessere Fehlermeldung für ungültiges JSON zu geben
		return Manifest{}, fmt.Errorf("failed to parse manifest JSON (from %s): %w", manifestURL, err)
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
		bodyBytes, _ := io.ReadAll(resp.Body) // Versuch, Fehlermeldung zu lesen
		return nil, fmt.Errorf("failed to download file from %s: status code %d, message: %s", url, resp.StatusCode, string(bodyBytes))
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
	ref string, // Geändert zu ref
	installed map[string]bool,
	requiredUtils map[string]bool,
) error {
	// Verhindere Endlosschleifen und doppelte Installation
	if installed[comp.Name] {
		// Bereits verarbeitet in diesem Durchlauf
		return nil
	}
	installed[comp.Name] = true // Markiere als (versucht) installiert

	fmt.Printf("Processing component: %s (from ref: %s)\n", comp.Name, ref)

	// 1. Abhängigkeiten zuerst installieren (rekursiv)
	for _, depName := range comp.Dependencies {
		// Prüfe, ob Abhängigkeit bereits verarbeitet wurde, um Zyklen zu vermeiden
		// (obwohl `installed` dies bereits tut, zusätzliche Klarheit)
		if installed[depName] {
			continue
		}

		depComp, exists := componentMap[depName]
		if !exists {
			// Sollte nicht passieren, wenn Manifest korrekt ist und alle Deps auflistet
			fmt.Printf("Warning: Dependency '%s' for component '%s' not found in manifest for ref '%s'. Skipping dependency.\n", depName, comp.Name, ref)
			continue
		}
		// Installiere Abhängigkeit (vom gleichen Ref wie die Hauptkomponente)
		err := installComponent(config, depComp, componentMap, ref, installed, requiredUtils)
		if err != nil {
			// Fehler bei Abhängigkeit weitergeben
			return fmt.Errorf("failed to install dependency '%s' for '%s': %w", depName, comp.Name, err)
		}
		fmt.Printf(" -> Installed dependency: %s\n", depName) // Info nach erfolgreicher Installation
	}

	// 2. Dateien der aktuellen Komponente herunterladen und schreiben
	fmt.Printf(" Installing files for: %s\n", comp.Name)
	repoComponentBasePath := "internal/components/" // Basispfad im Repo

	for _, repoFilePath := range comp.Files { // z.B. "internal/components/aspectratio/aspect_ratio.templ"

		// --- Zielpfad bestimmen ---
		var destPath string
		if strings.HasPrefix(repoFilePath, repoComponentBasePath) {
			// Relativen Pfad extrahieren (z.B. "aspectratio/aspect_ratio.templ")
			relativePath := repoFilePath[len(repoComponentBasePath):]
			// Zielpfad konstruieren (z.B. "./components/aspectratio/aspect_ratio.templ")
			destPath = filepath.Join(config.ComponentsDir, relativePath)
		} else {
			// Fallback für unerwartete Pfade (sollte nicht passieren bei korrekter Manifest-Pflege)
			fmt.Printf("  Warning: File path '%s' does not start with '%s'. Placing it directly in '%s'.\n", repoFilePath, repoComponentBasePath, config.ComponentsDir)
			fileName := filepath.Base(repoFilePath)
			destPath = filepath.Join(config.ComponentsDir, fileName)
		}

		// Stelle sicher, dass das *gesamte* Zielverzeichnis existiert (inkl. Unterordner)
		compDestDir := filepath.Dir(destPath) // z.B. "./components/aspectratio"
		err := os.MkdirAll(compDestDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create destination directory '%s': %w", compDestDir, err)
		}

		// --- Überschreiben-Logik ---
		if _, err := os.Stat(destPath); err == nil {
			// Datei existiert
			// TODO: Implementiere Logik zum Lesen der Version aus der Datei
			// TODO: Frage den Benutzer, ob überschrieben werden soll (bufio.NewReader)
			fmt.Printf("  Warning: File '%s' already exists. Skipping overwrite for now.\n", destPath)
			continue // Überspringe diese Datei
		}

		// --- Download ---
		fileURL := rawContentBaseURL + ref + "/" + repoFilePath
		fmt.Printf("   Downloading %s...\n", fileURL)
		data, err := downloadFile(fileURL)
		if err != nil {
			fileNameForError := filepath.Base(repoFilePath) // Nur Dateiname für Fehlermeldung
			return fmt.Errorf("failed to download file '%s' for component '%s' from %s: %w", fileNameForError, comp.Name, fileURL, err)
		}

		// --- Modifikationen ---
		// a) Versionskommentar hinzufügen
		versionComment := fmt.Sprintf("// templui component %s - version: %s installed by templui v%s\n", comp.Name, ref, version)
		modifiedData := append([]byte(versionComment), data...)

		// b) Importpfade ersetzen (nur für .templ und .go Dateien relevant?)
		if strings.HasSuffix(repoFilePath, ".templ") || strings.HasSuffix(repoFilePath, ".go") {
			modifiedData = replaceImports(modifiedData, config.ModuleName, comp.Name) // Pass component name for better context
		}

		// --- Schreiben ---
		err = os.WriteFile(destPath, modifiedData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file '%s': %w", destPath, err)
		}
		fmt.Printf("   Installed %s\n", destPath)
	}

	// 3. Benötigte Utils für diese Komponente sammeln (Pfade im Repo)
	for _, repoUtilPath := range comp.RequiredUtils {
		requiredUtils[repoUtilPath] = true
	}

	return nil
}

// installUtils anpassen, um config.UtilsDir zu verwenden
func installUtils(config Config, utilPaths []string, ref string) error { // Geändert zu ref
	if len(utilPaths) == 0 {
		return nil // Nichts zu tun
	}

	// Verwende UtilsDir direkt aus der Config
	utilsBaseDestDir := config.UtilsDir
	fmt.Printf("Ensuring utils are installed in: %s (from ref: %s)\n", utilsBaseDestDir, ref)
	repoUtilBasePath := "internal/utils/" // Basispfad im Repo

	// Stelle sicher, dass das Basis-Utils-Verzeichnis existiert
	// MkdirAll im Loop unten kümmert sich auch darum, aber hier explizit schadet nicht.
	err := os.MkdirAll(utilsBaseDestDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create base utils directory '%s': %w", utilsBaseDestDir, err)
	}

	for _, repoUtilPath := range utilPaths { // z.B. "internal/utils/classname.go" oder "internal/utils/forms/validation.go"

		// --- Zielpfad bestimmen ---
		var destPath string
		if strings.HasPrefix(repoUtilPath, repoUtilBasePath) {
			relativePath := repoUtilPath[len(repoUtilBasePath):]
			destPath = filepath.Join(utilsBaseDestDir, relativePath)
		} else {
			fmt.Printf("  Warning: Util path '%s' does not start with '%s'. Placing it directly in '%s'.\n", repoUtilPath, repoUtilBasePath, utilsBaseDestDir)
			fileName := filepath.Base(repoUtilPath)
			destPath = filepath.Join(utilsBaseDestDir, fileName)
		}

		// Stelle sicher, dass das *gesamte* Zielverzeichnis existiert
		utilDestDir := filepath.Dir(destPath)
		err := os.MkdirAll(utilDestDir, 0755) // <- Doppelt geprüft, schadet aber nicht
		if err != nil {
			return fmt.Errorf("failed to create destination utils directory '%s': %w", utilDestDir, err)
		}

		// --- Überschreiben-Logik (ähnlich wie bei Komponenten) ---
		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("  Info: Util file '%s' already exists. Skipping download.\n", destPath) // Geändert zu Info
			continue                                                                            // Überspringe diese Datei
		}

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
		fmt.Printf("   Installed %s\n", destPath)
	}

	return nil
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
