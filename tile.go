package main



func convertCardTypeToTerrain(cardType int) string {
	switch cardType {
	case 0:
		return "l"
	case 1:
		return "b"
	case 2:
		return "w"
	case 3:
		return "g"
	case 4:
		return "o"
	case -1:
		return "?"
	case -2:
		return "d"
	case -3:
		return "s"
	case -4:
		return "-"
	}
	return ""
}
