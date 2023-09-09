//go:build !js || !wasm
// +build !js !wasm

package main

import (
	"flag"

	"code.rocketnine.space/tslocum/citylimits/world"
)

func parseFlags() {
	flag.Parse()
	world.StartGame()
}
