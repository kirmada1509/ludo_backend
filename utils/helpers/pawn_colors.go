package helpers


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
