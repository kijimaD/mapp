package system

import (
	"fmt"
	"image/color"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kijimaD/mapp/world"
)

// デバッグテキストは左下に出る項目のこと

type RenderDebugTextSystem struct {
	player   gohan.Entity
	op       *ebiten.DrawImageOptions
	debugImg *ebiten.Image
}

func NewRenderDebugTextSystem(player gohan.Entity) *RenderDebugTextSystem {
	s := &RenderDebugTextSystem{
		player:   player,
		op:       &ebiten.DrawImageOptions{},
		debugImg: ebiten.NewImage(100, 200),
	}
	return s
}

func (s *RenderDebugTextSystem) Update(_ gohan.Entity) error {
	return gohan.ErrUnregister
}

func (s *RenderDebugTextSystem) Draw(e gohan.Entity, screen *ebiten.Image) error {
	// TODO: カメラの位置、タイルの位置を表示したい

	if world.World.Debug <= 0 {
		return nil
	}
	s.debugImg.Fill(color.RGBA{0, 0, 0, 80})
	ebitenutil.DebugPrintAt(
		s.debugImg,
		fmt.Sprintf("[DEBUG]\nENV %d\nENT %d\nUPD %d\nDRA %d\nTPS %0.0f\nFPS %0.0f\n",
			world.World.EnvironmentSprites,
			gohan.CurrentEntities(),
			gohan.CurrentUpdates(),
			gohan.CurrentDraws(),
			ebiten.CurrentTPS(),
			ebiten.CurrentFPS(),
		),
		0,
		0,
	)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(0, 800)
	screen.DrawImage(s.debugImg, op)
	return nil
}
