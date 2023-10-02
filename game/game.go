package game

import (
	"os"

	"github.com/kijimaD/mapp/system"
	"github.com/sedyh/mizu/pkg/engine"
)

type game struct{}

// NewGame returns a new isometric demo game.
func NewGame() *game {
	return &game{}
}

// Main scene, you can use that for settings or main menu
func (g *game) Setup(w engine.World) {
	w.AddComponents()
	w.AddEntities()
	w.AddSystems(
		&system.GeneralSystem{},
		&system.RenderDebugTextSystem{},
		&system.RenderHudSystem{},
		system.NewPlayerMoveSystem(),
		system.NewProfileSystem(),
		system.NewTickSystem(),
	)
}

func (g *game) Exit() {
	os.Exit(0)
}
