package game_handlers

import (
	models "ludo_backend/models/game_models"
	"ludo_backend/app/service"
)

type GameHandler struct {
	Game models.Game
	GameService service.GameService
}