package system

import (
	"os"
	"runtime/pprof"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type profileSystem struct {
	cpuProfile *os.File
}

func NewProfileSystem() *profileSystem {
	return &profileSystem{}
}

func (s *profileSystem) Update(_ gohan.Entity) error {
	if ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if s.cpuProfile == nil {
			cpuProfile, err := os.Create("mapp.prof")
			s.cpuProfile = cpuProfile
			if err != nil {
				return err
			}

			err = pprof.StartCPUProfile(s.cpuProfile)
			if err != nil {
				return err
			}
		} else {
			pprof.StopCPUProfile()

			s.cpuProfile.Close()
			s.cpuProfile = nil
		}
	}
	return nil
}

func (s *profileSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrUnregister
}
