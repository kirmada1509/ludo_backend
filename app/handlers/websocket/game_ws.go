package websocket

import (
	"log"
	models "ludo_backend/models/game_models"
	game_constants "ludo_backend/utils/constants"
)

func (handler WebsocketHandler) HandleRoomJoin(client *Client) {
	room := getAvailableRoom(Rooms, client)
	client_rooms[client.ID] = room.ID
	room.Clients[client.ID] = client

	if len(room.Clients) == game_constants.MaxPlayers {
		var userIds []string
		for id := range room.Clients {
			userIds = append(userIds, id)
		}
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
			})
		}

	} else {
		playerNames := make([]string, 0)
		for id := range room.Clients {
			playerNames = append(playerNames, id)
		}
		client.Conn.WriteJSON(map[string]interface{}{
			"message": "Waiting...",
			"roomId":  room.ID,
			"players": playerNames,
		})
	}
}

func (handler WebsocketHandler) HandleDiceRoll(DiceRollRequest models.DiceRollRequest) {
	diceResult, err := handler.GameService.HandleDiceRoll(DiceRollRequest)
	if err != nil {
		clients[DiceRollRequest.UserId].Conn.WriteJSON(map[string]interface{}{
			"success": false,
			"event":   "dice_roll",
			"message": "Error handling dice roll: " + err.Error(),
		})
		return
	}

	for _, c := range Rooms[game_rooms[DiceRollRequest.GameId]].Clients {
		c.Conn.WriteJSON(map[string]interface{}{
			"success":     true,
			"event":       "dice_roll",
			"dice_result": diceResult,
		})
	}
}

func (handler WebsocketHandler) HandlePawnMovement(pawnMovement models.PawnMovementRequest) {
	pawnMovementResponse, err := handler.GameService.HandlePawnMovement(pawnMovement)
	if err != nil {
		log.Println("Error handling pawn movement:", err)
	}

	for _, c := range Rooms[game_rooms[pawnMovement.GameId]].Clients {
		c.Conn.WriteJSON(map[string]interface{}{
			"event":    "pawn_movement",
			"movement": pawnMovementResponse,
		})
	}
}
