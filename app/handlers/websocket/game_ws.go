package websocket

import (
	"fmt"
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
					"event": "error",
					"payload": map[string]interface{}{
						"message": "Error creating game: " + err.Error(),
						"roomId":  room.ID,
					},
				})
			}
			return
		}

		room.Game = game
		game_rooms[game.GameID] = room.ID

		for i, c := range room.Clients {
			c.Conn.WriteJSON(map[string]interface{}{
				"event": "game_started",
				"payload": map[string]interface{}{
					"roomId": room.ID,
					"game":   game,
					"player_id": i,
				},
			})
		}
	} else {
		// playerNames := make([]string, 0)
		// for id := range room.Clients {
		// 	playerNames = append(playerNames, id)
		// }
		client.Conn.WriteJSON(map[string]interface{}{
			"event": "waiting",
			"payload": map[string]interface{}{
				"roomId":  room.ID,
				"player_id": len(room.Clients) - 1,
			},
		})
	}
}

func (handler WebsocketHandler) HandleDiceRoll(req models.DiceRollRequest) {
	diceResult, err := handler.GameService.HandleDiceRoll(req)
	if err != nil {
		clients[req.UserId].Conn.WriteJSON(map[string]interface{}{
			"event": "dice_roll",
			"payload": map[string]interface{}{
				"success": false,
				"message": "Error handling dice roll: " + err.Error(),
			},
		})
		return
	}
	fmt.Println("Dice rolled:", diceResult)
	for _, c := range Rooms[game_rooms[req.GameId]].Clients {
		c.Conn.WriteJSON(map[string]interface{}{
			"event": "dice_roll",
			"payload": map[string]interface{}{
				"success":     true,
				"dice_result": diceResult,
				"current_player": req.PlayerId,
			},
		})
	}
}

func (handler WebsocketHandler) HandlePawnMovement(req models.PawnMovementRequest) {
	resp, err := handler.GameService.HandlePawnMovement(req)
	if err != nil {
		log.Println("Error handling pawn movement:", err)
		clients[req.UserId].Conn.WriteJSON(map[string]interface{}{
			"event": "pawn_movement",
			"payload": map[string]interface{}{
				"success": false,
				"message": "Error handling pawn movement: " + err.Error(),
			},
		})
		return
	}

	for _, c := range Rooms[game_rooms[req.GameId]].Clients {
		c.Conn.WriteJSON(map[string]interface{}{
			"event": "pawn_movement",
			"payload": map[string]interface{}{
				"movement": resp,
			},
		})
	}
}
