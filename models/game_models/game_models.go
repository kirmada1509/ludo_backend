package models

type Game struct {
	GameID        string    `json:"game_id"`
	RoomUd        string `json:"room_id"`
	CurrentPlayer int    `json:"current_player_id"` //r g b y
	Board         Board  `json:"board"`
}

type Board struct {
	Players []Player `json:"players"`
}

type Player struct {
	Uid      string `json:"uid"`
	PlayerId int    `json:"player_id"` //r g b y
	Color    string `json:"color"`
	Pawns    []Pawn `json:"pawns"`
}

type Pawn struct {
	Id       int    `json:"id"`
	Color    string `json:"color"` // r g b y
	Position int    `json:"position"`
}
