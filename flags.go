//go:build !js || !wasm
// +build !js !wasm

package main

import (
	"flag"

	"github.com/kijimaD/mapp/world"
)

func parseFlags() {
	flag.Parse()
	world.StartGame()
}
