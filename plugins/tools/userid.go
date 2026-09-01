package tools

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterUserIDHandlers(b *tb.Bot) {
	b.Handle("/id", func(m *tb.Message) {
		text := fmt.Sprintf("<b>[ʏᴏᴜʀ ɪᴅ:]</b> <code>%d</code>\n", m.Sender.ID)
		text += fmt.Sprintf("<b>[ᴄʜᴀᴛ ɪᴅ:]</b> <code>%d</code>\n", m.Chat.ID)
		text += fmt.Sprintf("<b>[ᴍᴇssᴀɢᴇ ɪᴅ:]</b> <code>%d</code>\n", m.ID)

		if m.ReplyTo != nil {
			text += "\n<b>--- ʀᴇᴘʟɪᴇᴅ ᴍᴇssᴀɢᴇ ---</b>\n"
			text += fmt.Sprintf("<b>[ʀᴇᴘʟɪᴇᴅ ᴜsᴇʀ ɪᴅ:]</b> <code>%d</code>\n", m.ReplyTo.Sender.ID)
			text += fmt.Sprintf("<b>[ʀᴇᴘʟɪᴇᴅ ᴍᴇssᴀɢᴇ ɪᴅ:]</b> <code>%d</code>\n", m.ReplyTo.ID)
		}

		args := strings.Split(m.Text, " ")
		if len(args) > 1 {
			targetChat, err := b.ChatByID(args[1])
			if err == nil {
				text += fmt.Sprintf("\n<b>[ǫᴜᴇʀɪᴇᴅ ɪᴅ:]</b> <code>%d</code> (%s)", targetChat.ID, targetChat.Title)
			}
		}

		b.Send(m.Chat, text, tb.ModeHTML)
	})
}
