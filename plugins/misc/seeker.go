package misc

import (
	"time"

	"ANJALI/utils/database"
	"ANJALI/utils/stream"
)

// StartSeekerLoop starts a background loop to track stream progress
func StartSeekerLoop() {
	go func() {
		for {
			time.Sleep(1 * time.Second)

			// Logic translated: Update 'played' seconds in DB/Memory
			activeChats := database.GetActiveChats() // Assuming this returns []int64
			for _, chatID := range activeChats {
				if !database.IsMusicPlaying(chatID) {
					continue
				}

				queue := stream.GetQueue(chatID)
				if len(queue) == 0 {
					continue
				}

				duration := queue[0].Seconds
				if duration == 0 {
					continue
				}

				if queue[0].Played >= duration {
					continue
				}

				// Increment played seconds safely via stream package setter
				stream.IncrementPlayed(chatID)
			}
		}
	}()
}
