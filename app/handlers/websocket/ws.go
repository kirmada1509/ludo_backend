package websocket

import (
	"encoding/json"
	"log"
	"ludo_backend/app/repository"
	services "ludo_backend/app/service"
	"ludo_backend/database"
	models "ludo_backend/models/game_models"
	game_constants "ludo_backend/utils/constants"
	"net/http"

	"fmt"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

type WSMessage struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type Room struct {
	ID      string
	Clients map[string]*Client
	Game    models.Game
}

var clients = make(map[string]*Client)
var Rooms = make(map[string]*Room)
var client_rooms = make(map[string]string) // clientId : roomId
var game_rooms = make(map[string]string)   // gameId : roomId

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketHandler struct {
	Db          *mongo.Database
	GameRepo    *repository.GameRepository
	GameService *services.GameService
}

func NewWebsocketHandler(db *mongo.Database) *WebsocketHandler {
	return &WebsocketHandler{
		Db:          db,
		GameRepo:    repository.NewGameRepository(db),
		GameService: services.NewGameService(repository.NewGameRepository(db)),
	}
}

func InitWebsockets(w http.ResponseWriter, r *http.Request) {
	WsHandler := NewWebsocketHandler(database.MongoClient.Database("ludo"))
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	// Upgrade the connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	defer func() {
		conn.Close()
		delete(Rooms[client_rooms[userId]].Clients, userId)
		delete(client_rooms, userId)
		delete(clients, userId)

		for _, c := range Rooms[client_rooms[userId]].Clients {
			c.Conn.WriteJSON(map[string]interface{}{
				"message": "User disconnected",
				"userId":  userId,
				"users":   len(Rooms[client_rooms[userId]].Clients),
			})
		}
	}()

	client := &Client{
		ID:   userId,
		Conn: conn,
	}
	clients[userId] = client

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage error:", err)
			break
		}
		
		var msg WSMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("Unmarshal error:", err)
			conn.WriteJSON(map[string]string{"error": "invalid message format"})
			continue
		}

		switch msg.Action {
		case "join":
			WsHandler.HandleRoomJoin(client)
		case "dice_roll":
			var DiceRollRequest models.DiceRollRequest
			if err := json.Unmarshal([]byte(msg.Payload), &DiceRollRequest); err != nil {
				conn.WriteJSON(map[string]interface{}{"error": "Invalid payload"})
			}
			WsHandler.HandleDiceRoll(DiceRollRequest)
		case "move":
			var pawnMovement models.PawnMovementRequest
			if err := json.Unmarshal([]byte(msg.Payload), &pawnMovement); err != nil {
				conn.WriteJSON(map[string]interface{}{"error": "Invalid payload"})
				continue
			}
			WsHandler.HandlePawnMovement(pawnMovement)

		default:
			log.Println("Unknown action:", msg.Action)
			conn.WriteJSON(map[string]interface{}{
				"message": fmt.Sprint("Unknown action", msg.Action),
			})
		}
	}
}

func getAvailableRoom(Rooms map[string]*Room, client *Client) *Room {
	for _, room := range Rooms {
		if len(room.Clients) < game_constants.MaxPlayers {
			return room
		}
	}
	room := Room{
		ID:      fmt.Sprint(len(Rooms) + 1),
		Clients: map[string]*Client{},
	}
	Rooms[room.ID] = &room
	client_rooms[client.ID] = room.ID
	room.Clients[client.ID] = client
	return &room
}
