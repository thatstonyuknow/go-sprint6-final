package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

// Server encapsulates the HTTP server and logger.
type Server struct {
	Logger     *log.Logger
	HTTPServer *http.Server
}

// loggingMiddleware logs each incoming HTTP request.
func loggingMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Request: %s %s from %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// MyServer creates a new server with the given logger and registers the handlers.
func MyServer(logger *log.Logger) *Server {
	router := http.NewServeMux()

	// Register the root handler to serve index.html.
	router.HandleFunc("/", handlers.RootHandler)
	// Register the upload handler to process file uploads.
	router.HandleFunc("/upload", handlers.UploadHandler)

	// Wrap the router with the logging middleware.
	loggedRouter := loggingMiddleware(logger, router)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggedRouter,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Logger:     logger,
		HTTPServer: server,
	}
}
