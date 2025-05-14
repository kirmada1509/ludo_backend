package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type Room struct {
	ID      int
	Clients []*Client
}

var clients = make(map[string]*Client)
var Rooms = []*Room{}
var client_rooms = make(map[string]int)

func HandleGameWebSocket(w http.ResponseWriter, r *http.Request) {

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

	var room *Room
	if client_rooms[userId] == 0 {
		room = getAvailableRoom()
		if room == nil {
			room = &Room{
				ID:      len(Rooms) + 1,
				Clients: []*Client{},
			}
			Rooms = append(Rooms, room)
		}
		room.Clients = append(room.Clients, client)
	}else {
		room = Rooms[client_rooms[userId]]
	}

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

func getAvailableRoom() *Room {
	for _, room := range Rooms {
		if len(room.Clients) < 4 {
			return room
		}
	}
	return nil
}