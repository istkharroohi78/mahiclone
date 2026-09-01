package sudo

import (
	"os"

	"ANJALI/utils/inline"

	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterRestartHandlers(b *tb.Bot) {
	b.Handle("/restart", func(m *tb.Message) {
		if !IsSudoer(int64(m.Sender.ID)) {
			return
		}

		markup := &tb.ReplyMarkup{}
		markup.Inline(
			markup.Row(
				inline.CreateBtn(markup, "🤖 ᴍᴀɪɴ ʙᴏᴛ", "restart_main", "", inline.Primary, 0, false),
				inline.CreateBtn(markup, "👥 ᴄʟᴏɴᴇs ʙᴏᴛ", "restart_clones", "", inline.Primary, 0, false),
			),
			markup.Row(inline.CreateBtn(markup, "🔄 ʙᴏᴛʜ (ᴍᴀɪɴ + ᴄʟᴏɴᴇs)", "restart_both", "", inline.Primary, 0, false)),
			markup.Row(inline.CreateBtn(markup, "❌ ᴄᴀɴᴄᴇʟ", "cancel_restart", "", inline.Danger, 0, false)),
		)

		b.Send(m.Chat, "**⚠️ ᴋɪsᴇ ʀᴇsᴛᴀʀᴛ ᴋᴀʀɴᴀ ᴄʜᴀʜᴛᴇ ʜᴀɪɴ?**\n\n*(Neeche diye gaye options me se select karein)*", markup, tb.ModeMarkdown)
	})

	b.Handle("\frestart_main", func(c *tb.Callback) {
		if !IsSudoer(int64(c.Sender.ID)) {
			return
		}
		b.Edit(c.Message, "🔄 **ʀᴇsᴛᴀʀᴛɪɴɢ ᴍᴀɪɴ ʙᴏᴛ...**")
		os.Exit(0) // Assuming a supervisor script restarts the process
	})

	b.Handle("\fcancel_restart", func(c *tb.Callback) {
		if !IsSudoer(int64(c.Sender.ID)) {
			return
		}
		b.Edit(c.Message, "**❌ ʀᴇsᴛᴀʀᴛ ᴘʀᴏᴄᴇss ᴄᴀɴᴄᴇʟʟᴇᴅ.**")
	})
}
