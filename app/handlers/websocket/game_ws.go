package websocket

import (
	"log"
	models "ludo_backend/models/game_models"
	game_constants "ludo_backend/utils/constants"
)


func (handler WebsocketHandler) HandleRoomJoin(client *Client) {
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

		game, err := handler.GameService.CreateGame(room.ID, client.ID, userIds)
		if err != nil {
			for _, c := range room.Clients {
				c.Conn.WriteJSON(map[string]interface{}{
					"message": "Error creating game: " + err.Error(),
					"roomId":  room.ID,
				})
				return
			}
		}

		room.Game = game
		game_rooms[game.GameID] = room.ID
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
}


func (handler WebsocketHandler) HandlePawnMovement(pawnMovement models.PawnMovementRequest) {
	pawnMovementResponse, err := handler.GameService.HandlePawnMovement(pawnMovement)
	if err != nil {
		log.Println("Error handling pawn movement:", err)
	}
	game, err := handler.GameRepo.GetGameById(pawnMovement.GameId)
	if err != nil {
		log.Println("Error getting game by ID:", err)
		return
	}
	
	for _, c := range Rooms[game_rooms[pawnMovement.GameId]].Clients {
		c.Conn.WriteJSON(map[string]interface{}{
			"event": "pawn_movement",
			"movement": pawnMovementResponse,
			"board": game.Board,
		})
	}
}