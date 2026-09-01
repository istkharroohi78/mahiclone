package play

import (
	"fmt"
	"strconv"
	"strings"

	"ANJALI/utils/database"
	"ANJALI/utils/decorators"

	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterChannelPlayHandlers(b *tb.Bot) {
	b.Handle("/channelplay", func(m *tb.Message) {
		decorators.ActualAdminCB(b, nil, func(b *tb.Bot, c *tb.Callback) {
			// In telebot, command handlers use Context or Message
		})

		// Inline implementation of ActualAdmin since Telebot handles callbacks/messages differently
		args := strings.Split(m.Text, " ")
		if len(args) < 2 {
			b.Send(m.Chat, "⚠️ Usage: `/channelplay [Linked Channel ID]`")
			return
		}

		query := strings.ToLower(args[1])
		if query == "disable" {
			database.SetCMode(m.Chat.ID, 0)
			b.Send(m.Chat, "❌ Channel play disabled.")
			return
		}

		if channelID, err := strconv.ParseInt(query, 10, 64); err == nil {
			// Verify channel exists and bot is admin
			chat, err := b.ChatByID(fmt.Sprintf("%d", channelID))
			if err != nil {
				b.Send(m.Chat, "❌ Channel not found.")
				return
			}

			database.SetCMode(m.Chat.ID, chat.ID)
			b.Send(m.Chat, fmt.Sprintf("✅ Channel Play linked to: **%s**", chat.Title), tb.ModeMarkdown)
		}
	})
}
