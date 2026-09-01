package tools

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"ANJALI/config"

	tb "gopkg.in/tucnak/telebot.v2"
)

var MaliciousPatterns = []string{
	`webhook\.site`,
	`requestbin\.com`,
	`ngrok\.io`,
	`localhost`,
	`127\.0\.0\.1`,
}

func isMaliciousQuery(text string) bool {
	for _, pattern := range MaliciousPatterns {
		if matched, _ := regexp.MatchString("(?i)"+pattern, text); matched {
			return true
		}
	}
	return false
}

func RegisterSecurityHandlers(b *tb.Bot) {
	// Middleware inspection on commands
	b.Handle(tb.OnText, func(m *tb.Message) {
		if !strings.HasPrefix(m.Text, "/") {
			return
		}

		if isMaliciousQuery(m.Text) {
			b.Delete(m)
			warnMsg, _ := b.Send(m.Chat, "⚠️ **Security Alert: Malicious request blocked and reported.**", tb.ModeMarkdown)

			// Send Report to Logger Group
			cfg := config.LoadConfig()
			if cfg.LoggerID != 0 {
				logChat, err := b.ChatByID(fmt.Sprintf("%d", cfg.LoggerID))
				if err == nil {
					logText := fmt.Sprintf(`🚨 **MALICIOUS ATTEMPT DETECTED**

👤 **User:** [%s](tg://user?id=%d)
🆔 **ID:** `+"`%d`"+`
👥 **Chat:** %s (`+"`%d`"+`)
💬 **Query:** `+"`%s`", m.Sender.FirstName, m.Sender.ID, m.Sender.ID, m.Chat.Title, m.Chat.ID, m.Text)

					markup := &tb.ReplyMarkup{}
					markup.Inline(
						markup.Row(
							markup.Data("🚫 Block User", fmt.Sprintf("block_user_%d", m.Sender.ID)),
							markup.Data("🛑 Block Chat", fmt.Sprintf("block_chat_%d", m.Chat.ID)),
						),
					)
					b.Send(logChat, logText, markup, tb.ModeMarkdown)
				}
			}

			go func(msg *tb.Message) {
				time.Sleep(30 * time.Second)
				b.Delete(msg)
			}(warnMsg)
		}
	})
}
