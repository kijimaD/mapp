package system

import (
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/mapp/world"
)

const CameraMoveSpeed = 0.132

type CameraSystem struct{}

func NewCameraSystem() *CameraSystem {
	s := &CameraSystem{}

	return s
}

func (s *CameraSystem) Update(e gohan.Entity) error {
	if !world.World.GameStarted || world.World.GameOver {
		return nil
	}

	// TODO
	return nil
}

func (_ *CameraSystem) Draw(_ gohan.Entity, screen *ebiten.Image) error {
	return gohan.ErrUnregister
}
