package game

import (
	"os"

	"github.com/kijimaD/mapp/system"
	"github.com/sedyh/mizu/pkg/engine"
)

const sampleRate = 44100

// game is an isometric demo game.
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
		system.NewPlayerMoveSystem(),
		system.NewProfileSystem(),
	)
}

func (g *game) Exit() {
	os.Exit(0)
}
