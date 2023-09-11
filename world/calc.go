package world

func ValidXY(x, y int) bool {
	return x >= 0 && y >= 0 && x < 256 && y < 256
}
