// 世界
package world

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/mapp/asset"
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

var HUDButtons []*HUDButton
var CameraMinZoom = 0.1
var CameraMaxZoom = 10.0
var GrassTile = uint32(0)

var World = &GameWorld{
	CamScale:       startingZoom,
	CamScaleTarget: startingZoom,

	TileImages: make(map[uint32]*ebiten.Image),
	ResetGame:  true,
	Level:      NewLevel(256),

	BuildDragX: -1,
	BuildDragY: -1,
	LastBuildX: -1,
	LastBuildY: -1,

	Printer: message.NewPrinter(language.English),
	IsDebug: true,

	PreviewTileType: PlainTile,
}

type PowerPlant struct {
	Type int
	X, Y int
}

func Reset() {
	rand.Seed(time.Now().UnixNano())

	World.Funds = startingFunds

	World.CamX = float64((32 * TileSize) - rand.Intn(64*TileSize))
	World.CamY = float64((32 * TileSize) + rand.Intn(32*TileSize))
}

func loadMap(structureType int) (*tiled.Map, error) {
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

func DrawStructure(structureType int) *ebiten.Image {
	img := ebiten.NewImage(128, 128)

	m, err := loadMap(structureType)
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

				scale := 1.0 / float64(m.Width)
				paddingX := 64.0
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(
					xi+(paddingX*(float64(m.Width)-1)),
					(yi+float64(i*-40))+50)
				op.GeoM.Scale(scale, scale)
				img.DrawImage(tileImg, op)
			}
		}
	}

	return img
}

func LoadTileset() error {
	m, err := loadMap(StructureRoad)
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

func bulldozeArea(x int, y int, size int) {
	for dx := 0; dx < size; dx++ {
		for dy := 0; dy < size; dy++ {
			BuildStructure(StructureBulldozer, false, x-dx, y-dy, true)
		}
	}
}

func BuildStructure(structureType int, hover bool, placeX int, placeY int, internal bool) (*Structure, error) {
	m, err := loadMap(structureType)
	if err != nil {
		return nil, err
	}

	// TODO: w, hは1, 1とかが入る。何かわからない
	w := m.Width
	h := m.Height

	// 後の工程で(placeX-w, placeY-h), (placeX, placeY) を使う。これらが負の値になるとindexエラーになるのでチェックする
	if !ValidXY(placeX-w, placeY-h) || !ValidXY(placeX, placeY) {
		return nil, ErrInvalidBuildingNotFit
	}

	structure := &Structure{
		Type: structureType,
		X:    placeX,
		Y:    placeY,
	}

	// ブルドーザーを選択中に押すと削除する
	if structureType == StructureBulldozer && !hover {
		// TODO: 現在はタイル削除だけ。上にある建物削除をやる
		// TODO: タイルの階層は1層にする予定
		// 破壊する = その階層のタイルをnilに設定する
		bulldozed := false
		if World.Level.Tiles[0][placeX-w][placeY-w].TileType != PlainTile {
			World.Level.Tiles[0][placeX-w][placeY-w].TileType = PlainTile
			bulldozed = true
		}
		if !bulldozed {
			return nil, ErrNothingToBulldoze // FIXME: このエラーをメッセージに出したい
		}
		return structure, nil
	}

	valid := true
	// 道のタイルがすでにあるか判定
	// TODO: バス停の場合は道路上にのみ建設できる
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			tx, ty := (x+placeX)-w, (y+placeY)-h
			if !Buildable(structureType, tx, ty) && structureType != StructureBulldozer {
				valid = false
			}
		}
	}
	if hover {
		if structureType == StructureBulldozer {
			// ブルドーザーの場合は常に破壊できる
			World.HoverValid = true
		} else {
			World.HoverValid = valid
		}
	} else if !valid {
		return nil, ErrLocationOccupied
	}

	// ホバー時のタイルのプレビュー表示
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			tx, ty := (x+placeX)-w, (y+placeY)-h
			if hover {
				if Buildable(structureType, tx, ty) || structureType == StructureBulldozer {
					// タイルを平原にする
					// レベルは0なので、建設物の後ろのタイルを平原にセットする効果がある
					World.Level.Tiles[0][tx][ty].Hover = true
					World.PreviewTileType = PlainTile
				}
			}
		}
	}

	var t *tiled.LayerTile
	for i, layer := range m.Layers {
		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				t = layer.Tiles[y*m.Width+x]
				if t == nil {
					continue // No tile at this position.
				}

				tileImg := World.TileImages[t.Tileset.FirstGID+t.ID]
				if tileImg == nil {
					return nil, ErrTileImageNotFound
				}

				// 道路以外の建設物はベースタイル(階層0)の上に存在する
				// TODO: 道路破壊後に元のタイルを維持したいから、道路も自然タイルの上に置きたい
				layerNum := i // copy
				if structureType != StructureRoad {
					layerNum++
				}

				for layerNum > len(World.Level.Tiles)-1 {
					World.Level.AddLayer()
				}

				tx, ty := (x+placeX)-w, (y+placeY)-h
				if hover {
					if Buildable(structureType, tx, ty) || structureType == StructureBulldozer {
						// クリック中に出る建設プレビュー画像をセットする
						World.Level.Tiles[layerNum][tx][ty].Hover = true
						World.PreviewTileType = structureToTile(structureType)
					}
				} else {
					// クリックを離して建設する
					tiletype := structureToTile(structureType)
					World.Level.Tiles[layerNum][tx][ty].TileType = tiletype
				}
				// TODO handle flipping
			}
		}
	}

	return structure, nil
}

func structureToTile(structureType int) TileType {
	var tiletype TileType
	if structureType == StructureRoad {
		tiletype = RoadTile
	} else if structureType == StationBusStop {
		tiletype = BusStopTile
	} else if structureType == StructurePlain {
		tiletype = PlainTile
	}
	return tiletype
}

func StartGame() {
	if World.GameStarted {
		return
	}
	World.GameStarted = true

	// ヘルプページ非表示
	SetHelpPage(-1)
}
