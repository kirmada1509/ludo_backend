package models

type Game struct {
	GameID        string `json:"game_id" bson:"game_id"`
	RoomUd        string `json:"room_id" bson:"room_id"`
	CurrentPlayer int    `json:"current_player_id" bson:"current_player_id"` //r g b y
	Board         Board  `json:"board" bson:"board"`
}

type Board struct {
	Players []Player `json:"players" bson:"players"`
}

type Player struct {
	Uid      string `json:"uid" bson:"uid"`
	PlayerId int    `json:"player_id" bson:"player_id"` //r g b y
	Color    string `json:"color" bson:"color"`
	Pawns    []Pawn `json:"pawns" bson:"pawns"`
}

type Pawn struct {
	Id       int    `json:"id" bson:"id"`
	Color    string `json:"color" bson:"color"` // r g b y
	Position int    `json:"position" bson:"position"`
}
