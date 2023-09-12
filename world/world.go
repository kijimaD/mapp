// 世界
package world

import (
	"errors"
	"fmt"
	"image"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/mapp/asset"
	"github.com/kijimaD/mapp/component"
	"github.com/lafriks/go-tiled"
)

const startingYear = 1950
const maxPopulation = 100000
const (
	MonthTicks = 144 * 5
	YearTicks  = MonthTicks * 12
)
const TileSize = 64
const startingFunds = 100000
const startingZoom = 2.0
const SidebarWidth = 199
const startingTax = 0.12

type HUDButton struct {
	Sprite                       *ebiten.Image
	SpriteOffsetX, SpriteOffsetY float64
	Label                        string
	StructureType                int
}

var ErrNothingToBulldoze = errors.New("nothing to bulldoze")
var HUDButtons []*HUDButton
var CameraMinZoom = 0.1
var CameraMaxZoom = 10.0
var GrassTile = uint32(0)

var World = &GameWorld{
	CamScale:       startingZoom,
	CamScaleTarget: startingZoom,
	CamMoving:      true,

	PlayerWidth:  8,
	PlayerHeight: 32,

	TileImages: make(map[uint32]*ebiten.Image),
	ResetGame:  true,
	Level:      NewLevel(256),

	Power:     newPowerMap(),
	PowerOuts: newPowerOuts(),

	BuildDragX: -1,
	BuildDragY: -1,
	LastBuildX: -1,
	LastBuildY: -1,

	Printer: message.NewPrinter(language.English),
}

type PowerPlant struct {
	Type int
	X, Y int
}

func Reset() {
	for _, e := range gohan.AllEntities() {
		e.Remove()
	}
	World.Player = 0

	rand.Seed(time.Now().UnixNano())

	World.Funds = startingFunds

	World.ObjectGroups = nil
	World.HazardRects = nil
	World.CreepRects = nil
	World.CreepEntities = nil
	World.TriggerEntities = nil
	World.TriggerRects = nil
	World.TriggerNames = nil

	World.CamX = float64((32 * TileSize) - rand.Intn(64*TileSize))
	World.CamY = float64((32 * TileSize) + rand.Intn(32*TileSize))
}

func LoadMap(structureType int) (*tiled.Map, error) {
	filePath := StructureFilePaths[structureType]
	if filePath == "" {
		panic(fmt.Sprintf("unknown structure %d", structureType))
	}

	// Parse .tmx file.
	m, err := tiled.LoadFile(filePath, tiled.WithFileSystem(asset.FS))
	if err != nil {
		log.Fatalf("error parsing world: %+v", err)
	}

	return m, err
}

func DrawMap(structureType int) *ebiten.Image {
	img := ebiten.NewImage(128, 128)

	m, err := LoadMap(structureType)
	if err != nil {
		panic(err)
	}

	var t *tiled.LayerTile
	for i, layer := range m.Layers {
		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				t = layer.Tiles[y*m.Width+x]
				if t == nil || t.Nil {
					continue // No tile at this position.
				}

				tileImg := World.TileImages[t.Tileset.FirstGID+t.ID]
				if tileImg == nil {
					continue
				}

				xi, yi := CartesianToIso(float64(x), float64(y))

				scale := 0.9 / float64(m.Width)
				if m.Width < 2 {
					scale = 0.6
				}

				paddingX := 64.0
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(xi+(paddingX*(float64(m.Width)-1)), (yi+float64(i*-40))+92)
				op.GeoM.Scale(scale, scale)
				img.DrawImage(tileImg, op)
			}
		}
	}

	return img
}

func LoadTileset() error {
	m, err := LoadMap(StructureRoad)
	if err != nil {
		return err
	}

	// Load tileset.

	if len(World.tilesets) != 0 {
		return nil // Already loaded.
	}

	tileset := m.Tilesets[0]
	imgPath := filepath.Join("./image/tileset/", tileset.Image.Source)
	f, err := asset.FS.Open(filepath.ToSlash(imgPath))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	World.tilesets = append(World.tilesets, ebiten.NewImageFromImage(img))

	// Load tiles.

	for i := uint32(0); i < uint32(tileset.TileCount); i++ {
		rect := tileset.GetTileRect(i)
		World.TileImages[i+tileset.FirstGID] = World.tilesets[0].SubImage(rect).(*ebiten.Image)
	}

	World.TileImagesFirstGID = tileset.FirstGID
	return nil
}

func ShowBuildCost(structureType int, cost int) {
	if structureType == StructureBulldozer {
		ShowMessage(World.Printer.Sprintf("Bulldozed area (-$%d)", cost), 3)
	} else {
		ShowMessage(World.Printer.Sprintf("Built %s (-$%d)", strings.ToLower(StructureTooltips[World.HoverStructure]), cost), 3)
	}
}

func bulldozeArea(x int, y int, size int) {
	for dx := 0; dx < size; dx++ {
		for dy := 0; dy < size; dy++ {
			BuildStructure(StructureBulldozer, false, x-dx, y-dy, true)
		}
	}
}

func BuildStructure(structureType int, hover bool, placeX int, placeY int, internal bool) (*Structure, error) {
	m, err := LoadMap(structureType)
	if err != nil {
		return nil, err
	}

	if m.Width != 1 || m.Height != 1 {
		if placeX == 0 {
			placeX = 1
		}
		if placeY == 0 {
			placeY = 1
		}
	}

	w := m.Width - 1
	h := m.Height - 1

	if placeX-w < 0 || placeY-h < 0 || placeX >= 256 || placeY >= 256 {
		return nil, errors.New("invalid location: building does not fit")
	}

	structure := &Structure{
		Type: structureType,
		X:    placeX,
		Y:    placeY,
	}

	if structureType == StructureBulldozer && !hover {
		// TODO bulldoze entire structure, remove from zones
		var bulldozed bool
		for i := range World.Level.Tiles {
			if World.Level.Tiles[i][placeX][placeY].Sprite != nil {
				World.Level.Tiles[i][placeX][placeY].Sprite = nil
				bulldozed = true
			}

			var img *ebiten.Image
			if i == 0 {
				img = World.TileImages[GrassTile+World.TileImagesFirstGID]
			}
			if World.Level.Tiles[i][placeX][placeY].EnvironmentSprite != img {
				World.Level.Tiles[i][placeX][placeY].EnvironmentSprite = img
			}
		}
		if !bulldozed {
			return nil, ErrNothingToBulldoze
		}
		if !internal {
			checkSpaces := 2
			// PowerPlantはいらないので消していいが、参考になりそうなので残しておく
		REMOVEPOWER:
			for i, plant := range World.PowerPlants {
				for dx := 0; dx < checkSpaces; dx++ {
					for dy := 0; dy < checkSpaces; dy++ {
						if placeX == plant.X-dx && placeY == plant.Y-dy {
							World.PowerPlants = append(World.PowerPlants[:i], World.PowerPlants[i+1:]...)
							bulldozeArea(plant.X, plant.Y, 5)
							World.PowerUpdated = true
							break REMOVEPOWER
						}
					}
				}
			}
		}
		World.Power.SetTile(placeX, placeY, false)
		return structure, nil
	}

	createTileEntity := func(t *tiled.LayerTile, x float64, y float64) gohan.Entity {
		mapTile := gohan.NewEntity()
		mapTile.AddComponent(&component.Position{
			X: x,
			Y: y,
		})

		mapTile.AddComponent(&component.Sprite{
			Image:          World.TileImages[t.Tileset.FirstGID+t.ID],
			HorizontalFlip: t.HorizontalFlip,
			VerticalFlip:   t.VerticalFlip,
			DiagonalFlip:   t.DiagonalFlip,
		})

		return mapTile
	}
	_ = createTileEntity

	// TODO Add entity

	tileOccupied := func(tx int, ty int) bool {
		return World.Level.Tiles[1][tx][ty].Sprite != nil ||
			(World.Level.Tiles[0][tx][ty].Sprite != nil &&
				(structureType != StructureRoad ||
					World.Level.Tiles[0][tx][ty].Sprite != World.TileImages[World.TileImagesFirstGID]))
	}

	valid := true
	var existingRoadTiles int
VALIDBUILD:
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			tx, ty := (x+placeX)-w, (y+placeY)-h
			if structureType == StructureRoad && World.Level.Tiles[0][tx][ty].Sprite == World.TileImages[World.TileImagesFirstGID] {
				existingRoadTiles++
			}
			if tileOccupied(tx, ty) && structureType != StructureBulldozer {
				valid = false
				break VALIDBUILD
			}
		}
	}
	if structureType == StructureRoad && existingRoadTiles == 4 {
		valid = false
	}
	if hover {
		if structureType == StructureBulldozer {
			World.HoverValid = true
		} else {
			World.HoverValid = valid
		}
	} else if !valid {
		return nil, errors.New("invalid location: space already occupied")
	}

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			tx, ty := (x+placeX)-w, (y+placeY)-h
			if hover {
				if !tileOccupied(tx, ty) || structureType == StructureBulldozer {
					if structureType != StructureBulldozer {
						World.Level.Tiles[0][tx][ty].HoverSprite = World.TileImages[World.TileImagesFirstGID]
					}
					// Hide environment sprites temporarily.
					for i := 1; i < len(World.Level.Tiles); i++ {
						World.Level.Tiles[i][tx][ty].HoverSprite = asset.ImgBlank
					}
				}
			} else {
				World.Level.Tiles[0][tx][ty].Sprite = World.TileImages[World.TileImagesFirstGID]
				World.Level.Tiles[0][tx][ty].EnvironmentSprite = nil
				World.Level.Tiles[1][tx][ty].EnvironmentSprite = nil
			}
		}
	}

	var t *tiled.LayerTile
	for i, layer := range m.Layers {
		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				t = layer.Tiles[y*m.Width+x]
				if t == nil || t.Nil {
					continue // No tile at this position.
				}

				tileImg := World.TileImages[t.Tileset.FirstGID+t.ID]
				if tileImg == nil {
					continue
				}

				layerNum := i
				if structureType != StructureRoad {
					layerNum++
				}

				for layerNum > len(World.Level.Tiles)-1 {
					World.Level.AddLayer()
				}

				tx, ty := (x+placeX)-w, (y+placeY)-h
				if hover {
					if !tileOccupied(tx, ty) || structureType == StructureBulldozer {
						World.Level.Tiles[layerNum][tx][ty].HoverSprite = World.TileImages[t.Tileset.FirstGID+t.ID]
					}
				} else {
					World.Level.Tiles[layerNum][tx][ty].Sprite = World.TileImages[t.Tileset.FirstGID+t.ID]

					if structureType == StructureRoad {
						World.Power.SetTile(tx, ty, true)
					}
				}

				// TODO handle flipping
			}
		}
	}

	return structure, nil
}

func StartGame() {
	if World.GameStarted {
		return
	}
	World.GameStarted = true

	// Show initial help page.
	SetHelpPage(0)
}

func HUDButtonAt(x, y int) *HUDButton {
	point := image.Point{x, y}
	for i, rect := range World.HUDButtonRects {
		if point.In(rect) {
			return HUDButtons[i]
		}
	}
	return nil
}

func AltButtonAt(x, y int) int {
	point := image.Point{x, y}
	if point.In(World.RCIButtonRect) {
		return 0
	}
	return -1
}

func SetHoverStructure(structureType int) {
	World.HoverStructure = structureType
	World.HUDUpdated = true
}

func Tooltip() string {
	tooltipText := StructureTooltips[World.HoverStructure]
	cost := StructureCosts[World.HoverStructure]
	if cost > 0 {
		tooltipText += World.Printer.Sprintf("\n$%d", cost)
	}
	return tooltipText
}
