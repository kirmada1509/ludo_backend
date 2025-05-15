package websocket

import (
	models "ludo_backend/models/game_models"
	"net/http"

	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type Room struct {
	ID      string
	Clients map[string]*Client
	Game    models.Game
}

var clients = make(map[string]*Client)
var Rooms = make(map[string]*Room)
var client_rooms = make(map[string]string) // clientId : roomId

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebsockets(w http.ResponseWriter, r *http.Request) {
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
		conn.Close()
		delete(Rooms[client_rooms[userId]].Clients, userId)
		delete(client_rooms, userId)
		delete(clients, userId)

		for _, c := range Rooms[client_rooms[userId]].Clients {
			c.Conn.WriteJSON(map[string]interface{}{
				"message": "User disconnected",
				"userId":  userId,
				"users":   len(Rooms[client_rooms[userId]].Clients),
			})
		}
	}()

	// Create a new client
	client := &Client{
		ID:   userId,
		Conn: conn,
	}
	clients[userId] = client
	HandleGameWebSocket(w, r, client)
}

func getAvailableRoom(Rooms map[string]*Room, client *Client) *Room {
	for _, room := range Rooms {
		if len(room.Clients) < 4 {
			return room
		}
	}
	room := Room{
		ID:      fmt.Sprint(len(Rooms) + 1),
		Clients: map[string]*Client{},
	}
	Rooms[room.ID] = &room
	client_rooms[client.ID] = room.ID
	room.Clients[client.ID] = client
	return &room
}
