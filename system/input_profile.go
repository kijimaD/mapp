package system

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sedyh/mizu/pkg/engine"
)

type profileSystem struct {
	cpuProfile *os.File
}

func NewProfileSystem() *profileSystem {
	return &profileSystem{}
}

func (s *profileSystem) Update(w engine.World) {
	if ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if s.cpuProfile == nil {
			fmt.Println("profile active")
			cpuProfile, err := os.Create("mapp.prof")
			s.cpuProfile = cpuProfile
			if err != nil {
				return
			}

			err = pprof.StartCPUProfile(s.cpuProfile)
			if err != nil {
				return
			}
		} else {
			pprof.StopCPUProfile()

			s.cpuProfile.Close()
			s.cpuProfile = nil
		}
	}
	return
}
