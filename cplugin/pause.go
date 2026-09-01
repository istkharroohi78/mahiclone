package cplugin

import (
	"fmt"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// PauseCommand handles /pause and /cpause
func PauseCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	chatID := m.Chat.ID

	// BULLETPROOF ADMIN CHECK
	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Sirf Admins he is command ko use kar sakte hain!**", tb.ModeMarkdown)
		return
	}

	// CHECK IF MUSIC IS PLAYING
	if !utils.IsMusicPlaying(chatID) {
		b.Send(m.Chat, "⚠️ **Nothing is currently playing.**", tb.ModeMarkdown)
		return
	}

	// Database off & stream pause
	utils.MusicOff(chatID)
	utils.PauseStream(chatID) // Calls pytgcalls wrapper

	// Inline Buttons setup leveraging the Premium Emojis and Button Styles defined previously
	markup := &tb.ReplyMarkup{}

	// CreateBtn arguments: (markup, text, callback_data, url, style, no_emoji)
	btnResume := CreateBtn(markup, "ʀᴇsᴜᴍᴇ ▷", fmt.Sprintf("ADMIN Resume|%d", chatID), "", Success, false)
	btnReplay := CreateBtn(markup, "ʀᴇᴘʟᴀʏ ↺", fmt.Sprintf("ADMIN Replay|%d", chatID), "", Primary, false)

	// Uses use_emoji=True equivalent logic in Go
	btnClone := CreateBtn(markup, "✯ CLONE NOW ✯", "", "https://t.me/clone_MUSICrobot", Primary, false)

	markup.Inline(
		markup.Row(btnResume, btnReplay),
		markup.Row(btnClone),
	)

	text := fmt.Sprintf("⏸ **Stream Paused**\n\n**Admin:** [%s](tg://user?id=%d)", m.Sender.FirstName, m.Sender.ID)
	b.Send(m.Chat, text, markup, tb.ModeMarkdown)
}
