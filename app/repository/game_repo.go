package repository

import (
	"context"
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


