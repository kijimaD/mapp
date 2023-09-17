package world

import "github.com/beefsack/go-astar"

// 電力がないのでこの箇所は必要ないが、後で参考になりそうなので残す
// 電力を運んでいるタイル

type PowerMapTile struct {
	X            int
	Y            int
	CarriesPower bool // Set to true for roads and all building tiles (even power plants)
}

// そのタイルの上下左右が伝送路であればそのタイルを返す
func (t *PowerMapTile) Up() *PowerMapTile {
	tx, ty := t.X, t.Y-1
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Power[tx][ty]
	if !n.CarriesPower {
		return nil
	}
	return n
}

func (t *PowerMapTile) Down() *PowerMapTile {
	tx, ty := t.X, t.Y+1
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Power[tx][ty]
	if !n.CarriesPower {
		return nil
	}
	return n
}

func (t *PowerMapTile) Left() *PowerMapTile {
	tx, ty := t.X-1, t.Y
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Power[tx][ty]
	if !n.CarriesPower {
		return nil
	}
	return n
}

func (t *PowerMapTile) Right() *PowerMapTile {
	tx, ty := t.X+1, t.Y
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Power[tx][ty]
	if !n.CarriesPower {
		return nil
	}
	return n
}

// aster.Patherのinterface構成関数のひとつ
// 近隣にあるタイルを返す
func (t *PowerMapTile) PathNeighbors() []astar.Pather {
	var neighbors []astar.Pather
	n := t.Up()
	if n != nil {
		neighbors = append(neighbors, n)
	}
	n = t.Down()
	if n != nil {
		neighbors = append(neighbors, n)
	}
	n = t.Left()
	if n != nil {
		neighbors = append(neighbors, n)
	}
	n = t.Right()
	if n != nil {
		neighbors = append(neighbors, n)
	}
	return neighbors
}

// aster.Patherのinterface構成関数のひとつ
// 経路地への1つの移動あたりのコスト計算。電気が運べるならコストは0
func (t *PowerMapTile) PathNeighborCost(to astar.Pather) float64 {
	toT := to.(*PowerMapTile)
	if !toT.CarriesPower {
		return 0
	}
	return 1
}

// aster.Patherのinterface構成関数のひとつ
// xとyの絶対値。移動距離がわかりそうなのだが、なぜこれを返すのかよくわからない
func (t *PowerMapTile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*PowerMapTile)
	absX := toT.X - t.X // 先 - 元
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Y - t.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}
