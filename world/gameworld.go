package world

import (
	"image"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"golang.org/x/text/message"
)

type GameWorld struct {
	Level *GameLevel

	Player gohan.Entity

	ScreenW, ScreenH int

	IsDebug bool

	GameStarted      bool
	GameStartedTicks int
	GameOver         bool

	PlayerX, PlayerY float64

	CamX, CamY     float64
	CamScale       float64
	CamScaleTarget float64

	HoverStructure         int // 選択中の建設物
	HoverX, HoverY         int
	HoverLastX, HoverLastY int
	HoverValid             bool

	Map *tiled.Map

	TileImages         map[uint32]*ebiten.Image
	TileImagesFirstGID uint32

	ResetGame bool

	GotCursorPosition bool

	tilesets []*ebiten.Image

	EnvironmentSprites int

	HUDUpdated     bool
	HUDButtonRects []image.Rectangle

	HelpUpdated     bool
	HelpPage        int
	HelpButtonRects []image.Rectangle

	Ticks int
	Funds int

	Printer *message.Printer

	Messages      []string // 右上に一時的に表示するメッセージ。MessagesとMessagesTicksのスライスの数は対応している
	MessagesTicks []int    // 右上に一時的に表示するメッセージの残り秒

	BuildDragX int
	BuildDragY int

	LastBuildX int
	LastBuildY int
}

func (w *GameWorld) SetGameOver(vx, vy float64) {
	if w.GameOver {
		return
	}

	w.GameOver = true
}
