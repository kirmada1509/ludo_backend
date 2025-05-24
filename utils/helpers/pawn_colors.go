package helpers

import "math/rand"


func GetColor(playerId int) string {
	switch playerId {
	case 0:
		return "green"
	case 1:
		return "yellow"
	case 2:
		return "blue"
	case 3:
		return "red"
	default:
		return ""
	}
}

func GetHomePosition(playerId int) int {
	switch playerId {
	case 0:
		return 1
	case 1:
		return 14
	case 2:
		return 27
	case 3:
		return 40
	default:
		return -1
	}
}

func RollDice() int {
	return 1 + rand.Intn(6)
}