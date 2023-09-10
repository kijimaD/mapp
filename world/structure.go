package world

import "code.rocketnine.space/tslocum/gohan"

const (
	StructureToggleHelp = iota + 1
	StructureBulldozer
	StructureRoad
	StructurePlain
	StationBusStop
)

var StructureFilePaths = map[int]string{
	StructureBulldozer: "map/bulldozer.tmx",
	StructureRoad:      "map/road.tmx",
	StructurePlain:     "map/plain.tmx",
	StationBusStop:     "map/busstop.tmx",
}

type Structure struct {
	Type int
	X, Y int

	Entity   gohan.Entity
	Children []gohan.Entity
}
