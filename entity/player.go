package entity

import (
	"github.com/kijimaD/mapp/component"
	"code.rocketnine.space/tslocum/gohan"
)

func NewPlayer() gohan.Entity {
	player := gohan.NewEntity()

	player.AddComponent(&component.Position{})
	player.AddComponent(&component.Velocity{})

	return player
}
