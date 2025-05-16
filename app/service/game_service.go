package service

import (
	"fmt"
	"ludo_backend/app/repository"
	models "ludo_backend/models/game_models"
	helpers "ludo_backend/utils/helpers"
)

type GameService struct {
	GameRepo *repository.GameRepository
}

func NewGameService(gameRepo *repository.GameRepository) *GameService {
	return &GameService{
		GameRepo: gameRepo,
	}
}

func (service GameService) CreateGame(roomId string, creator string, uids []string) (models.Game, error) {
	var game models.Game
	game.GameID = roomId + "_" + creator
	game.RoomId = roomId
	game.CurrentPlayer = 0
	game.DiceResult = 3
	var players []models.Player
	for index, uid := range uids {
		player := models.Player{
			Uid:      uid,
			PlayerId: index,
			Color:    helpers.GetColor(index),
			Pawns: []models.Pawn{
				{Id: 0, Color: helpers.GetColor(index), Position: 0},
				{Id: 1, Color: helpers.GetColor(index), Position: 0},
				{Id: 2, Color: helpers.GetColor(index), Position: 0},
				{Id: 3, Color: helpers.GetColor(index), Position: 0},
			},
		}
		players = append(players, player)
	}
	game.Board = models.Board{Players: players}
	err := service.GameRepo.CreateGame(game)
	return game, err
}

func (service GameService) GetGameById(gameId string) (models.Game, error) {
	game, err := service.GameRepo.GetGameById(gameId)
	if err != nil {
		return models.Game{}, err
	}
	return game, nil
}


func (service GameService) HandlePawnMovement(pawnMovement models.PawnMovementRequest)  (models.PawnMovementResponse, error) {
	game, err := service.GameRepo.GetGameById(pawnMovement.GameId)
	var pawnMovementResponse models.PawnMovementResponse
	if err != nil {
		return pawnMovementResponse, err
	}
	
	if(game.CurrentPlayer != pawnMovement.PlayerId){
		return pawnMovementResponse, fmt.Errorf("it's not your turn")
	}
	currentPawnPosition := game.Board.Players[pawnMovement.PlayerId].Pawns[pawnMovement.PawnId].Position
	err = service.GameRepo.MovePawm(pawnMovement, currentPawnPosition + game.DiceResult)
	if err != nil {
		return pawnMovementResponse, err
	}
	pawnMovementResponse.GameId = pawnMovement.GameId
	pawnMovementResponse.PlayerId = pawnMovement.PlayerId
	pawnMovementResponse.PawnId = pawnMovement.PawnId
	pawnMovementResponse.Position = currentPawnPosition + game.DiceResult
	return pawnMovementResponse, nil
}