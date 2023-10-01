package system

import (
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/mapp/component"
	"github.com/kijimaD/mapp/world"
)

type TickSystem struct {
	Position *component.Position
}

func NewTickSystem() *TickSystem {
	s := &TickSystem{}

	return s
}

func (s *TickSystem) Update(_ gohan.Entity) error {
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
