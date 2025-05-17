package helpers

import "math/rand"


func GetColor(playerId int) string {
	switch playerId {
	case 0:
		return "red"
	case 1:
		return "green"
	case 2:
		return "blue"
	case 3:
		return "yellow"
	default:
		return ""
	}
}

func RollDice() int {
	return 1 + rand.Intn(6)
}