package service

import (
	"fmt"
	"ludo_backend/app/repository"
	models "ludo_backend/models/game_models"
	helpers "ludo_backend/utils/helpers"

	"github.com/google/uuid"
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
	game.GameID = uuid.New().String()
	game.RoomId = roomId
	game.CurrentPlayer = 0
	game.DiceResult = 0
	var players []models.Player
	for index, uid := range uids {
		player := models.Player{
			Uid:      uid,
			PlayerId: index,
			Color:    helpers.GetColor(index),
			Pawns: []models.Pawn{
				{PawnId: uid + "0", Index: 0, Color: helpers.GetColor(index), HomePathIndex: helpers.GetHomePosition(index), CurrentPathIndex: helpers.GetHomePosition((index)), IsInHome: true},
				{PawnId: uid + "1", Index: 1, Color: helpers.GetColor(index), HomePathIndex: helpers.GetHomePosition(index), CurrentPathIndex: helpers.GetHomePosition((index)), IsInHome: true},
				{PawnId: uid + "2", Index: 2, Color: helpers.GetColor(index), HomePathIndex: helpers.GetHomePosition(index), CurrentPathIndex: helpers.GetHomePosition((index)), IsInHome: true},
				{PawnId: uid + "3", Index: 3, Color: helpers.GetColor(index), HomePathIndex: helpers.GetHomePosition(index), CurrentPathIndex: helpers.GetHomePosition((index)), IsInHome: true},
			},
		}
		players = append(players, player)
	}
	game.Board = models.Board{Players: players}
	game.AllowMovement = false
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
	if game.AllowMovement {
		return 0, fmt.Errorf("you have already rolled the dice")
	}
	if game.CurrentPlayer != diceRoll.PlayerId {
		return 0, fmt.Errorf("it's not your turn")
	}
	diceResult := helpers.RollDice()
	game.DiceResult = diceResult
	game.AllowMovement = true
	game, err = service.GameRepo.UpdateGame(game)
	if err != nil {
		return 0, err
	}
	return diceResult, nil
}

func (service GameService) HandlePawnMovement(req models.PawnMovementRequest) (models.PawnMovementResponse, error) {
	game, err := service.GameRepo.GetGameById(req.GameId)
	var res models.PawnMovementResponse
	if err != nil {
		return res, err
	}

	if game.CurrentPlayer != req.PlayerId {
		return res, fmt.Errorf("it's not your turn, current player is %d", game.CurrentPlayer)
	}

	pawn := &game.Board.Players[req.PlayerId].Pawns[req.PawnIndex]
	pawn.CurrentPathIndex += game.DiceResult

	game.AllowMovement = false
	game.CurrentPlayer = (game.CurrentPlayer + 1) % len(game.Board.Players)

	_, err = service.GameRepo.UpdateGame(game)
	if err != nil {
		return res, err
	}

	fmt.Printf("Moved pawn %d of player %d to position %d\n", req.PawnIndex, req.PlayerId, pawn.CurrentPathIndex)

	res = models.PawnMovementResponse{
		GameId:        req.GameId,
		PlayerId:      req.PlayerId,
		CurrentPlayer: game.CurrentPlayer,
		PawnIndex:     pawn.Index,
		PathIndex:     pawn.CurrentPathIndex,
		PawnId:        pawn.PawnId,
	}
	return res, nil
}
