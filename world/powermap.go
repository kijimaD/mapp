package world

// mapはx, yのtileで構成される
type PowerMap [][]*PowerMapTile

func newPowerMap() PowerMap {
	m := make(PowerMap, 256)
	for x := 0; x < 256; x++ {
		m[x] = make([]*PowerMapTile, 256)
		for y := 0; y < 256; y++ {
			m[x][y] = &PowerMapTile{
				X: x,
				Y: y,
			}
		}
	}
	return m
}

func newPowerOuts() [][]bool {
	m := make([][]bool, 256)
	for x := 0; x < 256; x++ {
		m[x] = make([]bool, 256)
	}
	return m
}

func ResetPowerOuts() {
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			World.PowerOuts[x][y] = false
		}
	}
	World.HavePowerOut = false
}

func (m PowerMap) GetTile(x, y int) *PowerMapTile {
	if !ValidXY(x, y) {
		return nil
	}
	return m[x][y]
}

func (m PowerMap) SetTile(x, y int, carriesPower bool) {
	t := m[x][y]
	if t.CarriesPower == carriesPower {
		return
	}
	t.CarriesPower = carriesPower

	World.PowerUpdated = true
}
