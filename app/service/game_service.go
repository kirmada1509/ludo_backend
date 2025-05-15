package service

import (
	models "ludo_backend/models/game_models"
)

type GameService struct {
}

func (gameService GameService) CreateGame(roomId string, creator string, uids []string) models.Game {
	var game models.Game
	game.GameID = roomId + "_"+ creator
	game.RoomUd = roomId
	game.CurrentPlayer = 0
	var players []models.Player
	for index, uid := range uids {
		player := models.Player{
			Uid: uid,
			PlayerId: index,
			Color: getColor(index),
			Pawns: []models.Pawn{
				{Id: 0, Color: getColor(index), Position: -1},
				{Id: 1, Color: getColor(index), Position: -1},
				{Id: 2, Color: getColor(index), Position: -1},
				{Id: 3, Color: getColor(index), Position: -1},
			},
		}
		players = append(players, player)
	}
	game.Board = models.Board{Players: players}
	return game
}


func getColor(playerId int) string {
	switch playerId {
	case 0:
		return "red"
	case 1:
		return "green"
	case 2:
		return "blue"
	case 3:
		return "yellow"
	default:
		return ""
	}
}