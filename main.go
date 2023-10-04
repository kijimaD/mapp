package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/mapp/scene"
	"github.com/kijimaD/mapp/world"
	"github.com/sedyh/mizu/pkg/engine"
)

func main() {
	ebiten.SetWindowTitle("Mapp")
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetMaxTPS(144)
	ebiten.SetRunnableOnUnfocused(true) // Note - this currently does nothing in ebiten
	ebiten.SetWindowClosingHandled(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)

	scene := scene.NewScene()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM)
	go func() {
		<-sigc

		scene.Exit()
	}()

	world.StartGame()
	err := ebiten.RunGame(NewGame(engine.NewGame(scene)))
	if err != nil {
		log.Fatal(err)
	}
}

// Layout()の結果をebitenから取り出すために定義している
// https://github.com/sedyh/mizu/issues/8#issuecomment-1528772092
type Game struct {
	e ebiten.Game
}

func NewGame(e ebiten.Game) *Game {
	return &Game{e}
}

func (g *Game) Update() error {
	return g.e.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.e.Draw(screen)
}

func (g *Game) Layout(w, h int) (int, int) {
	world.World.ScreenW = w
	world.World.ScreenH = h
	return w, h
}
