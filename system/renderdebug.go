package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kijimaD/mapp/world"
	"github.com/sedyh/mizu/pkg/engine"
)

// デバッグテキストは左上に出る項目のこと

type renderDebugTextSystem struct{}

func NewRenderDebugTextSystem() *renderDebugTextSystem {
	return &renderDebugTextSystem{}
}

func (r *renderDebugTextSystem) Draw(w engine.World, screen *ebiten.Image) {
	if world.World.IsDebug == false {
		return
	}
	r.debugPrint(screen)
}

func (r *renderDebugTextSystem) debugPrint(screen *ebiten.Image) {
	mouseX, mouseY := ebiten.CursorPosition()
	tileX, tileY := world.ScreenToCartesian(mouseX, mouseY)
	msg := fmt.Sprintf(`[DEBUG] - V%s
TPS %.0f
FPS %.0f
Mouse (%d,%d)
Tile (%.0f,%.0f)
HoverType %d
`,
		"               ",
		ebiten.CurrentTPS(),
		ebiten.CurrentFPS(),
		mouseX,
		mouseY,
		tileX,
		tileY,
		world.World.HoverStructure,
	)
	ebitenutil.DebugPrint(screen, msg)
}
