package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func handleFileDownload(w http.ResponseWriter, r *http.Request) {
	fileName := filepath.Base(r.URL.Path)
	filePath := filepath.Join(storageDir, fileName)

	// Prevent path traversal attacks
	absDir, _ := filepath.Abs(storageDir)
	absFile, _ := filepath.Abs(filePath)
	rel, err := filepath.Rel(absDir, absFile)
	if err != nil || strings.HasPrefix(rel, "..") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found: "+fileName, http.StatusNotFound)
		return
	}
	defer file.Close()

	// Detect content type from the first 512 bytes
	buffer := make([]byte, 512)
	n, _ := file.Read(buffer)
	contentType := http.DetectContentType(buffer[:n])

	// Seek back to the beginning
	file.Seek(0, 0)

	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Header().Set("Content-Type", contentType)
	io.Copy(w, file)
}
