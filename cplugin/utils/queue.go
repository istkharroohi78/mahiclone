package utils

import (
	"sync"
)

// QueueItem stores all the metadata for a single track in the queue
type QueueItem struct {
	Title      string
	Dur        string
	File       string
	VidID      string
	By         string
	UserID     int64
	StreamType string
	Played     int
	Seconds    int
}

var (
	// DBMutex ensures thread-safe operations on the DB map
	DBMutex sync.Mutex
	// DB acts as the global in-memory database mapping Chat IDs to their Queues
	DB = make(map[int64][]QueueItem)
)

// Put adds a new track to the chat's queue
func Put(chatID int64, title, duration, vidID, filePath, rUser string, userID int64, streamType string) {
	if streamType == "" {
		streamType = "audio" // Default fallback
	}

	putF := QueueItem{
		Title:      title,
		Dur:        duration,
		File:       filePath,
		VidID:      vidID,
		By:         rUser,
		UserID:     userID,
		StreamType: streamType,
		Played:     0,
		Seconds:    0,
	}

	DBMutex.Lock()
	defer DBMutex.Unlock()

	// Go's append automatically handles empty or uninitialized slices
	DB[chatID] = append(DB[chatID], putF)
}
