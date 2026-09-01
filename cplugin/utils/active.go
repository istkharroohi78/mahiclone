package utils

import (
	"log"
	"sync"
)

var (
	// Maps provide O(1) instant lookups compared to iterating through slices/lists
	activeChats = make(map[int64]bool)
	streamState = make(map[int64]bool)

	// mu ensures thread safety when multiple goroutines/clones access these maps
	mu sync.RWMutex
)

// IsActiveChat checks if a chat is currently active
func IsActiveChat(chatID int64) bool {
	mu.RLock() // Read lock
	defer mu.RUnlock()
	return activeChats[chatID]
}

// AddActiveChat adds a chat to the active list
func AddActiveChat(chatID int64) {
	mu.Lock() // Write lock
	defer mu.Unlock()
	activeChats[chatID] = true
}

// RemoveActiveChat removes a chat from the active list
func RemoveActiveChat(chatID int64) {
	mu.Lock()
	defer mu.Unlock()
	delete(activeChats, chatID)
}

// GetActiveChats returns a slice of all active chat IDs
func GetActiveChats() []int64 {
	mu.RLock()
	defer mu.RUnlock()

	var chats []int64
	for chatID := range activeChats {
		chats = append(chats, chatID)
	}
	return chats
}

// IsStreaming checks if a chat is currently streaming
// Note: Covers both is_streaming and iss_streaming from Python
func IsStreaming(chatID int64) bool {
	mu.RLock()
	defer mu.RUnlock()
	return streamState[chatID]
}

// StreamOn sets the stream status to true
func StreamOn(chatID int64) {
	mu.Lock()
	defer mu.Unlock()
	streamState[chatID] = true
}

// StreamOff sets the stream status to false
func StreamOff(chatID int64) {
	mu.Lock()
	defer mu.Unlock()
	streamState[chatID] = false
}

// Clear cleans up the active state, stream state, and in-memory DB queue for a chat
func Clear(chatID int64) {
	mu.Lock()
	defer mu.Unlock()

	// Safe execution using deferred recovery (Go's equivalent to try-except)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error in _clear_: %v\n", r) // Error terminal mein dikhega
		}
	}()

	// ✅ FIX 1: db clear
	// (Note: Replace this with your actual in-memory queue clearing logic)
	// core.MemoryDB[chatID] = nil

	// ✅ FIX 2: Active list se remove karna
	delete(activeChats, chatID)

	// ✅ FIX 3: Stream status ko dictionary se hatana
	delete(streamState, chatID)
}
