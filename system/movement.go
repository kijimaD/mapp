package system

import (
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/mapp/component"
)

type MovementSystem struct {
	Position *component.Position
}

func NewMovementSystem() *MovementSystem {
	s := &MovementSystem{}
	return s
}

func (s *MovementSystem) Update(e gohan.Entity) error {
	return nil
}

func (_ *MovementSystem) Draw(_ gohan.Entity, screen *ebiten.Image) error {
	return gohan.ErrUnregister
}
