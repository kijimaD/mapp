package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kijimaD/mapp/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/mapp/game"
)

func main() {
	ebiten.SetWindowTitle("Mapp")
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetMaxTPS(144)
	ebiten.SetRunnableOnUnfocused(true) // Note - this currently does nothing in ebiten
	ebiten.SetWindowClosingHandled(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)

	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM)
	go func() {
		<-sigc

		g.Exit()
	}()

	world.StartGame()

	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
