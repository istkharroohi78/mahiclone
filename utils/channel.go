package utils

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

// GetChannelPlayCB fetches channel play mode settings
func GetChannelPlayCB(b *tb.Bot, command string, chatID int64) (int64, string) {
	if command == "c" {
		targetChatID := GetCMode(chatID) // Defined in database.go
		if targetChatID == 0 {
			return 0, ""
		}

		chat, err := b.ChatByID(fmt.Sprintf("%d", targetChatID))
		if err != nil {
			return 0, ""
		}
		return targetChatID, chat.Title
	}

	return chatID, ""
}
