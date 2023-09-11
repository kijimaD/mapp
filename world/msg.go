package world

import "sync"

var messageLock = &sync.Mutex{}

const messageDuration = 144 * 3

func TickMessages() {
	messageLock.Lock()
	defer messageLock.Unlock()

	var removed int
	for j := 0; j < len(World.MessagesTicks); j++ {
		i := j - removed
		if World.MessagesTicks[i] == 0 {
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
