package system

import (
	"github.com/kijimaD/mapp/world"
	"github.com/sedyh/mizu/pkg/engine"
)

type TickSystem struct{}

func NewTickSystem() *TickSystem {
	s := &TickSystem{}

	return s
}

func (s *TickSystem) Update(w engine.World) {
	// Update date display.
	if world.World.Ticks%world.MonthTicks == 0 {
		world.World.HUDUpdated = true
	}
	if world.World.Ticks%144 == 0 {
		// 右上に出るメッセージ
		world.TickMessages()
	}
	world.World.Ticks++
}
