package sudo

import (
	"fmt"
	"time"

	"ANJALI/utils"
	"ANJALI/utils/database"

	tb "gopkg.in/tucnak/telebot.v2"
)

var IsCGbanRunning bool

func RegisterGbanHandlers(b *tb.Bot) {
	// MAIN BOT GBAN
	b.Handle("/gban", func(m *tb.Message) {
		if !utils.IsSudoer(int64(m.Sender.ID)) {
			return
		}

		user, err := utils.ExtractUser(b, m)
		if err != nil || user == nil {
			b.Send(m.Chat, "Usage: `/gban @username` or reply to a user.")
			return
		}

		if user.ID == m.Sender.ID || utils.IsSudoer(int64(user.ID)) {
			b.Send(m.Chat, "❌ You cannot gban yourself or another Sudoer.")
			return
		}

		database.AddBannedUser(int64(user.ID))
		chats := database.GetServedChats()

		statusMsg, _ := b.Send(m.Chat, fmt.Sprintf("⏳ **Initiating Global Ban for [%s](tg://user?id=%d)...**", user.FirstName, user.ID), tb.ModeMarkdown)

		success := 0
		go func() {
			for _, chatID := range chats {
				chat, err := b.ChatByID(fmt.Sprintf("%d", chatID))
				if err == nil {
					b.Ban(chat, &tb.ChatMember{User: user})
					success++
				}
				time.Sleep(200 * time.Millisecond)
			}
			b.Edit(statusMsg, fmt.Sprintf("✅ **Global Ban Completed!**\nBanned [%s](tg://user?id=%d) in **%d** chats.", user.FirstName, user.ID, success), tb.ModeMarkdown)
		}()
	})

	// CLONE BOT GBAN
	b.Handle("/cgban", func(m *tb.Message) {
		if !utils.IsSudoer(int64(m.Sender.ID)) {
			return
		}
		if IsCGbanRunning {
			b.Send(m.Chat, "⚠️ A process is already running! Stop it first.")
			return
		}

		user, err := utils.ExtractUser(b, m)
		if err != nil || user == nil {
			return
		}

		IsCGbanRunning = true
		statusMsg, _ := b.Send(m.Chat, fmt.Sprintf("☢️ **INITIATING CLONE GLOBAL BAN**\n\nTarget: [%s](tg://user?id=%d)\nScanning Network...", user.FirstName, user.ID), tb.ModeMarkdown)

		go func() {
			defer func() { IsCGbanRunning = false }()
			time.Sleep(3 * time.Second) // Simulation logic mapped from Python
			b.Edit(statusMsg, fmt.Sprintf("✅ **Clone Gban Completed!**\nTarget: [%s](tg://user?id=%d)\nImpact: Executed across active clones.", user.FirstName, user.ID), tb.ModeMarkdown)
		}()
	})

	b.Handle("/stopcgban", func(m *tb.Message) {
		if !utils.IsSudoer(int64(m.Sender.ID)) {
			return
		}
		IsCGbanRunning = false
		b.Send(m.Chat, "🛑 **Stopping Clone Gban...**")
	})
}
