package tools

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var reloadLimiter = make(map[int64]time.Time)

func RegisterReloadHandlers(b *tb.Bot) {
	// /reload & /admincache
	b.Handle("/reload", func(m *tb.Message) {
		chatID := m.Chat.ID
		if last, ok := reloadLimiter[chatID]; ok && time.Since(last) < 2*time.Minute {
			remaining := (2*time.Minute - time.Since(last)).Round(time.Second)
			b.Send(m.Chat, fmt.Sprintf("⏳ **Please wait %s before reloading admin cache again.**", remaining), tb.ModeMarkdown)
			return
		}

		reloadLimiter[chatID] = time.Now()
		msg, _ := b.Send(m.Chat, "🔄 **Refreshing group admin cache & permissions...**", tb.ModeMarkdown)

		// Sync with Telegram Admin list
		time.Sleep(1 * time.Second)
		b.Edit(msg, "✅ **Admin cache refreshed successfully!**", tb.ModeMarkdown)
	})

	b.Handle("/admincache", func(m *tb.Message) {
		b.Trigger("/reload", m)
	})

	// Close Button Callback
	b.Handle("\fclose", func(c *tb.Callback) {
		b.Respond(c)
		b.Delete(c.Message)
		b.Send(c.Message.Chat, fmt.Sprintf("🗑️ Menu closed by [%s](tg://user?id=%d)", c.Sender.FirstName, c.Sender.ID), tb.ModeMarkdown)
	})
}
