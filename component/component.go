package component

import "github.com/hajimehoshi/ebiten/v2"

// 位置を持つ
type Position struct {
	X, Y float64
}

// 移動できる
type Moveable struct {
	X, Y float64
}

// 描画可能なもの
type Renderable struct {
	Image          *ebiten.Image // スプライト画像
	HorizontalFlip bool
	VerticalFlip   bool
	DiagonalFlip   bool // TODO unimplemented

	RenderOrder int // 描画の優先順位 下1 < 2 < 3上
	// Angle float64

	// Overlay            *ebiten.Image
	// OverlayX, OverlayY float64 // Overlay offset

	// Frame     int
	// Frames    []*ebiten.Image
	// FrameTime time.Duration
	// LastFrame time.Time
	// NumFrames int

	// OverrideColorScale bool
	// ColorScale         float64
}

// 説明を持つ
type Describable struct {
	Info string
}

// 名前を持つ
type Name struct {
	Name string
}
