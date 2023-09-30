package world

// 駅となるタイル

type StationMapTile struct {
	X        int
	Y        int
	Opened   bool
	Capacity int
	Name     string
}

func (t *StationMapTile) Up() *StationMapTile {
	tx, ty := t.X, t.Y-1
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Station[tx][ty]
	return n
}

func (t *StationMapTile) Down() *StationMapTile {
	tx, ty := t.X, t.Y+1
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Station[tx][ty]
	return n
}

func (t *StationMapTile) Left() *StationMapTile {
	tx, ty := t.X-1, t.Y
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Station[tx][ty]
	return n
}

func (t *StationMapTile) Right() *StationMapTile {
	tx, ty := t.X+1, t.Y
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Station[tx][ty]
	return n
}
