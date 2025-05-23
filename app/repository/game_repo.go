package repository

import (
	"context"
	"fmt"
	models "ludo_backend/models/game_models"

	"go.mongodb.org/mongo-driver/bson"
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

func (repo GameRepository) GetPawn(gameId string, playerId int, pawnId int) (models.Pawn, error) {
	var game models.Game
	err := repo.db.Collection("games").FindOne(context.TODO(), map[string]interface{}{"game_id": gameId}).Decode(&game)
	if err != nil {
		return models.Pawn{}, err
	}
	return game.Board.Players[playerId].Pawns[pawnId], nil
}

func (repo GameRepository) UpdateGame(game models.Game) (models.Game, error) {
	// Replace the entire game document in MongoDB
	filter := bson.M{"game_id": game.GameID}
	_, err := repo.db.Collection("games").ReplaceOne(context.TODO(), filter, game)
	if err != nil {
		return models.Game{}, err
	}
	updated, err := repo.GetGameById(game.GameID)
	if err != nil {
		return models.Game{}, err
	}
	fmt.Println("Updated game:", updated)
	return updated, nil
}
