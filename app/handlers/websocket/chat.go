package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)



func HandleChatWebsocket(w http.ResponseWriter, r *http.Request) {

	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	// Upgrade the connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	defer func() {
		delete(clients, userId)
		conn.Close()
		log.Printf("Client disconnected: %s", userId)
	}()

	// Create a new client
	client := &Client{
		ID:   userId,
		Conn: conn,
	}
	clients[userId] = client

	room := getAvailableRoom(Rooms, client)

	client.Conn.WriteJSON(map[string]interface{}{
		"roomId": room.ID,
		"users":  (room.Clients),
	})

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Error reading message" + err.Error()))
			break
		}
		
		for _, c := range room.Clients {
			c.Conn.WriteMessage(websocket.TextMessage, msg)
		}

	}

}