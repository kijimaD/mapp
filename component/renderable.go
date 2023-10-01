package component

import (
	"github.com/hajimehoshi/ebiten/v2"
)

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
