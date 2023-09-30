package world

import (
	"image"
	"strings"
	"sync"
)

var messageLock = &sync.Mutex{}

// メッセージは、建設時などに右上に一時的に出るメッセージのこと
func TickMessages() {
	messageLock.Lock()
	defer messageLock.Unlock()

	var removed int
	for j := 0; j < len(World.MessagesTicks); j++ {
		i := j - removed // 削除された分短くなるのを考慮する

		// Ticksの中身が0になったものは、MessagesとMessagesTicksから削除していく。0以上の場合はデクリメントする
		if World.MessagesTicks[i] == 0 {
			// 前を削除
			// [古1, 古2, 新1, 新2]
			// [古2, 新1, 新2]
			World.Messages = append(World.Messages[:i], World.Messages[i+1:]...)
			World.MessagesTicks = append(World.MessagesTicks[:i], World.MessagesTicks[i+1:]...)
			removed++
			World.HUDUpdated = true
		} else if World.MessagesTicks[i] > 0 {
			World.MessagesTicks[i]--
		}
	}
}

func ShowMessage(message string, duration int) {
	messageLock.Lock()
	defer messageLock.Unlock()

	World.Messages = append(World.Messages, message)
	World.MessagesTicks = append(World.MessagesTicks, duration)

	World.HUDUpdated = true
}

func ShowBuildCost(structureType int, cost int) {
	if structureType == StructureBulldozer {
		ShowMessage(World.Printer.Sprintf("Bulldozed area (-$%d)", cost), 3)
	} else {
		ShowMessage(World.Printer.Sprintf("Built %s (-$%d)", strings.ToLower(StructureTooltips[World.HoverStructure]), cost), 3)
	}
}

// 指定座標に該当するボタンを返す
func HUDButtonAt(x, y int) *HUDButton {
	point := image.Point{x, y}
	for i, rect := range World.HUDButtonRects {
		if point.In(rect) {
			return HUDButtons[i]
		}
	}
	return nil
}

// 建設を選択中
func SetHoverStructure(structureType int) {
	World.HoverStructure = structureType
	World.HUDUpdated = true
}

// 選択中の建物のツールチップテキストを取得する
func TooltipText() string {
	tooltipText := StructureTooltips[World.HoverStructure]
	cost := StructureCosts[World.HoverStructure]
	if cost > 0 {
		tooltipText += World.Printer.Sprintf("\n$%d", cost)
	}
	return tooltipText
}
