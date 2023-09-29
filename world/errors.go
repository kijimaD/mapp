package world

import "errors"

var ErrNothingToBulldoze = errors.New("nothing to bulldoze")                         // 取り壊すものがない
var ErrInvalidBuildingNotFit = errors.New("invalid location: building does not fit") // 建設場所が不適
var ErrLocationOccupied = errors.New("invalid location: space already occupied")
var ErrTileImageNotFound = errors.New("tile image not found")
