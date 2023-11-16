package system

import (
	"image/color"

	uiimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kijimaD/mapp/world"
)

func breakbtn() *widget.Button {
	buttonImage, _ := loadButtonImage()
	face, _ := loadFont(20)
	btn := widget.NewButton(
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
			world.SetHoverStructure(world.StructureBulldozer)
		}),
	)

	return btn
}

func roadbtn() *widget.Button {
	buttonImage, _ := loadButtonImage()
	face, _ := loadFont(20)
	btn := widget.NewButton(
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
			world.SetHoverStructure(world.StructureRoad)
		}),
	)

	return btn
}

func busstopbtn() *widget.Button {
	buttonImage, _ := loadButtonImage()
	face, _ := loadFont(20)
	btn := widget.NewButton(
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

		widget.ButtonOpts.Text("bus stop", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   4,
			Right:  4,
			Top:    4,
			Bottom: 4,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			world.SetHoverStructure(world.StationBusStop)
		}),
	)

	return btn
}

func helpbtn() *widget.Button {
	buttonImage, _ := loadButtonImage()
	face, _ := loadFont(20)
	btn := widget.NewButton(
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

		widget.ButtonOpts.Text("help", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   4,
			Right:  4,
			Top:    4,
			Bottom: 4,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if world.World.HelpPage != -1 {
				world.SetHelpPage(-1) // 閉じる
			} else {
				world.SetHelpPage(0) // 開く
			}
		}),
	)

	return btn
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
