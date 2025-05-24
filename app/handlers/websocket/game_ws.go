package websocket

import (
	"log"
	api_models "ludo_backend/models/api_models"
	game_models "ludo_backend/models/game_models"
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
			resp := api_models.WebSocketResponse{
				Success: false,
				Event:   "game_started",
				Error: &api_models.WSError{
					Message: "Error creating game: " + err.Error(),
				},
			}
			for _, c := range room.Clients {
				c.Conn.WriteJSON(resp)
			}
			return
		}

		room.Game = game
		game_rooms[game.GameID] = room.ID

		for _, c := range room.Clients {
			resp := api_models.WebSocketResponse{
				Success: true,
				Event:   "game_started",
				Payload: map[string]interface{}{
					"game":      game,
				},
			}
			c.Conn.WriteJSON(resp)
		}
	} else {
		resp := api_models.WebSocketResponse{
			Success: true,
			Event:   "waiting",
			Payload: map[string]interface{}{
				"room_id": room.ID,
				"player_id": len(room.Clients) - 1,
			},
		}
		client.Conn.WriteJSON(resp)
	}
}

func (handler WebsocketHandler) HandleDiceRoll(req game_models.DiceRollRequest) {
	diceResult, err := handler.GameService.HandleDiceRoll(req)
	if err != nil {
		resp := api_models.WebSocketResponse{
			Success: false,
			Event:   "dice_roll",
			Error: &api_models.WSError{
				Message: "Dice roll failed: " + err.Error(),
			},
		}
		clients[req.UserId].Conn.WriteJSON(resp)
		return
	}

	resp := api_models.WebSocketResponse{
		Success: true,
		Event:   "dice_roll",
		Payload: map[string]interface{}{
			"dice_result": diceResult,
		},
	}
	for _, c := range Rooms[game_rooms[req.GameId]].Clients {
		c.Conn.WriteJSON(resp)
	}
}

func (handler WebsocketHandler) HandlePawnMovement(req game_models.PawnMovementRequest) {
	movementResp, err := handler.GameService.HandlePawnMovement(req)
	if err != nil {
		log.Println("Error handling pawn movement:", err)
		resp := api_models.WebSocketResponse{
			Success: false,
			Event:   "pawn_movement",
			Error: &api_models.WSError{
				Message: "Pawn movement failed: " + err.Error(),
			},
		}
		clients[req.UserId].Conn.WriteJSON(resp)
		return
	}

	resp := api_models.WebSocketResponse{
		Success: true,
		Event:   "pawn_movement",
		Payload: map[string]interface{}{
			"movement": movementResp,
		},
	}

	for _, c := range Rooms[game_rooms[req.GameId]].Clients {
		c.Conn.WriteJSON(resp)
	}
}
