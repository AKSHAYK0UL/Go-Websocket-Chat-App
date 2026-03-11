package main

import (
	"log"
	"net/http"
)

func main() {
	rm := &RoomManager{rooms: make(map[string]*Room)}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ws", rm.HandleWebSocket)
	mux.HandleFunc("GET /", handleHome)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
