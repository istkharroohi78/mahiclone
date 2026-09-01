package misc

import (
	"strconv"
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
			activeChats := database.GetActiveChats() // This actually returns []string
			for _, chatID := range activeChats {
				
				// Convert chatID from string to int64
				chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
				if err != nil {
					continue
				}

				if !database.IsMusicPlaying(chatIDInt) {
					continue
				}

				queue := stream.GetQueue(chatIDInt)
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
				// The compiler expects 0 arguments for this function: want ()
				stream.IncrementPlayed()
			}
		}
	}()
}
