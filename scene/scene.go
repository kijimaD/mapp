package scene

import (
	"os"

	"github.com/kijimaD/mapp/system"
	"github.com/sedyh/mizu/pkg/engine"
)

type scene struct{}

// NewScene returns a new isometric demo scene.
func NewScene() *scene {
	return &scene{}
}

// Main scene, you can use that for settings or main menu
func (g *scene) Setup(w engine.World) {
	w.AddComponents()
	w.AddEntities()
	w.AddSystems(
		system.NewGeneralSystem(),
		system.NewRenderDebugTextSystem(),
		system.NewRenderHudSystem(),
		system.NewPlayerMoveSystem(),
		system.NewProfileSystem(),
		system.NewTickSystem(),
	)
}

func (g *scene) Exit() {
	os.Exit(0)
}
