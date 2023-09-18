package world

import "code.rocketnine.space/tslocum/gohan"

const (
	StructureToggleHelp = iota + 1
	StructureBulldozer
	StructureRoad
	StructurePlain
	StationBusStop
)

// タイルへのファイルパス
var StructureFilePaths = map[int]string{
	StructureBulldozer: "map/bulldozer.tmx",
	StructureRoad:      "map/road.tmx",
	StructurePlain:     "map/plain.tmx",
	StationBusStop:     "map/busstop.tmx",
}

// ツールチップの文字列
var StructureTooltips = map[int]string{
	StructureToggleHelp: "Help",
	StructureBulldozer:  "Bulldozer",
	StructureRoad:       "Road",
	StationBusStop:      "BusStop",
}

// 実行に必要な額
var StructureCosts = map[int]int{
	StructureBulldozer: 5,
	StructureRoad:      25,
	StationBusStop:     50,
}

type Structure struct {
	Type int
	X, Y int

	Entity   gohan.Entity
	Children []gohan.Entity
}
