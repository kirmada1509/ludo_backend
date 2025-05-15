package service

import (
	models "ludo_backend/models/game_models"
	"math/rand"
)

type GameService struct {
	Game models.Game
}

func (gameService *GameService) CreateGame(roomId string, uid string) models.Game {
	gameService.Game.GameID = roomId + "_"+ uid
	gameService.Game.RoomUd = roomId
	gameService.Game.CurrentPlayer = 0

	// Initialize players
	players := []models.Player{
		{Uid: uid, PlayerId: 0, Color: "red", Pawns: []models.Pawn{{Id: 1, Color: "red", Position: 0}}},
	}

	gameService.Game.Board = models.Board{Players: players}
	return gameService.Game
}

func (service GameService) DiceRoll() int {
	service.Game.CurrentPlayer = (service.Game.CurrentPlayer + 1) % 4
	return rand.Intn(6) + 1
}

func (service GameService) MovePawn(pawnPosition int, diceValue int) int {
	newPosition := pawnPosition + diceValue
	return newPosition
}

