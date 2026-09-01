package cplugin

import (
	"fmt"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// ResumeCommand handles /resume
func ResumeCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}
	chatID := m.Chat.ID

	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Sirf Admins he is command ko use kar sakte hain!**", tb.ModeMarkdown)
		return
	}

	if utils.IsMusicPlaying(chatID) {
		b.Send(m.Chat, "⚠️ **Stream is already playing.**", tb.ModeMarkdown)
		return
	}

	utils.MusicOn(chatID)
	utils.ResumeStream(chatID) // Calls pytgcalls wrapper

	// Dynamic Buttons
	markup := &tb.ReplyMarkup{}
	btnSkip := CreateBtn(markup, "sᴋɪᴘ", fmt.Sprintf("ADMIN Skip|%d", chatID), "", Primary, false)
	btnStop := CreateBtn(markup, "sᴛᴏᴘ", fmt.Sprintf("ADMIN Stop|%d", chatID), "", Danger, false)
	btnPause := CreateBtn(markup, "ᴘᴀᴜsᴇ", fmt.Sprintf("ADMIN Pause|%d", chatID), "", Success, false)

	markup.Inline(
		markup.Row(btnSkip, btnStop),
		markup.Row(btnPause),
	)

	text := fmt.Sprintf("▶️ **Stream Resumed**\n\n**Admin:** [%s](tg://user?id=%d)", m.Sender.FirstName, m.Sender.ID)
	b.Send(m.Chat, text, markup, tb.ModeMarkdown)
}
