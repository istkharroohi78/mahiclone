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

			// GetActiveChats directly returns []int64 now, so no conversion needed!
			activeChats := database.GetActiveChats() 
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

				// FIXED: Removed 'chatID' argument to match stream package
				stream.IncrementPlayed()
			}
		}
	}()
}
