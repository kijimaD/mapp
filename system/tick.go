package system

import (
	"github.com/kijimaD/mapp/component"
	"github.com/kijimaD/mapp/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type TickSystem struct {
	Position *component.Position
	Velocity *component.Velocity
}

func NewTickSystem() *TickSystem {
	s := &TickSystem{}

	return s
}

func (s *TickSystem) Update(_ gohan.Entity) error {
	if world.World.Paused {
		return nil
	}

	// Update date display.
	if world.World.Ticks%world.MonthTicks == 0 {
		world.World.HUDUpdated = true
	}
	if world.World.Ticks%144 == 0 {
		world.TickMessages()
	}
	world.World.Ticks++
	return nil
}

func (s *TickSystem) Draw(e gohan.Entity, screen *ebiten.Image) error {
	return gohan.ErrUnregister
}
