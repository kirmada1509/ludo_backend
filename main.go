package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"ludo_backend/app/handlers/game_handlers"
	"ludo_backend/app/handlers/websocket"
	"ludo_backend/db"
)

func main() {
	db.InitMongoDB("mongodb://localhost:27017")
	
	r := gin.Default()

	gameHandler := &game_handlers.GameHandler{}
	r.GET("/ws", func(c *gin.Context) {
		websocket.HandleWebsockets(c.Writer, c.Request)
	})

	r.POST("/game/roll_dice", gameHandler.RollDice)

	fmt.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
