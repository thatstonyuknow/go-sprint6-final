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

// RootHandler serves the HTML file "index.html" for the root endpoint.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Return 404 if the URL path is not exactly "/".
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// UploadHandler processes the file upload.
// It parses the multipart form data, retrieves the uploaded file,
// writes its content to a temporary file, and then reads the data
// using os.ReadFile. The content is passed to the auto-detection function
// in the service package to obtain the converted string.
// A local file is then created (using the current UTC time and the uploaded file's extension)
// to save the conversion result, which is then returned.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Return 404 if the URL path is not exactly "/".
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
	// Clean up the temporary file when we're done.
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

	// Create a local file to save the conversion result.
	localFile, err := os.Create(filename)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	// Write the conversion result to the local file.
	if _, err = localFile.WriteString(converted); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return the conversion result in the HTTP response.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(converted))
}
