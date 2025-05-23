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