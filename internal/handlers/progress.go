package handlers

import (
	"net/http"
	"sync"

	"github.com/axzilla/templui/components"
)

// Simple in-memory progress storage
var (
	progressValue      = 0
	progressValueMutex sync.RWMutex
)

// ProgressHandler returns the current progress
func ProgressHandler(w http.ResponseWriter, r *http.Request) {
	// Get current progress
	progressValueMutex.RLock()
	currentProgress := progressValue
	progressValueMutex.RUnlock()

	// Increment progress for next request (simulate progress)
	progressValueMutex.Lock()
	if progressValue < 100 {
		progressValue += 10
		if progressValue > 100 {
			progressValue = 100
		}
	}
	progressValueMutex.Unlock()

	// Render the progress component
	components.Progress(components.ProgressProps{
		Value:     currentProgress,
		Label:     "Processing data...",
		ShowValue: true,
		HxGet:     "/api/progress",
		HxTrigger: "every 2s",
		HxTarget:  "#progress-container",
	}).Render(r.Context(), w)
}

// ProgressResetHandler resets the progress
func ProgressResetHandler(w http.ResponseWriter, r *http.Request) {
	// Reset progress
	progressValueMutex.Lock()
	progressValue = 0
	progressValueMutex.Unlock()

	// Render the reset progress component
	components.Progress(components.ProgressProps{
		Value:     0,
		Label:     "Processing data...",
		ShowValue: true,
		HxGet:     "/api/progress",
		HxTrigger: "every 2s",
		HxTarget:  "#progress-container",
	}).Render(r.Context(), w)
}
