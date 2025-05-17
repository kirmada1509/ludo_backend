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

func (service GameService) HandleDiceRoll(diceRoll models.DiceRollRequest) (int, error) {
	game, err := service.GameRepo.GetGameById(diceRoll.GameId)
	if err != nil {
		return 0, err
	}
	if game.CurrentPlayer != diceRoll.PlayerId {
		return 0, fmt.Errorf("it's not your turn")
	}
	game.DiceResult = helpers.RollDice()
	game, err = service.GameRepo.UpdateGame(game)
	if err != nil {
		return 0, err
	}
	return game.DiceResult, err
}

func (service GameService) HandlePawnMovement(req models.PawnMovementRequest) (models.PawnMovementResponse, error) {
	game, err := service.GameRepo.GetGameById(req.GameId)
	var res models.PawnMovementResponse
	if err != nil {
		return res, err
	}

	if game.CurrentPlayer != req.PlayerId {
		return res, fmt.Errorf("it's not your turn")
	}

	player := game.Board.Players[req.PlayerId]
	pawn := player.Pawns[req.PawnId]
	currentPawnPosition := pawn.Position
	newPosition := currentPawnPosition + game.DiceResult

	err = service.GameRepo.MovePawm(req, newPosition)
	if err != nil {
		return res, err
	}

	game.CurrentPlayer = (game.CurrentPlayer + 1) % len(game.Board.Players)
	game, err = service.GameRepo.UpdateGame(game)
	if err != nil {
		return res, err
	}
	res = models.PawnMovementResponse{
		GameId:   req.GameId,
		PlayerId: req.PlayerId,
		PawnId:   req.PawnId,
		Position: newPosition,
	}
	return res, nil
}
