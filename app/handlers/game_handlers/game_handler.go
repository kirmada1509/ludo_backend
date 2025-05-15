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
