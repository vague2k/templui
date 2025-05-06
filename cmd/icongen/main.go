package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	outputFile     = "./icon/icon_defs.go"
	iconContentDir = "./icon/content" // Directory for individual icon contents
	lucideVersion  = "0.452.0"        // Current Lucide version - update as needed
)

type GitHubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	DownloadUrl string `json:"download_url"`
}

func main() {
	// Fetch the list of files from GitHub
	files, err := fetchGitHubContents("")
	if err != nil {
		panic(fmt.Errorf("failed to fetch GitHub contents: %w", err))
	}

	// Initialize slice for icon definitions
	var iconDefs []string
	iconDefs = append(iconDefs, "package icon\n")
	iconDefs = append(iconDefs, "// This file is auto generated\n")
	iconDefs = append(iconDefs, fmt.Sprintf("// Using Lucide icons version %s\n", lucideVersion))

	// Create the content directory if it doesn't exist
	err = os.MkdirAll(iconContentDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Create a wait group to synchronize the goroutines
	var wg sync.WaitGroup
	// Create a channel to limit the number of concurrent downloads
	semaphore := make(chan struct{}, 10) // Allow up to 10 concurrent downloads
	// Create a mutex to protect the iconDefs slice
	var mu sync.Mutex

	// Filter SVG files
	var svgFiles []GitHubContent
	for _, file := range files {
		if filepath.Ext(file.Name) == ".svg" {
			svgFiles = append(svgFiles, file)
		}
	}

	// Setup progress tracking
	total := len(svgFiles)
	fmt.Printf("Found %d SVG icons to process\n", total)

	var processed int32
	statusChan := make(chan bool)
	doneChan := make(chan bool)

	// Start progress reporter goroutine
	go func() {
		for {
			select {
			case <-statusChan:
				current := atomic.AddInt32(&processed, 1)
				percent := float64(current) / float64(total) * 100
				fmt.Printf("\rProcessing icons... [%d/%d] %.1f%% complete", current, total, percent)
			case <-doneChan:
				fmt.Println("\rProcessing icons... [100%] Complete!                ")
				return
			}
		}
	}()

	// Process each SVG file
	for _, file := range svgFiles {
		wg.Add(1)
		go func(file GitHubContent) {
			defer wg.Done()
			// Acquire a token from the semaphore
			semaphore <- struct{}{}
			defer func() {
				// Release the token back to the semaphore when done
				<-semaphore
			}()

			name := strings.TrimSuffix(file.Name, ".svg")
			funcName := toPascalCase(name)

			// Add icon definition (thread-safe)
			mu.Lock()
			iconDefs = append(iconDefs, fmt.Sprintf("var %s = Icon(%q)\n", funcName, name))
			mu.Unlock()

			// Download icon content
			content, err := downloadFile(file.DownloadUrl)
			if err != nil {
				fmt.Printf("\nError downloading %s: %v\n", file.Name, err)
				statusChan <- true
				return
			}

			// Save icon content to a separate file
			err = os.WriteFile(filepath.Join(iconContentDir, name+".svg"), content, 0644)
			if err != nil {
				fmt.Printf("\nError writing %s: %v\n", file.Name, err)
			}

			// Update status
			statusChan <- true
		}(file)
	}

	// Wait for all downloads to complete
	wg.Wait()

	// Signal the progress reporter to stop
	doneChan <- true

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

// fetchGitHubContents fetches the list of files from a GitHub repository directory
func fetchGitHubContents(url string) ([]GitHubContent, error) {
	fmt.Println("Fetching icons from GitHub repository...")

	// Instead of using the API, let's use a direct approach to get all icons
	// The API has pagination limits, but we can get file listings from the tree endpoint
	treeUrl := "https://api.github.com/repos/lucide-icons/lucide/git/trees/main?recursive=1"

	req, err := http.NewRequest("GET", treeUrl, nil)
	if err != nil {
		return nil, err
	}

	// Add GitHub API headers
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the tree response
	type TreeItem struct {
		Path string `json:"path"`
		Type string `json:"type"`
		Url  string `json:"url"`
	}

	type TreeResponse struct {
		Tree []TreeItem `json:"tree"`
	}

	var treeResp TreeResponse
	if err := json.Unmarshal(body, &treeResp); err != nil {
		return nil, err
	}

	// Filter for SVG files in the icons directory
	var contents []GitHubContent
	for _, item := range treeResp.Tree {
		if strings.HasPrefix(item.Path, "icons/") && strings.HasSuffix(item.Path, ".svg") {
			// Extract the filename from the path
			filename := filepath.Base(item.Path)

			// Create a GitHubContent entry
			content := GitHubContent{
				Name:        filename,
				Path:        item.Path,
				DownloadUrl: fmt.Sprintf("https://raw.githubusercontent.com/lucide-icons/lucide/main/%s", item.Path),
			}
			contents = append(contents, content)
		}
	}

	if len(contents) == 0 {
		return nil, fmt.Errorf("no SVG files found in the icons directory")
	}

	fmt.Printf("Found %d icons in the repository\n", len(contents))
	return contents, nil
}

// downloadFile downloads a file from a URL
func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download returned status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
