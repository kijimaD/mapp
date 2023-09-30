package world

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	Sprite            *ebiten.Image
	EnvironmentSprite *ebiten.Image
	HoverSprite       *ebiten.Image
	TileType          TileType
}

// levelは高さのことっぽい
type GameLevel struct {
	Tiles [][][]*Tile

	size int
}

func NewLevel(size int) *GameLevel {
	l := &GameLevel{
		size: size,
	}
	const startingLayers = 2
	for i := 0; i < startingLayers; i++ {
		l.AddLayer()
	}
	return l
}

// 平面にタイルを敷き詰めた配列を作る
func (l *GameLevel) AddLayer() {
	tileMap := make([][]*Tile, l.size)
	for x := 0; x < l.size; x++ {
		tileMap[x] = make([]*Tile, l.size)
		for y := 0; y < l.size; y++ {
			tileMap[x][y] = &Tile{}
		}
	}
	l.Tiles = append(l.Tiles, tileMap)
}

// タイル全体のHoverSpriteフラグをオフにする
func (l *GameLevel) ClearHoverSprites() {
	for i := range l.Tiles {
		for x := range l.Tiles[i] {
			for _, tile := range l.Tiles[i][x] {
				if tile == nil {
					continue
				}
				tile.HoverSprite = nil
			}
		}
	}
}
