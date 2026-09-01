package cplugin

import (
	"fmt"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// StopCommand handles /stop, /end, /cstop, /cend
func StopCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	chatID := m.Chat.ID

	// BULLETPROOF ADMIN CHECK
	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Sirf Admins he is command ko use kar sakte hain!**", tb.ModeMarkdown)
		return
	}

	// Stop Stream
	utils.StopStream(chatID)

	// Reset Loop
	utils.SetLoop(chatID, 0)

	// Clear Queue
	utils.ClearChatQueue(chatID)

	text := fmt.Sprintf("⏹ **Stream Ended/Stopped**\n\n**Admin:** [%s](tg://user?id=%d)", m.Sender.FirstName, m.Sender.ID)
	b.Send(m.Chat, text, utils.CloseMarkup(), tb.ModeMarkdown)
}
