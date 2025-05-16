package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"ludo_backend/app/handlers/websocket"

	"ludo_backend/database"
)

func main() {
	database.InitMongoDB("mongodb://localhost:27017")
	// db := database.MongoClient.Database("ludo")
	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		websocket.InitWebsockets(c.Writer, c.Request)
	})

	fmt.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
