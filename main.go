package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"ludo_backend/app/handlers/websocket"

	"github.com/gin-gonic/gin"

	"ludo_backend/database"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	database.InitMongoDB("mongodb://localhost:27017")
	// db := database.MongoClient.Database("ludo")
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Ludo Backend API",
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		websocket.InitWebsockets(c.Writer, c.Request)
	})

	fmt.Println("Server running on http://localhost:8080")
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
