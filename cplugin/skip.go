package cplugin

import (
	"fmt"
	"strconv"
	"strings"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// SkipCommand handles /skip, /next
func SkipCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}
	chatID := m.Chat.ID

	// BULLETPROOF ADMIN CHECK
	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **You don't have permission to use this command. Only Admins can skip.**", tb.ModeMarkdown)
		return
	}

	queue := utils.GetQueue(chatID)
	if len(queue) == 0 {
		b.Send(m.Chat, "⚠️ **The queue is empty!**", tb.ModeMarkdown)
		return
	}

	if utils.GetLoop(chatID) != 0 {
		b.Send(m.Chat, "⚠️ **Please disable loop to skip the track.**", tb.ModeMarkdown)
		return
	}

	skipCount := 1
	args := strings.Split(m.Text, " ")
	if len(args) > 1 {
		if val, err := strconv.Atoi(args[1]); err == nil {
			if val >= 1 && val <= len(queue) {
				skipCount = val
			} else {
				b.Send(m.Chat, fmt.Sprintf("⚠️ Provide a number between 1 and %d", len(queue)), tb.ModeMarkdown)
				return
			}
		}
	}

	// Pop tracks based on skip count
	for i := 0; i < skipCount-1; i++ {
		utils.PopQueue(chatID)
	}

	// Change Stream via Call module
	err := utils.SkipStream(chatID) // Wraps pytgcalls change_stream
	if err != nil {
		b.Send(m.Chat, "❌ Error skipping stream.")
		utils.StopStream(chatID)
		return
	}

	b.Send(m.Chat, fmt.Sprintf("➻ sᴛʀᴇᴀᴍ sᴋɪᴩᴩᴇᴅ 🎄\n└ʙʏ : [%s](tg://user?id=%d)", m.Sender.FirstName, m.Sender.ID), tb.ModeMarkdown)
}
