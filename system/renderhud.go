package system

import (
	"image"
	"image/color"
	"strings"

	"github.com/ebitenui/ebitenui"
	uiimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kijimaD/mapp/component"
	"github.com/kijimaD/mapp/world"
)

const (
	helpW       = 480
	helpH       = 220
	columns     = 3
	buttonWidth = world.SidebarWidth / columns
)

type renderHudSystem struct {
	Position *component.Position

	op           *ebiten.DrawImageOptions
	hudImg       *ebiten.Image
	tmpImg       *ebiten.Image // ボタンを描画するベースになるのかな
	tmpImg2      *ebiten.Image
	helpImg      *ebiten.Image
	sidebarColor color.RGBA
	ui           *ebitenui.UI
}

func NewRenderHudSystem() *renderHudSystem {
	s := &renderHudSystem{
		op:      &ebiten.DrawImageOptions{},
		hudImg:  ebiten.NewImage(1, 1),
		tmpImg:  ebiten.NewImage(1, 1),
		tmpImg2: ebiten.NewImage(1, 1),
		helpImg: ebiten.NewImage(helpW, helpH),
		ui:      generateUI(),
	}

	sidebarShade := uint8(108)
	s.sidebarColor = color.RGBA{sidebarShade, sidebarShade, sidebarShade, 255}

	return s
}

func generateUI() *ebitenui.UI {
	buttonImage, _ := loadButtonImage()
	face, _ := loadFont(20)
	insets := widget.Insets{
		Top:    0,
		Left:   200,
		Right:  0,
		Bottom: 0,
	}
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(insets),
			widget.RowLayoutOpts.Spacing(4),
		)),
	)
	breakbtn := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				},
			),
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionStart,
					Stretch:  false,
				},
			),
			widget.WidgetOpts.MinSize(40, 40),
		),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("break", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   4,
			Right:  4,
			Top:    4,
			Bottom: 4,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button clicked")
		}),
	)
	breakbtn.ClickedEvent.AddHandler(func(args interface{}) {
		world.SetHoverStructure(world.StructureBulldozer)
	})

	roadbtn := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				},
			),
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionStart,
					Stretch:  false,
				},
			),
			widget.WidgetOpts.MinSize(40, 40),
		),

		widget.ButtonOpts.Image(buttonImage),

		widget.ButtonOpts.Text("road", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   4,
			Right:  4,
			Top:    4,
			Bottom: 4,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button clicked")
		}),
	)
	roadbtn.ClickedEvent.AddHandler(func(args interface{}) {
		world.SetHoverStructure(world.StructureRoad)
	})

	rootContainer.AddChild(breakbtn)
	rootContainer.AddChild(roadbtn)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	return &ui
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := uiimage.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})
	hover := uiimage.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})
	pressed := uiimage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func (r *renderHudSystem) Draw(w engine.World, screen *ebiten.Image) {
	if world.World.HUDUpdated {
		r.hudImg.Clear()
		r.drawSidebar()
		r.drawMessages()
		r.drawTooltip()
		r.drawHelp()
		world.World.HUDUpdated = false
	}
	// 1つの画像を作って描画するという感じか。それぞれ分けて作成した画像を重ね合わせるという感じではなく
	screen.DrawImage(r.hudImg, nil)

	r.ui.Draw(screen)
}

func (r *renderHudSystem) Update(w engine.World) {
	r.ui.Update()
}

func (s *renderHudSystem) drawSidebar() {
	bounds := s.hudImg.Bounds()
	if bounds.Dx() != world.World.ScreenW || bounds.Dy() != world.World.ScreenH {
		s.hudImg = ebiten.NewImage(world.World.ScreenW, world.World.ScreenH)
		s.tmpImg = ebiten.NewImage(world.World.ScreenW, world.World.ScreenH)
		s.tmpImg2 = ebiten.NewImage(world.SidebarWidth, world.World.ScreenH)
	} else {
		s.hudImg.Clear()
		s.tmpImg.Clear()
		s.tmpImg2.Clear()
	}

	// Fill background.
	s.hudImg.SubImage(image.Rect(0, 0, world.SidebarWidth, world.World.ScreenH)).(*ebiten.Image).Fill(s.sidebarColor)

	// Draw buttons.

	const paddingSize = 1
	const buttonHeight = buttonWidth
	world.World.HUDButtonRects = make([]image.Rectangle, len(world.HUDButtons))
	var lastButtonY int
	for i, button := range world.HUDButtons {
		row := i / columns
		x, y := (i%columns)*buttonWidth, row*buttonHeight
		r := image.Rect(x+paddingSize, y+paddingSize, x+buttonWidth-paddingSize, y+buttonHeight-paddingSize)

		if button != nil {
			selected := world.World.HoverStructure == button.StructureType
			if button.StructureType == world.StructureToggleHelp {
				selected = world.World.HelpPage != -1
			}

			// Draw background.
			s.drawButtonBackground(s.tmpImg, r, selected)

			// Draw sprite.
			colorScale := 1.0
			if selected {
				colorScale = 0.9
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x+paddingSize)+button.SpriteOffsetX, float64(y+paddingSize)+button.SpriteOffsetY)
			op.ColorM.Scale(colorScale, colorScale, colorScale, 1)
			s.tmpImg.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).DrawImage(button.Sprite, op)

			s.drawButtonBorder(s.tmpImg, r, selected)
		}

		world.World.HUDButtonRects[i] = r
		if button != nil {
			nonHUDButton := button.StructureType == world.StructureToggleHelp
			if !nonHUDButton {
				lastButtonY = y
			}
		}
	}

	dateY := lastButtonY + buttonHeight*2 - buttonHeight/2 - 16
	s.drawDate(dateY)
	s.drawFunds(dateY + 50)
	s.hudImg.DrawImage(s.tmpImg, nil)
	s.hudImg.SubImage(image.Rect(world.SidebarWidth-1, 0, world.SidebarWidth, world.World.ScreenH)).(*ebiten.Image).Fill(color.Black)
}

func (s *renderHudSystem) drawButtonBackground(img *ebiten.Image, r image.Rectangle, selected bool) {
	buttonShade := uint8(142)
	colorButton := color.RGBA{buttonShade, buttonShade, buttonShade, 255}

	bgColor := colorButton
	if selected {
		bgColor = s.sidebarColor
	}

	img.SubImage(r).(*ebiten.Image).Fill(bgColor)
}

func (s *renderHudSystem) drawButtonBorder(img *ebiten.Image, r image.Rectangle, selected bool) {
	borderSize := 2

	lightBorderShade := uint8(216)
	colorLightBorder := color.RGBA{lightBorderShade, lightBorderShade, lightBorderShade, 255}

	mediumBorderShade := uint8(56)
	colorMediumBorder := color.RGBA{mediumBorderShade, mediumBorderShade, mediumBorderShade, 255}

	darkBorderShade := uint8(42)
	colorDarkBorder := color.RGBA{darkBorderShade, darkBorderShade, darkBorderShade, 255}

	topLeftBorder := colorLightBorder
	bottomRightBorder := colorMediumBorder
	if selected {
		topLeftBorder = colorDarkBorder
		bottomRightBorder = colorLightBorder
	}

	// Draw top and left border.
	img.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+borderSize)).(*ebiten.Image).Fill(topLeftBorder)
	img.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Min.X+borderSize, r.Max.Y)).(*ebiten.Image).Fill(topLeftBorder)

	// Draw bottom and right border.
	img.SubImage(image.Rect(r.Min.X, r.Max.Y-borderSize, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(bottomRightBorder)
	img.SubImage(image.Rect(r.Max.X-borderSize, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(bottomRightBorder)
}

func (s *renderHudSystem) drawTooltip() {
	label := world.TooltipText()
	if label == "" {
		return
	}

	lines := 1 + strings.Count(label, "\n")
	max := maxLen(strings.Split(label, "\n"))

	scale := 3.0
	x, y := world.SidebarWidth, 0
	w, h := (max*6+10)*int(scale), 16*(int(scale))*lines+10
	r := image.Rect(x, y, x+w, y+h)
	s.hudImg.SubImage(r).(*ebiten.Image).Fill(color.RGBA{0, 0, 0, 120})

	s.tmpImg.Clear()
	ebitenutil.DebugPrint(s.tmpImg, label)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(world.SidebarWidth+(4*scale), 2)
	s.hudImg.DrawImage(s.tmpImg, op)
}

func maxLen(v []string) int {
	max := 0
	for _, line := range v {
		l := len(line)
		if l > max {
			max = l
		}
	}
	return max
}

func (s *renderHudSystem) drawMessages() {
	lines := len(world.World.Messages)
	if lines == 0 {
		return
	}
	/*var label string
	max := maxLen(world.World.Messages)
	for i := lines - 1; i >= 0; i-- {
		if i != lines-1 {
			label += "\n"
		}
		for j := max - len(world.World.Messages[i]); j > 0; j-- {
			label += " "
		}
		label += world.World.Messages[i]
	}*/

	label := world.World.Messages[len(world.World.Messages)-1]
	max := len(label)
	lines = 1

	const padding = 10

	scale := 2.0
	w, h := (max*6+10)*int(scale), 16*(int(scale))*lines+10
	x, y := world.World.ScreenW-w, 0
	r := image.Rect(x, y, x+w, y+h)
	s.hudImg.SubImage(r).(*ebiten.Image).Fill(color.RGBA{0, 0, 0, 120})

	s.tmpImg.Clear()
	ebitenutil.DebugPrint(s.tmpImg, label)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x)+padding, 3)
	s.hudImg.DrawImage(s.tmpImg, op)
}

func (s *renderHudSystem) drawDate(y int) {
	const datePadding = 10
	month, year := world.Date()
	label := month

	scale := 2.0
	x, y := datePadding, y

	s.tmpImg2.Clear()
	ebitenutil.DebugPrint(s.tmpImg2, label)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	s.hudImg.DrawImage(s.tmpImg2, op)

	label = year

	x = world.SidebarWidth - 1 - datePadding - (len(label) * 6 * int(scale))

	s.tmpImg2.Clear()
	ebitenutil.DebugPrint(s.tmpImg2, label)
	op.GeoM.Reset()
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	s.hudImg.DrawImage(s.tmpImg2, op)
}

func (s *renderHudSystem) drawFunds(y int) {
	label := world.World.Printer.Sprintf("$%d", world.World.Funds)

	scale := 2.0
	x, y := world.SidebarWidth/2-(len(label)*12)/2, y

	s.tmpImg2.Clear()
	ebitenutil.DebugPrint(s.tmpImg2, label)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	s.hudImg.DrawImage(s.tmpImg2, op)
}

func (s *renderHudSystem) drawHelp() {
	if world.World.HelpPage < 0 {
		return
	}

	if world.World.HelpUpdated {
		s.helpImg.Fill(s.sidebarColor)

		label := strings.TrimSpace(world.HelpText[world.World.HelpPage])

		s.tmpImg.Clear()
		ebitenutil.DebugPrint(s.tmpImg, label)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(5, 0)
		s.helpImg.DrawImage(s.tmpImg, op)

		s.helpImg.SubImage(image.Rect(0, 0, helpW, 1)).(*ebiten.Image).Fill(color.Black)
		s.helpImg.SubImage(image.Rect(0, 0, 1, helpH)).(*ebiten.Image).Fill(color.Black)

		// Draw prev/next buttons.
		buttonSize := 32
		buttonPadding := 4
		prevRect := image.Rect(buttonPadding+2, helpH-buttonSize-buttonPadding+1, buttonSize+buttonPadding+2, helpH-buttonPadding+1)
		closeRect := image.Rect(helpW/2-buttonSize/2, helpH-buttonSize-buttonPadding+1, helpW/2+buttonSize/2, helpH-buttonPadding+1)
		nextRect := image.Rect(helpW-buttonPadding, helpH-buttonSize-buttonPadding+1, helpW-buttonSize-buttonPadding, helpH-buttonPadding+1)

		drawButton := func(r image.Rectangle, l string) {
			s.drawButtonBackground(s.helpImg, r, false)
			ebitenutil.DebugPrintAt(s.helpImg, l, r.Min.X+buttonSize/2-4, r.Min.Y+buttonSize/2-10)
			s.drawButtonBorder(s.helpImg, r, false)
		}

		if world.World.HelpPage > 0 {
			drawButton(prevRect, "<")
		}
		drawButton(closeRect, "X")
		if world.World.HelpPage < len(world.HelpText)-1 {
			drawButton(nextRect, ">")
		}

		world.World.HelpButtonRects = []image.Rectangle{
			prevRect,
			closeRect,
			nextRect,
		}

		world.World.HelpUpdated = false
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(world.World.ScreenW)-helpW, float64(world.World.ScreenH)-helpH)
	s.hudImg.DrawImage(s.helpImg, op)
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
