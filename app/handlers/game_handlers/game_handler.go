package game_handlers

import (
	models "ludo_backend/models/game_models"
	"ludo_backend/app/service"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	Game models.Game
	GameService service.GameService
}

func (handler *GameHandler) CreateGame(c *gin.Context) {

}

func (handler *GameHandler) RollDice(c *gin.Context) {
	//TODO: MIDDLEWARE FOR AUTHENTICATION BEFORE DICE ROLL
	type PlayerId struct {
		PlayerId int `json:"player_id"`
	}
	var playerId PlayerId
	if err := c.ShouldBindJSON(&playerId); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
	}
	if handler.Game.CurrentPlayer != playerId.PlayerId {
		c.JSON(403, gin.H{"error": "Not your turn"})
	}
	
	diceValue := handler.GameService.DiceRoll()
	c.JSON(200, gin.H{"dice_value": diceValue})

}
