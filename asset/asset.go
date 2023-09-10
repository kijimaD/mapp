package asset

import (
	"embed"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

const sampleRate = 44100

//go:embed image map
var FS embed.FS

var (
	ImgBlank       = ebiten.NewImage(1, 1)
	ImgWhiteSquare = ebiten.NewImage(64, 64)
	ImgBlackSquare = ebiten.NewImage(64, 64)
	ImgHelp        = LoadImage("image/help.png")
)

func init() {
	ImgWhiteSquare.Fill(color.White)
	ImgBlackSquare.Fill(color.Black)
}

func LoadImage(p string) *ebiten.Image {
	f, err := FS.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	baseImg, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(baseImg)
}
