package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

// RootHandler serves the HTML file "index.html"
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "./index.html")
}

// UploadHandler processes the file upload and writes the conversion result
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Return 404 if the URL path is not exactly "/upload".
	if r.URL.Path != "/upload" {
		http.NotFound(w, r)
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

	// Invoke file content into memory
	data, err := io.ReadAll(file)
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

	// Create a local file within the safe directory using the safe root.
	localFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
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

	if _, err := w.Write([]byte(converted)); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}
