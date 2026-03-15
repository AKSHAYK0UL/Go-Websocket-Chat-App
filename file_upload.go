package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Max upload size: 100 MB
const maxUploadSize = 100 << 20

// Blocked file extensions for security
var blockedExtensions = map[string]bool{
	".exe": true, ".sh": true, ".bat": true,
	".cmd": true, ".ps1": true, ".msi": true,
}

func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File too large (max 100MB)", http.StatusRequestEntityTooLarge)
		return
	}

	file, metadata, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Block dangerous file types
	ext := strings.ToLower(filepath.Ext(metadata.Filename))
	if blockedExtensions[ext] {
		http.Error(w, "File type not allowed", http.StatusUnsupportedMediaType)
		return
	}

	// Prefix with timestamp to avoid name collisions
	safeName := fmt.Sprintf("%d_%s", time.Now().UnixMilli(), filepath.Base(metadata.Filename))
	destPath := filepath.Join(storageDir, safeName)

	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, file); err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, safeName)
}
