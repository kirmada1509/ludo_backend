package websocket

import (
	"log"
	game_constants "ludo_backend/utils/constants"
	services "ludo_backend/app/service"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleGameWebSocket(w http.ResponseWriter, r *http.Request, client *Client) {
	gameService := &services.GameService{}
	var room *Room
	if client_rooms[client.ID] != "" {
		room = Rooms[client_rooms[client.ID]]
	} else {
		room = getAvailableRoom(Rooms, client)
		client_rooms[client.ID] = room.ID
		room.Clients[client.ID] = client
	}

	client.Conn.WriteJSON(map[string]interface{}{
		"message": "Welcome to the game room!",
		"roomId":  room.ID,
		"users":   len(room.Clients),
	})
	var userIds []string
	for id := range room.Clients {
		userIds = append(userIds, id)
	}
	if len(room.Clients) == game_constants.MaxPlayers {
		game := gameService.CreateGame(room.ID, client.ID, userIds)
		room.Game = game
		for _, c := range room.Clients {
			c.Conn.WriteJSON(map[string]interface{}{
				"message": "Game is starting!",
				"roomId":  room.ID,
			})
		}

		for _, c := range room.Clients {
			c.Conn.WriteJSON(map[string]interface{}{
				"message": "Game Started!",
				"roomId":  room.ID,
				"game":    game,
				"users":   len(room.Clients),
			})
		}

	} else {
		client.Conn.WriteJSON(map[string]interface{}{
			"message": "Waiting for more players to join...",
			"roomId":  room.ID,
			"users":   len(room.Clients),
		})
	}

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Error reading message"+err.Error()))
			break
		}

		for _, c := range room.Clients {
			c.Conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
