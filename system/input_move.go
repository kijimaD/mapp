package system

import (
	"errors"
	"os"
	"strings"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/mapp/component"
	"github.com/kijimaD/mapp/world"
)

type playerMoveSystem struct {
	Position *component.Position

	player   gohan.Entity
	movement *MovementSystem

	scrollDragX, scrollDragY         int
	scrollCamStartX, scrollCamStartY float64
}

func NewPlayerMoveSystem(player gohan.Entity, m *MovementSystem) *playerMoveSystem {
	return &playerMoveSystem{
		player:      player,
		movement:    m,
		scrollDragX: -1,
		scrollDragY: -1,
	}
}
func (s *playerMoveSystem) buildStructure(structureType int, tileX int, tileY int, playSound bool) (*world.Structure, error) {
	cost := world.StructureCosts[structureType]
	if world.World.Funds < cost {
		world.ShowMessage("Insufficient funds", 3)
		return nil, errors.New("insufficient funds")
	}

	structure, err := world.BuildStructure(world.World.HoverStructure, false, tileX, tileY, false)
	if err == nil || world.World.HoverStructure == world.StructureBulldozer {
		world.World.LastBuildX, world.World.LastBuildY = tileX, tileY

		if err == nil {
			world.World.Funds -= cost
		}

		world.World.HUDUpdated = true
	} else {
		dX := tileX - world.World.LastBuildX
		if dX < 0 {
			dX *= -1
		}
		dY := tileY - world.World.LastBuildY
		if dY < 0 {
			dY *= -1
		}
		if (dX > 1 || dY > 1) && err != world.ErrNothingToBulldoze {
			errMessage := err.Error()
			if len(errMessage) > 0 {
				errMessage = strings.ToUpper(errMessage[0:1]) + errMessage[1:]
			}
			world.ShowMessage(errMessage, 3)
		}
	}
	return structure, err
}

func (s *playerMoveSystem) Update(e gohan.Entity) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
		return nil
	}

	// デバッグモード
	if ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyV) {
		world.World.IsDebug = !world.World.IsDebug
		return nil
	}

	if world.World.GameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			world.World.ResetGame = true
		}
		return nil
	}

	// Update target zoom level.
	var scrollY float64
	if ebiten.IsKeyPressed(ebiten.KeyC) || ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		scrollY = -0.25
	} else if ebiten.IsKeyPressed(ebiten.KeyE) || ebiten.IsKeyPressed(ebiten.KeyPageUp) {
		scrollY = .25
	} else {
		_, scrollY = ebiten.Wheel()
		if scrollY < -1 {
			scrollY = -1
		} else if scrollY > 1 {
			scrollY = 1
		}
	}
	world.World.CamScaleTarget += scrollY * (world.World.CamScaleTarget / 7)
	if world.World.CamScaleTarget < world.CameraMinZoom {
		world.World.CamScaleTarget = world.CameraMinZoom
	} else if world.World.CamScaleTarget > world.CameraMaxZoom {
		world.World.CamScaleTarget = world.CameraMaxZoom
	}

	// Smooth zoom transition.
	div := 10.0
	if world.World.CamScaleTarget > world.World.CamScale {
		world.World.CamScale += (world.World.CamScaleTarget - world.World.CamScale) / div
	} else if world.World.CamScaleTarget < world.World.CamScale {
		world.World.CamScale -= (world.World.CamScale - world.World.CamScaleTarget) / div
	}

	pressLeft := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	pressRight := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	pressUp := ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	pressDown := ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS)

	const camSpeed = 10
	if (pressLeft && !pressRight) ||
		(pressRight && !pressLeft) {
		if pressLeft {
			world.World.CamX -= camSpeed
		} else {
			world.World.CamX += camSpeed
		}
	}

	if (pressUp && !pressDown) ||
		(pressDown && !pressUp) {
		if pressUp {
			world.World.CamY -= camSpeed
		} else {
			world.World.CamY += camSpeed
		}
	}

	x, y := ebiten.CursorPosition()
	if !world.World.GotCursorPosition {
		if x != 0 || y != 0 {
			world.World.GotCursorPosition = true
		} else {
			return nil
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		if s.scrollDragX == -1 && s.scrollDragY == -1 {
			// TODO Disabled due to possible ebiten bug.
			//ebiten.SetCursorMode(ebiten.CursorModeCaptured)

			s.scrollDragX, s.scrollDragY = x, y
			s.scrollCamStartX, s.scrollCamStartY = world.World.CamX, world.World.CamY
		} else {
			dx, dy := float64(x-s.scrollDragX)/world.World.CamScale, float64(y-s.scrollDragY)/world.World.CamScale
			world.World.CamX, world.World.CamY = s.scrollCamStartX-dx, s.scrollCamStartY-dy
		}
	} else {
		if s.scrollDragX != -1 && s.scrollDragY != -1 {
			s.scrollDragX, s.scrollDragY = -1, -1
			//ebiten.SetCursorMode(ebiten.CursorModeVisible)
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			vX, vY := world.World.ScreenW/2-x, world.World.ScreenH/2-y
			dx, dy := float64(vX)/world.World.CamScale, float64(vY)/world.World.CamScale
			world.World.CamX, world.World.CamY = world.World.CamX-dx, world.World.CamY-dy
		}
	}
	// Clamp viewport.
	minCam := -256.0 * world.TileSize / 2
	maxCam := 256.0 * world.TileSize / 2
	if world.World.CamX < minCam {
		world.World.CamX = minCam
	} else if world.World.CamX > maxCam {
		world.World.CamX = maxCam
	}
	if world.World.CamY < 0 {
		world.World.CamY = 0
	} else if world.World.CamY > maxCam {
		world.World.CamY = maxCam
	}

	// サイドバー
	if x < world.SidebarWidth {
		world.World.Level.ClearHoverSprites()
		world.World.HoverX, world.World.HoverY = 0, 0

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			button := world.HUDButtonAt(x, y)
			if button != nil {
				if button.StructureType == world.StructureToggleHelp {
					if world.World.HelpPage != -1 {
						world.SetHelpPage(-1) // 閉じる
					} else {
						world.SetHelpPage(0) // 開く
					}
				} else {
					if world.World.HoverStructure == button.StructureType {
						world.SetHoverStructure(0) // Deselect.
					} else {
						world.SetHoverStructure(button.StructureType)
					}
				}
			}
		}
		return nil
	}

	// ヘルプページ
	if x >= world.World.ScreenW-helpW && y >= world.World.ScreenH-helpH {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			const (
				helpPrev = iota
				helpClose
				helpNext
			)

			helpButton := world.HelpButtonAt(x-(world.World.ScreenW-helpW), y-(world.World.ScreenH-helpH))
			var updated bool
			switch helpButton {
			case helpPrev:
				if world.World.HelpPage > 0 {
					world.World.HelpPage--
					updated = true
				}
			case helpClose:
				world.World.HelpPage = -1
				updated = true
			case helpNext:
				if world.World.HelpPage < len(world.HelpText)-1 {
					world.World.HelpPage++
					updated = true
				}
			}
			if updated {
				world.World.HelpUpdated = true
				world.World.HUDUpdated = true
			}
		}
		return nil
	}

	if world.World.HoverStructure != 0 {
		roadTiles := func(fromX, fromY, toX, toY int) [][2]int {
			var tiles [][2]int
			fx, fy := float64(fromX), float64(fromY)
			tx, ty := float64(toX), float64(toY)
			dx, dy := tx-fx, ty-fy
			for dx < -1 || dx > 1 || dy < -1 || dy > 1 {
				dx /= 2
				dy /= 2
			}
			tiles = append(tiles, [2]int{fromX, fromY})
			for fx != tx || fy != ty {
				fx, fy = fx+dx, fy+dy
				tiles = append(tiles, [2]int{int(fx), int(fy)})
			}
			return tiles
		}

		tileX, tileY := world.ScreenToCartesian(x, y)
		if tileX >= 0 && tileY >= 0 && tileX < 256 && tileY < 256 {
			multiUseStructure := world.World.HoverStructure == world.StructureBulldozer || world.World.HoverStructure == world.StructureRoad
			dragStarted := world.World.BuildDragX != -1 || world.World.BuildDragY != -1
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || (multiUseStructure && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)) || (multiUseStructure && dragStarted) {
				if !dragStarted && world.World.Funds >= world.StructureCosts[world.World.HoverStructure] {
					world.World.BuildDragX, world.World.BuildDragY = int(tileX), int(tileY)
				}

				if world.World.HoverStructure == world.StructureRoad {
					tiles := roadTiles(world.World.BuildDragX, world.World.BuildDragY, int(tileX), int(tileY))

					if dragStarted && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
						// TODO build all tiles
						world.World.Level.ClearHoverSprites()
						var builtRoad bool
						cost := world.StructureCosts[world.World.HoverStructure] * len(tiles) / 2
						if cost <= world.World.Funds {
							cost = 0
							for _, tile := range tiles {
								_, err := s.buildStructure(world.World.HoverStructure, tile[0], tile[1], !builtRoad)
								if err == nil {
									cost += world.StructureCosts[world.World.HoverStructure]
									builtRoad = true
								}
							}
							if cost > 0 {
								world.ShowBuildCost(world.World.HoverStructure, cost)
							}
						} else {
							world.ShowMessage("Insufficient funds", 3)
						}

						world.World.BuildDragX, world.World.BuildDragY = -1, -1
						dragStarted = false
					} else {
						// TODO draw hover sprites
						// TODO move below into shared func
						world.World.Level.ClearHoverSprites()
						world.BuildStructure(world.World.HoverStructure, true, int(tileX), int(tileY), false)
						var cost int
						for _, tile := range tiles {
							world.BuildStructure(world.World.HoverStructure, true, tile[0], tile[1], false)
							cost += world.StructureCosts[world.World.HoverStructure]
						}
						world.World.HoverValid = cost <= world.World.Funds
					}
					return nil
				} else if dragStarted && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
					world.World.BuildDragX, world.World.BuildDragY = -1, -1
				}

				cost := world.StructureCosts[world.World.HoverStructure]
				if world.World.Funds < cost {
					world.ShowMessage("Insufficient funds", 3)
				} else {
					world.World.Level.ClearHoverSprites()

					// TODO draw hovers and build all roads in a line from drag start
					structure, err := s.buildStructure(world.World.HoverStructure, int(tileX), int(tileY), true)
					if err == nil {
						tileX, tileY = float64(structure.X), float64(structure.Y)
						world.ShowBuildCost(world.World.HoverStructure, cost)
						if !multiUseStructure {
							world.World.HoverStructure = 0
							world.World.BuildDragX, world.World.BuildDragY = -1, -1
							world.World.LastBuildX, world.World.LastBuildY = -1, -1
						}
					}

					if world.World.HoverStructure > 0 {
						world.BuildStructure(world.World.HoverStructure, true, int(tileX), int(tileY), false)
					}
				}
			} else {
				world.World.Level.ClearHoverSprites()

				world.BuildStructure(world.World.HoverStructure, true, int(tileX), int(tileY), false)
			}
			world.World.HoverX, world.World.HoverY = int(tileX), int(tileY)
		}
	}

	return nil
}

func (s *playerMoveSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrUnregister
}

func deltaXY(x1, y1, x2, y2 float64) (dx float64, dy float64) {
	dx, dy = x1-x2, y1-y2
	if dx < 0 {
		dx *= -1
	}
	if dy < 0 {
		dy *= -1
	}
	return dx, dy
}
