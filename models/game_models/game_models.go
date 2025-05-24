package models

type Game struct {
	GameID        string `json:"game_id" bson:"game_id"`
	RoomId        string `json:"room_id" bson:"room_id"`
	CurrentPlayer int    `json:"current_player_id" bson:"current_player_id"` //r g b y
	DiceResult    int    `json:"dice_result" bson:"dice_result"`
	Board         Board  `json:"board" bson:"board"`
	AllowMovement bool   `json:"allow_movement" bson:"allow_movement"`
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
	PawnId               string `json:"pawn_id" bson:"pawn_id"`
	Index            int    `json:"index" bson:"index"`
	Color            string `json:"color" bson:"color"` // r g b y
	HomePathIndex    int    `json:"home_path_index" bson:"home_path_index"`
	CurrentPathIndex int    `json:"current_path_index" bson:"current_path_index"`
	IsInHome         bool   `json:"is_in_home" bson:"is_in_home"`
}

type PawnMovementRequest struct {
	UserId    string `json:"user_id" bson:"user_id"`
	GameId    string `json:"game_id" bson:"game_id"`
	PlayerId  int    `json:"player_id" bson:"player_id"`
	PawnId    string `json:"pawn_id" bson:"pawn_id"`
	PawnIndex int    `json:"pawn_index" bson:"pawn_index"`
}
type PawnMovementResponse struct {
	GameId        string `json:"game_id" bson:"game_id"`
	PlayerId      int    `json:"player_id" bson:"player_id"`
	CurrentPlayer int    `json:"current_player" bson:"current_player"`
	PawnIndex     int    `json:"pawn_index" bson:"pawn_index"`
	PathIndex     int    `json:"new_path_index" bson:"new_path_index"`
	PawnId        string `json:"pawn_id" bson:"pawn_id"`
}

type DiceRollRequest struct {
	UserId   string `json:"user_id" bson:"user_id"`
	GameId   string `json:"game_id" bson:"game_id"`
	PlayerId int    `json:"player_id" bson:"player_id"`
}

type ChatMessage struct {
	UserId   string `json:"user_id" bson:"user_id"`
	Message  string `json:"message" bson:"message"`
	GameId   string `json:"game_id" bson:"game_id"`
	PlayerId int    `json:"player_id" bson:"player_id"`
}
