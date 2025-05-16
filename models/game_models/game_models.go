package models
type Game struct {
	GameID        string `json:"game_id" bson:"game_id"`
	RoomId        string `json:"room_id" bson:"room_id"`
	CurrentPlayer int    `json:"current_player_id" bson:"current_player_id"` //r g b y
	DiceResult    int    `json:"dice_result" bson:"dice_result"`
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

type PawnMovementRequest struct {
	GameId string `json:"game_id" bson:"game_id"`
	PlayerId int `json:"player_id" bson:"player_id"`
	PawnId int    `json:"pawn_id" bson:"pawn_id"`
}
type PawnMovementResponse struct {
	GameId string `json:"game_id" bson:"game_id"`
	PlayerId int `json:"player_id" bson:"player_id"`
	PawnId int    `json:"pawn_id" bson:"pawn_id"`
	Position int    `json:"position" bson:"position"`
}
