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

	DisableEsc bool

	IsDebug bool
	NoClip  bool

	GameStarted      bool
	GameStartedTicks int
	GameOver         bool

	PlayerX, PlayerY float64

	CamX, CamY     float64
	CamScale       float64
	CamScaleTarget float64
	CamMoving      bool

	PlayerWidth  float64
	PlayerHeight float64

	HoverStructure         int
	HoverX, HoverY         int
	HoverLastX, HoverLastY int
	HoverValid             bool

	Map             *tiled.Map
	ObjectGroups    []*tiled.ObjectGroup
	HazardRects     []image.Rectangle
	CreepRects      []image.Rectangle
	CreepEntities   []gohan.Entity
	TriggerEntities []gohan.Entity
	TriggerRects    []image.Rectangle
	TriggerNames    []string

	NativeResolution bool

	BrokenPieceA, BrokenPieceB gohan.Entity

	TileImages         map[uint32]*ebiten.Image
	TileImagesFirstGID uint32

	ResetGame bool

	MuteMusic        bool
	MuteSoundEffects bool // TODO

	GotCursorPosition bool

	tilesets []*ebiten.Image

	EnvironmentSprites int

	SelectedStructure *Structure

	HUDUpdated     bool
	HUDButtonRects []image.Rectangle

	RCIButtonRect image.Rectangle
	RCIWindowRect image.Rectangle
	ShowRCIWindow bool

	HelpUpdated     bool
	HelpPage        int
	HelpButtonRects []image.Rectangle

	PowerPlants []*PowerPlant

	HavePowerOut bool
	PowerOuts    [][]bool

	Ticks int

	Paused bool

	Funds int

	Printer *message.Printer

	TransparentStructures bool

	Messages      []string // 右上に一時的に表示するメッセージ。MessagesとMessagesTicksのスライスの数は対応している
	MessagesTicks []int    // 右上に一時的に表示するメッセージの残り秒

	Power          PowerMap
	PowerUpdated   bool
	PowerAvailable int
	PowerNeeded    int

	BuildDragX int
	BuildDragY int

	LastBuildX int
	LastBuildY int

	TaxR float64
	TaxC float64
	TaxI float64

	resetTipShown bool
}

func (w *GameWorld) SetGameOver(vx, vy float64) {
	if w.GameOver {
		return
	}

	w.GameOver = true
}
