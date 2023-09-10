package world

import "code.rocketnine.space/tslocum/gohan"

const (
	StructureToggleHelp = iota + 1
	StructureBulldozer
	StructureRoad
	StructureResidentialLow
	StructureResidentialMedium
	StructureResidentialHigh
	StructureCommercialLow
	StructureCommercialMedium
	StructureCommercialHigh
	StructureIndustrialLow
	StructureIndustrialMedium
	StructureIndustrialHigh
	StructurePoliceStation
	StructurePowerPlantCoal
	StructurePowerPlantSolar
)

var StructureFilePaths = map[int]string{
	StructureBulldozer:         "map/bulldozer.tmx",
	StructureRoad:              "map/road.tmx",
	StructureResidentialLow:    "map/residential_low1.tmx",
	StructureResidentialMedium: "map/residential_med1.tmx",
	StructureResidentialHigh:   "map/residential_high1.tmx",
	StructureCommercialLow:     "map/commercial_low1.tmx",
	StructureCommercialMedium:  "map/commercial_med1.tmx",
	StructureCommercialHigh:    "map/commercial_high1.tmx",
	StructureIndustrialLow:     "map/industrial_low1.tmx",
	StructureIndustrialMedium:  "map/industrial_med1.tmx",
	StructureIndustrialHigh:    "map/industrial_high1.tmx",
	StructurePoliceStation:     "map/policestation.tmx",
	StructurePowerPlantCoal:    "map/power_coal.tmx",
	StructurePowerPlantSolar:   "map/power_solar.tmx",
}

type Structure struct {
	Type int
	X, Y int

	Entity   gohan.Entity
	Children []gohan.Entity
}
