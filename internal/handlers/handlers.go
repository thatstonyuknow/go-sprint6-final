package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

const safeDir = "." // Use the project root as the safe directory

// RootHandler serves the HTML file "index.html" from the safe directory
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Return 404 if the URL path is not exactly "/".
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Open the safe directory (project root) as the root.
	root, err := os.OpenRoot(safeDir)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer root.Close()

	// Open index.html relative to the safe directory.
	file, err := root.Open("index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read the file contents.
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// UploadHandler processes the file upload and safely writes the conversion result
// into the safe directory using os.OpenRoot.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Return 404 if the URL path is not exactly "/upload".
	if r.URL.Path != "/upload" {
		http.NotFound(w, r)
		return
	}

	// Only accept POST requests.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form data with a 10 MB memory limit.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Retrieve the file from the form field "myFile".
	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create a temporary file to save the uploaded file's content.
	tmpFile, err := os.CreateTemp("", "upload-*")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Clean up the temporary file when done.
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Copy the uploaded file's content to the temporary file.
	if _, err := io.Copy(tmpFile, file); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Read the file content using os.ReadFile.
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert the file content using the auto-detection function from the service package.
	converted, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Generate a filename using the current UTC time and the extension of the uploaded file.
	ext := filepath.Ext(header.Filename)
	filename := strings.ReplaceAll(time.Now().UTC().String(), " ", "_")
	filename = strings.ReplaceAll(filename, ":", "-") + ext

	// Open the safe directory (project root) as the root.
	root, err := os.OpenRoot(safeDir)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer root.Close()

	// Create a local file within the safe directory using the safe root.
	localFile, err := root.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	// Write the conversion result to the local file.
	if _, err = localFile.Write([]byte(converted)); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return the conversion result in the HTTP response.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(converted))
}
