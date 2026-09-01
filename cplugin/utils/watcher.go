package utils

import (
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

// HandleVideoChatEvent listens for video_chat_started and video_chat_ended events.
// Ise main bot file me Telebot handlers (jaise tb.OnVideoChatStarted) ke sath connect karna hoga.
func HandleVideoChatEvent(m *tb.Message, stopStreamForce func(int64)) {
	chatID := m.Chat.ID

	// 1. Stream ko force stop karega (Call manager ka function pass karna hoga)
	if stopStreamForce != nil {
		stopStreamForce(chatID)
	}

	// 2. Queue ko manually clear kar rahe hain taaki bot fresh start kare
	DBMutex.Lock()
	defer DBMutex.Unlock()

	// Chat ID ki entry ko empty slice se replace kar diya (Queue Cleared)
	DB[chatID] = []QueueItem{}

	log.Printf("Video chat event trigger hua. Queue cleared for ChatID: %d", chatID)
}
