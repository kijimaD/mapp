package system

import (
	"fmt"
	"image/color"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kijimaD/mapp/world"
	"github.com/sedyh/mizu/pkg/engine"
)

// デバッグテキストは左下に出る項目のこと

type RenderDebugTextSystem struct {
	op       *ebiten.DrawImageOptions
	debugImg *ebiten.Image
}

func (s *RenderDebugTextSystem) Draw(w engine.World, screen *ebiten.Image) {
	if s.op == nil {
		s.op = &ebiten.DrawImageOptions{}
		s.debugImg = ebiten.NewImage(100, 200)
	}

	if world.World.IsDebug == false {
		return
	}
	s.debugImg.Fill(color.RGBA{0, 0, 0, 80})
	mouseX, mouseY := ebiten.CursorPosition()
	tileX, tileY := world.ScreenToCartesian(mouseX, mouseY)
	ebitenutil.DebugPrintAt(
		s.debugImg,
		fmt.Sprintf(`[DEBUG]
ENV %d
ENT %d
UPD %d
DRA %d
TPS %0.0f
FPS %0.0f
Mouse (%d,%d)
Tile (%.0f,%.0f)
Hover %d
`,
			world.World.EnvironmentSprites,
			gohan.CurrentEntities(),
			gohan.CurrentUpdates(),
			gohan.CurrentDraws(),
			ebiten.CurrentTPS(),
			ebiten.CurrentFPS(),
			mouseX,
			mouseY,
			tileX,
			tileY,
			world.World.HoverStructure,
		),
		0,
		0,
	)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(0, 800)
	screen.DrawImage(s.debugImg, op)
	return
}
