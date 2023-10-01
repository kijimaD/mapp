package entity

import (
	"code.rocketnine.space/tslocum/gohan"
	"github.com/kijimaD/mapp/component"
)

func NewPlayer() gohan.Entity {
	player := gohan.NewEntity()

	player.AddComponent(&component.Position{})

	return player
}
