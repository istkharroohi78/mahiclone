package tools

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	enabledVCLogChats = make(map[int64]bool)
	vcLogMutex        sync.RWMutex
)

func RegisterVCLoggerHandlers(b *tb.Bot) {
	b.Handle("/vclogger", func(m *tb.Message) {
		args := strings.Split(m.Text, " ")
		if len(args) < 2 {
			vcLogMutex.RLock()
			status := enabledVCLogChats[m.Chat.ID]
			vcLogMutex.RUnlock()

			statusStr := "❌ OFF"
			if status {
				statusStr = "✅ ON"
			}
			b.Send(m.Chat, fmt.Sprintf("📊 **Voice Chat Logger Status:** %s\n\n**Usage:** `/vclogger on` or `/vclogger off`", statusStr), tb.ModeMarkdown)
			return
		}

		action := strings.ToLower(args[1])
		vcLogMutex.Lock()
		defer vcLogMutex.Unlock()

		if action == "on" {
			enabledVCLogChats[m.Chat.ID] = true
			b.Send(m.Chat, "✅ **Voice Chat Join/Leave Logger Enabled.**", tb.ModeMarkdown)
		} else if action == "off" {
			delete(enabledVCLogChats, m.Chat.ID)
			b.Send(m.Chat, "❌ **Voice Chat Logger Disabled.**", tb.ModeMarkdown)
		}
	})
}

// SendVCJoinNotification posts an alert when user joins voice chat
func SendVCJoinNotification(b *tb.Bot, chatID, userID int64, userName string) {
	vcLogMutex.RLock()
	enabled := enabledVCLogChats[chatID]
	vcLogMutex.RUnlock()

	if !enabled {
		return
	}

	chat, _ := b.ChatByID(fmt.Sprintf("%d", chatID))
	text := fmt.Sprintf(`<b>#JoinVoiceChat</b>

<b>● ɴᴀᴍᴇ ➛</b> <a href="tg://user?id=%d">%s</a>
<b>● ɪᴅ ➛</b> <code>%d</code>`, userID, userName, userID)

	msg, err := b.Send(chat, text, tb.ModeHTML)
	if err == nil {
		go func() {
			time.Sleep(7 * time.Second)
			b.Delete(msg)
		}()
	}
}
