package main

import (
	"fmt"
	"log"
	"net/http"

	"ludo_backend/ws"
)

func main() {
	http.HandleFunc("/ws", ws.HandleGameWebSocket)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServve:", err)
	}
}
