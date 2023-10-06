package system

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/mapp/world"
	"github.com/sedyh/mizu/pkg/engine"
)

// デバッグテキストは右上に出る項目のこと

type renderDebugTextSystem struct{}

func NewRenderDebugTextSystem() *renderDebugTextSystem {
	return &renderDebugTextSystem{}
}

func (r *renderDebugTextSystem) Draw(w engine.World, screen *ebiten.Image) {
	if world.World.IsDebug == false {
		return
	}

	r.debugmenu(w).Draw(screen)
}

func (r *renderDebugTextSystem) debugmenu(w engine.World) *ebitenui.UI {
	face, _ := loadFont(20)
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.RowLayoutOpts.Spacing(4),
		)),
	)

	mouseX, mouseY := ebiten.CursorPosition()
	tileX, tileY := world.ScreenToCartesian(mouseX, mouseY)
	labeltext := fmt.Sprintf(`[DEBUG]%s
ENT %d
TPS %.0f
FPS %.0f
Mouse (%d,%d)
Tile (%.0f,%.0f)
HoverType %d
`,
		"               ",
		w.Entities(),
		ebiten.CurrentTPS(),
		ebiten.CurrentFPS(),
		mouseX,
		mouseY,
		tileX,
		tileY,
		world.World.HoverStructure,
	)
	label1 := widget.NewText(
		widget.TextOpts.Text(labeltext, face, color.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
			}),
		),
	)
	label1.GetWidget().Disabled = true
	rootContainer.AddChild(label1)

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	return &ui
}
