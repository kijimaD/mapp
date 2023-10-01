package game

import (
	"image/color"
	"os"
	"sync"

	"code.rocketnine.space/tslocum/gohan"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/mapp/asset"
	"github.com/kijimaD/mapp/entity"
	"github.com/kijimaD/mapp/system"
	"github.com/kijimaD/mapp/world"
)

const sampleRate = 44100

// game is an isometric demo game.
type game struct {
	w, h           int
	op             *ebiten.DrawImageOptions
	disableEsc     bool
	debugMode      bool
	cpuProfile     *os.File
	movementSystem *system.MovementSystem
	addedSystems   bool
	updateTicks    int
	sync.Mutex
}

// NewGame returns a new isometric demo game.
func NewGame() (*game, error) {
	g := &game{
		op:        &ebiten.DrawImageOptions{},
		debugMode: true,
	}

	err := g.loadAssets()
	if err != nil {
		panic(err)
	}

	const numEntities = 30000
	gohan.Preallocate(numEntities)

	return g, nil
}

// Layout is called when the game's layout changes.
func (g *game) Layout(w, h int) (int, int) {
	if w != g.w || h != g.h {
		world.World.ScreenW, world.World.ScreenH = w, h
		g.w, g.h = w, h
		world.World.HUDUpdated = true
	}
	return g.w, g.h
}

func (g *game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		g.Exit()
		return nil
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
			return err
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

		if world.World.Player == 0 {
			world.World.Player = entity.NewPlayer()
		}

		if !g.addedSystems {
			g.addSystems()
			g.addedSystems = true // TODO
		}

		world.World.ResetGame = false
		world.World.GameOver = false
	}

	err := gohan.Update()
	if err != nil {
		return err
	}
	return nil
}

// renderSprite renders a sprite on the screen.
func (g *game) renderSprite(
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

func (g *game) Draw(screen *ebiten.Image) {
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
				if tile.HoverSprite != nil {
					// プレビューは暗く描画する
					sprite = tile.HoverSprite
					colorScale = 0.6
					if !world.World.HoverValid {
						colorScale = 0.2
					}
				} else if tile.TileType == world.BusStopTile {
					img := world.GrassTile
					// TODO: インデックス直指定をやめる
					sprite = world.World.TileImages[img+world.World.TileImagesFirstGID+25]
				} else if tile.TileType == world.RoadTile {
					img := world.GrassTile
					sprite = world.World.TileImages[img+world.World.TileImagesFirstGID+24]
				} else if tile.TileType == world.PlainTile {
					img := world.GrassTile
					sprite = world.World.TileImages[img+world.World.TileImagesFirstGID]
				} else {
					continue
				}
				if i > 1 {
					alpha = 0.2
				}
				drawn += g.renderSprite(float64(x), float64(y), 0, float64(i*-heightFactor), 0, 1, colorScale, alpha, false, false, sprite, screen)
			}
		}
	}
	world.World.EnvironmentSprites = drawn

	err := gohan.Draw(screen)
	if err != nil {
		panic(err)
	}
}

func (g *game) addSystems() {
	// Simulation systems.
	gohan.AddSystem(system.NewTickSystem())

	// Input systems.
	g.movementSystem = system.NewMovementSystem()
	gohan.AddSystem(system.NewPlayerMoveSystem(world.World.Player, g.movementSystem))

	// Render systems.
	gohan.AddSystem(system.NewCameraSystem())
	gohan.AddSystem(system.NewRenderHudSystem())
	gohan.AddSystem(system.NewRenderDebugTextSystem(world.World.Player))
	gohan.AddSystem(system.NewProfileSystem(world.World.Player))
}

func (g *game) loadAssets() error {
	asset.ImgWhiteSquare.Fill(color.White)
	return nil
}

func (g *game) Exit() {
	os.Exit(0)
}
