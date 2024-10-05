// main.go
package main

import (
	"log"
	"net/http"
	"time"

	_ "live_editor/docs" // Import generated docs package
	"live_editor/handlers"
	"live_editor/storage"

	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/net/websocket"
)

// CustomResponseWriter wraps the http.ResponseWriter and captures the status code
type CustomResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader is used to capture the status code
func (rw *CustomResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// NewCustomResponseWriter initializes a new CustomResponseWriter
func NewCustomResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	// By default, we assume a 200 status code
	return &CustomResponseWriter{w, http.StatusOK}
}

// Logging Middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Use custom response writer to capture status code
		crw := NewCustomResponseWriter(w)

		var methodColors = map[string]string{
			"GET":     "\033[1;34m",       // Blue
			"POST":    "\033[1;32m",       // Green
			"PUT":     "\033[1;38;5;208m", // Orange
			"DELETE":  "\033[1;31m",       // Red
			"PATCH":   "\033[1;36m",       // Cyan
			"HEAD":    "\033[1;35m",       // Magenta
			"OPTIONS": "\033[1;37m",       // White
		}

		// Call the next handler
		next.ServeHTTP(crw, r)

		// Determine the log color based on the status code
		statusColor := "\033[32m" // Green for success
		if crw.statusCode >= 400 {
			statusColor = "\033[31m" // Red for errors (4xx, 5xx)
		}

		methodColor, ok := methodColors[r.Method]
		if !ok {
			methodColor = "\033[0m" // Default color
		}

		log.Printf("%s[%s]\033[0m %s %s%d %s\033[0m in %v",
			methodColor, r.Method,
			r.RequestURI,
			statusColor, crw.statusCode, http.StatusText(crw.statusCode), // Status code and text with color
			time.Since(start),
		)
	})
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

// @title Live Document Editor API
// @version 1.0
// @description This is a real-time document editor API with WebSockets.
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize in-memory storage and handlers
	store := storage.NewMemoryStorage()
	docHandler := handlers.NewDocumentHandler(store)
	wsHandler := handlers.NewWebSocketHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)

	// Define routes for document management
	mux.Handle("/documents", docHandler)

	// Define a WebSocket route
	mux.Handle("/ws", websocket.Handler(wsHandler.HandleConnection))

	// Serve the Swagger documentation
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Start the WebSocket broadcast in a separate goroutine
	go wsHandler.StartBroadcast()

	// Wrap the mux with a logger
	loggedMux := loggingMiddleware(mux)

	// Start the HTTP server
	log.Println("\033[1;33mStarting server on :8080...\033[0m")
	err := http.ListenAndServe(":8080", loggedMux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
