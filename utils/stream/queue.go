package stream

import (
	"sync"
)

type QueueItem struct {
	Title      string
	Duration   string
	StreamType string
	By         string
	UserID     int64
	ChatID     int64
	File       string
	VidID      string
	Seconds    int
	Played     int
	Source     string
}

var (
	globalQueue = make(map[int64][]QueueItem)
	queueMutex  sync.RWMutex
)

// PutQueue adds a track to the queue
func PutQueue(chatID, origChatID int64, file, title, duration, user, vidID string, userID int64, stream string, forcePlay bool, source string) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	item := QueueItem{
		Title:      title,
		Duration:   duration,
		StreamType: stream,
		By:         user,
		UserID:     userID,
		ChatID:     origChatID,
		File:       file,
		VidID:      vidID,
		Source:     source,
		Seconds:    0, // Convert duration string to seconds here in prod
	}

	queue := globalQueue[chatID]
	if forcePlay {
		globalQueue[chatID] = append([]QueueItem{item}, queue...)
	} else {
		globalQueue[chatID] = append(queue, item)
	}
}

// GetQueue retrieves the queue
func GetQueue(chatID int64) []QueueItem {
	queueMutex.RLock()
	defer queueMutex.RUnlock()
	return globalQueue[chatID]
}

// PopQueue removes the first element
func PopQueue(chatID int64) *QueueItem {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	queue := globalQueue[chatID]
	if len(queue) == 0 {
		return nil
	}
	popped := queue[0]
	globalQueue[chatID] = queue[1:]
	return &popped
}
