package middleware

import (
	"context"
	"net/http"
	"os"
	"os/exec"

	"github.com/axzilla/templui/internal/ctxkeys"
)

func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0, must-revalidate, no-cache, no-store, private")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// WithURLPathValue adds the current URL's path to the context.
func WithURLPathValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(
			r.Context(),
			ctxkeys.URLPathValue,
			r.URL.Path,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LatestVersion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("version.txt")
		if err != nil {
			// If version.txt doesn't exist, fall back to git commit hash
			commit, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
			if err != nil {
				http.Error(w, "Unable to determine version", http.StatusInternalServerError)
				return
			}
			data = commit
		}
		ctx := context.WithValue(r.Context(), ctxkeys.Version, string(data))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// --- Logging Middleware Components ---

// loggingResponseWriter wraps http.ResponseWriter to capture status code.
type loggingResponseWriter struct {
	http.ResponseWriter     // Embed the original ResponseWriter
	statusCode          int // Store the status code
}

// NewLoggingResponseWriter creates a new loggingResponseWriter.
func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// Default to 200 OK if WriteHeader is never called.
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

// WriteHeader captures the status code before calling the original WriteHeader.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs incoming request details and response status/duration.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start := time.Now()
		lrw := NewLoggingResponseWriter(w)
		// log.Printf("INFO: --> %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(lrw, r)
		// duration := time.Since(start)
		// statusCode := lrw.statusCode
		// log.Printf("INFO: <-- %s %s completed in %v (Status: %d)", r.Method, r.URL.Path, duration, statusCode)
	})
}
