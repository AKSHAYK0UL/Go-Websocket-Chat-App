package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	os.Mkdir(storageDir, os.ModePerm)
	rm := &RoomManager{rooms: make(map[string]*Room)}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ws", rm.HandleWebSocket)
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("POST /upload", HandleFileUpload)
	mux.HandleFunc("GET /download/", handleFileDownload)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
