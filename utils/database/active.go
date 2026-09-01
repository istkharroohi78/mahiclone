package database

import "fmt"

// ActiveChats stores currently active streams in memory
var ActiveChats = make(map[string]bool)

// AddActiveChat adds a chat to the active stream list
func AddActiveChat(chatID int64) {
	cid := fmt.Sprintf("%d", chatID)
	ActiveChats[cid] = true
}

// RemoveActiveChat removes a chat from the active stream list
func RemoveActiveChat(chatID int64) {
	cid := fmt.Sprintf("%d", chatID)
	delete(ActiveChats, cid)
}

// IsActiveChat checks if a chat is currently streaming
func IsActiveChat(chatID int64) bool {
	cid := fmt.Sprintf("%d", chatID)
	return ActiveChats[cid]
}
