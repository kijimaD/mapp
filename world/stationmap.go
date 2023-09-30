package world

// mapはx, yのtileで構成される
type StationMap [][]*StationMapTile

func newStationMap() StationMap {
	m := make(StationMap, 256)
	for x := 0; x < 256; x++ {
		m[x] = make([]*StationMapTile, 256)
		for y := 0; y < 256; y++ {
			m[x][y] = &StationMapTile{
				X:        x,
				Y:        y,
				Opened:   false,
				Capacity: 0,
				Name:     "No name",
			}
		}
	}
	return m
}

func (m StationMap) GetTile(x, y int) *StationMapTile {
	if !ValidXY(x, y) {
		return nil
	}
	return m[x][y]
}

func (m StationMap) SetTile(x, y int, opened bool) {
	t := m[x][y]
	if t.Opened == opened {
		return
	}
	t.Opened = opened

	// World.StationUpdated = true
}
