# templUI Component Installer

Ein CLI-Tool zum einfachen Installieren von templUI-Komponenten in dein Projekt, ähnlich wie shadcn/ui.

## Installation

Du kannst das Tool global mit Go installieren:

```bash
go install github.com/axzilla/templui/cmd/compinstall@latest
```

Nach der Installation kannst du das Tool von überall aus mit dem Befehl `compinstall` ausführen.

## Verwendung

### Initialisierung

Zuerst musst du das Tool in deinem Projekt initialisieren:

```bash
compinstall init
```

Dies erstellt eine Konfigurationsdatei `.templui.json` in deinem Projektverzeichnis, die den Pfad zu deinem Komponenten-Verzeichnis und deinen Go-Modul-Namen speichert.

Bei der Initialisierung wird versucht, deinen Go-Modul-Namen automatisch aus der `go.mod`-Datei zu erkennen. Du kannst diesen aber auch manuell angeben.

### Komponenten hinzufügen

Um eine Komponente zu deinem Projekt hinzuzufügen:

```bash
compinstall add button
```

Dies installiert die Button-Komponente und alle ihre Abhängigkeiten in das konfigurierte Verzeichnis. Dabei werden die Imports in den Komponenten automatisch an deinen Go-Modul-Namen angepasst.

Bei jeder Installation einer Komponente werden auch die benötigten `utils` automatisch in einen `utils`-Ordner neben dem Komponenten-Verzeichnis installiert.

### Hilfe anzeigen

Um die Hilfe und eine Liste aller verfügbaren Komponenten anzuzeigen:

```bash
compinstall -h
```

## Beispiel

```bash
# Navigiere zu deinem Projekt
cd mein-projekt

# Initialisiere TemplUI Component Installer
compinstall init
# Gib deinen Go-Modul-Namen ein, z.B. github.com/username/mein-projekt

# Füge die Button-Komponente hinzu
compinstall add button

# Füge die Card-Komponente hinzu
compinstall add card
```

## Optionen

- `init`: Initialisiert die Konfigurationsdatei
- `add <component>`: Fügt eine Komponente zu deinem Projekt hinzu
- `-v`, `--version`: Zeigt die Version des Tools an
- `-h`, `--help`: Zeigt die Hilfe und eine Liste aller verfügbaren Komponenten an

## Ähnlich wie shadcn/ui

Dieses Tool ist inspiriert von shadcn/ui und bietet eine ähnliche Erfahrung für templUI-Komponenten:

```bash
# shadcn/ui
npx shadcn@latest init
npx shadcn@latest add button

# templUI Component Installer
compinstall init
compinstall add button
```

Die Verwendung ist fast identisch, was es für Benutzer, die bereits mit shadcn/ui vertraut sind, sehr einfach macht.

## Anpassung der Imports

Eine besondere Funktion des Tools ist die automatische Anpassung der Imports in den Komponenten. Wenn eine Komponente ursprünglich Imports wie `"github.com/axzilla/templui/internal/utilss"` enthält, werden diese automatisch zu `"github.com/username/mein-projekt/utils"` (oder was auch immer dein Modul-Name ist) geändert.

Dies stellt sicher, dass die Komponenten in deinem Projekt ohne weitere Anpassungen funktionieren.

## Installation von Utils

Das Tool installiert bei jeder Komponenten-Installation automatisch auch die benötigten Utility-Funktionen. Diese werden in einen `utils`-Ordner neben dem Komponenten-Verzeichnis kopiert. Die Imports in den Utility-Dateien werden ebenfalls an deinen Modul-Namen angepasst.

Die `utils` werden immer installiert, unabhängig davon, welche Komponente du installierst, um sicherzustellen, dass alle Komponenten korrekt funktionieren.
