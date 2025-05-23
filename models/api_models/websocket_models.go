package models

type WebSocketResponse struct {
	Success bool        `json:"success"`
	Event   string      `json:"event"`
	Payload interface{} `json:"payload,omitempty"`
	Error   *WSError    `json:"error,omitempty"`
}

type WSError struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}
