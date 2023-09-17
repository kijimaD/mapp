package world

import (
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
