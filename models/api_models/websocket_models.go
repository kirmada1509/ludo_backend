package models

import "github.com/gorilla/websocket"

type WebSocketClient struct {
	ID   string
	Conn *websocket.Conn
}

type WebSocketRoom struct {
	ID      int
	Clients []*WebSocketClient
}