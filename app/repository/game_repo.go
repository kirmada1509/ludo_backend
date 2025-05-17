package repository

import (
	"context"
	"fmt"
	models "ludo_backend/models/game_models"

	"go.mongodb.org/mongo-driver/mongo"
)

type GameRepository struct {
	db *mongo.Database
}

func NewGameRepository(db *mongo.Database) *GameRepository {
	return &GameRepository{
		db: db,
	}
}

func (repo GameRepository) CreateGame(game models.Game) error {
	// Create a new game in the database
	_, err := repo.db.Collection("games").InsertOne(context.TODO(), game)
	if err != nil {
		return nil
	}
	return err
}

func (repo GameRepository) GetGameById(gameId string) (models.Game, error) {
	// Get a game by its ID from the database
	var game models.Game
	err := repo.db.Collection("games").FindOne(context.TODO(), map[string]interface{}{"game_id": gameId}).Decode(&game)
	if err != nil {
		return models.Game{}, err
	}
	return game, nil
}
//TODO: Remove this and use UpdateGame()
func (repo GameRepository) MovePawm(pawnMovement models.PawnMovementRequest, diceResult int) error {
	filter := map[string]interface{}{"game_id": pawnMovement.GameId}
	update := map[string]interface{}{
	"$set": map[string]interface{}{
		fmt.Sprintf("board.players.%d.pawns.%d.position", pawnMovement.PlayerId, pawnMovement.PawnId): diceResult,
	}}

	_, err := repo.db.Collection("games").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo GameRepository) GetPawn(gameId string, playerId int, pawnId int) (models.Pawn, error) {
	var game models.Game
	err := repo.db.Collection("games").FindOne(context.TODO(), map[string]interface{}{"game_id": gameId}).Decode(&game)
	if err != nil {
		return models.Pawn{}, err
	}
	return game.Board.Players[playerId].Pawns[pawnId], nil
}

func (repo GameRepository) UpdateGame(game models.Game) (models.Game, error) {
	// Update the game in the database
	var updatedGame models.Game
	result := repo.db.Collection("games").FindOneAndUpdate(
		context.TODO(),
		map[string]interface{}{"game_id": game.GameID},
		map[string]interface{}{"$set": game},
	)
	err := result.Decode(&updatedGame)
	if err != nil {
		return updatedGame, err
	}
	return updatedGame, nil
}