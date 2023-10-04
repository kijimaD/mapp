package system

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/mapp/asset"
	"github.com/kijimaD/mapp/world"
	"github.com/sedyh/mizu/pkg/engine"
)

const sampleRate = 44100

type generalSystem struct {
	w, h        int
	op          *ebiten.DrawImageOptions
	debugMode   bool
	updateTicks int
}

func NewGeneralSystem() *generalSystem {
	return &generalSystem{
		op:          &ebiten.DrawImageOptions{},
		updateTicks: 0,
	}
}

func (g *generalSystem) Update(w engine.World) {
	if ebiten.IsWindowBeingClosed() {
		g.Exit()
		return
	}

	const updateSidebarDelay = 144 * 3
	g.updateTicks++
	if g.updateTicks == updateSidebarDelay {
		world.World.HUDUpdated = true
		//g.updateTicks = 0
		// TODO
	}

	if world.World.ResetGame {
		world.Reset()

		err := world.LoadTileset()
		if err != nil {
			return
		}

		// 平原レイヤーで埋める
		for x := range world.World.Level.Tiles[0] {
			for y := range world.World.Level.Tiles[0][x] {
				world.World.Level.Tiles[0][x][y].TileType = world.PlainTile
			}
		}

		// Load HUD sprites.
		op := &ebiten.DrawImageOptions{}
		op.ColorM.Scale(1, 1, 1, 0.4)
		world.HUDButtons = []*world.HUDButton{
			{
				StructureType: world.StructureBulldozer,
				Sprite:        world.DrawStructure(world.StructureBulldozer),
				SpriteOffsetX: 0,
				SpriteOffsetY: -60,
			},
			nil,
			{
				StructureType: world.StructureRoad,
				Sprite:        world.DrawStructure(world.StructureRoad),
				SpriteOffsetX: 0,
				SpriteOffsetY: -60,
			},
			{
				StructureType: world.StationBusStop,
				Sprite:        world.DrawStructure(world.StationBusStop),
				SpriteOffsetX: 0,
				SpriteOffsetY: -60,
			},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			{
				StructureType: world.StructureToggleHelp,
				Sprite:        asset.ImgHelp,
				SpriteOffsetX: 0,
				SpriteOffsetY: -1,
			},
		}

		world.World.ResetGame = false
		world.World.GameOver = false
	}

	return
}

func (g *generalSystem) renderSprite(
	x float64,
	y float64,
	offsetx float64,
	offsety float64,
	angle float64,
	geoScale float64,
	colorScale float64,
	alpha float64,
	hFlip bool,
	vFlip bool,
	sprite *ebiten.Image,
	target *ebiten.Image,
) int {
	if alpha < .01 {
		return 0
	}

	xi, yi := world.CartesianToIso(float64(x), float64(y))

	padding := float64(world.TileSize) * world.World.CamScale
	cx, cy := float64(world.World.ScreenW/2), float64(world.World.ScreenH/2)

	// Skip drawing tiles that are out of the screen.
	drawX, drawY := world.IsoToScreen(xi, yi)
	if drawX+padding < 0 || drawY+padding < 0 || drawX-padding > float64(world.World.ScreenW) || drawY-padding > float64(world.World.ScreenH) {
		return 0
	}

	g.op.GeoM.Reset()

	if hFlip {
		g.op.GeoM.Scale(-1, 1)
		g.op.GeoM.Translate(world.TileSize, 0)
	}
	if vFlip {
		g.op.GeoM.Scale(1, -1)
		g.op.GeoM.Translate(0, world.TileSize)
	}

	// Move to current isometric position.
	g.op.GeoM.Translate(xi, yi+offsety)
	// Translate camera position.
	g.op.GeoM.Translate(-world.World.CamX, -world.World.CamY)
	// Zoom.
	g.op.GeoM.Scale(world.World.CamScale, world.World.CamScale)
	// Center.
	g.op.GeoM.Translate(cx, cy)

	g.op.ColorM.Reset()
	g.op.ColorM.Scale(colorScale, colorScale, colorScale, alpha)

	target.DrawImage(sprite, g.op)
	g.op.ColorM.Reset()

	return 1
}

func (g *generalSystem) Draw(w engine.World, screen *ebiten.Image) {
	const heightFactor = 10 // 1つ階層が上がるとどれだけ上方向にずらして表示するか
	// タイル描画。タイルはEntityになっていない
	// エンティティ個別に描くのではなく、タイルそれぞれについてイテレートして描画する
	// 空間と物があるとして、空間を先にループさせる、ような感じ。そっちのほうが効率的になりそうなのはわかる
	// タイルは静的であまり変わることがない。といいつつ地形改変もある
	var drawn int
	for i := range world.World.Level.Tiles {
		for x := range world.World.Level.Tiles[i] {
			for y, tile := range world.World.Level.Tiles[i][x] {
				if tile == nil {
					continue
				}
				var sprite *ebiten.Image
				colorScale := 1.0
				alpha := 1.0

				sprite, err := g.tileToImage(tile.TileType)
				if err != nil {
					continue
				}
				drawn += g.renderSprite(float64(x), float64(y), 0, float64(i*-heightFactor), 0, 1, colorScale, alpha, false, false, sprite, screen)

				// プレビュー表示
				if tile.Hover {
					colorScale = 0.6
					alpha := 0.8
					previewSprite, err := g.tileToImage(world.World.PreviewTileType)
					if err == nil {
						drawn += g.renderSprite(float64(x), float64(y), 0, float64(i*-heightFactor), 0, 1, colorScale, alpha, false, false, previewSprite, screen)
					}
				}
			}
		}
	}
	world.World.EnvironmentSprites = drawn
}

// TODO: インデックス直指定をやめる
func (g *generalSystem) tileToImage(tileType world.TileType) (*ebiten.Image, error) {
	var sprite *ebiten.Image
	img := world.GrassTile
	if tileType == world.BusStopTile {
		sprite = world.World.TileImages[img+world.World.TileImagesFirstGID+25]
	} else if tileType == world.RoadTile {
		sprite = world.World.TileImages[img+world.World.TileImagesFirstGID+24]
	} else if tileType == world.PlainTile {
		sprite = world.World.TileImages[img+world.World.TileImagesFirstGID]
	} else {
		return nil, fmt.Errorf("not found tile")
	}
	return sprite, nil
}

func (g *generalSystem) loadAssets() error {
	asset.ImgWhiteSquare.Fill(color.White)
	return nil
}

func (g *generalSystem) Exit() {
	os.Exit(0)
}
