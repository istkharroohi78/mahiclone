package cplugin

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// IDCommand handles /id
func IDCommand(b *tb.Bot, m *tb.Message) {
	text := fmt.Sprintf("**[ᴍᴇssᴀɢᴇ ɪᴅ:](https://t.me/c/%d/%d)** `%d`\n", m.Chat.ID, m.ID, m.ID)
	text += fmt.Sprintf("**[ʏᴏᴜʀ ɪᴅ:](tg://user?id=%d)** `%d`\n", m.Sender.ID, m.Sender.ID)

	args := strings.SplitN(m.Text, " ", 2)
	if len(args) == 2 {
		targetChat, err := b.ChatByID(args[1])
		if err == nil {
			text += fmt.Sprintf("**[ᴜsᴇʀ ɪᴅ:](tg://user?id=%d)** `%d`\n", targetChat.ID, targetChat.ID)
		} else {
			b.Send(m.Chat, "ᴛʜɪs ᴜsᴇʀ ᴅᴏᴇsɴ'ᴛ ᴇxɪsᴛ.", &tb.SendOptions{ReplyTo: m})
			return
		}
	}

	if m.Chat.Username != "" {
		text += fmt.Sprintf("**[ᴄʜᴀᴛ ɪᴅ:](https://t.me/%s)** `%d`\n\n", m.Chat.Username, m.Chat.ID)
	} else {
		text += fmt.Sprintf("**[ᴄʜᴀᴛ ɪᴅ:]** `%d`\n\n", m.Chat.ID)
	}

	if m.ReplyTo != nil {
		if m.ReplyTo.Sender != nil {
			text += fmt.Sprintf("**[ʀᴇᴘʟɪᴇᴅ ᴍᴇssᴀɢᴇ ɪᴅ:](https://t.me/c/%d/%d)** `%d`\n", m.Chat.ID, m.ReplyTo.ID, m.ReplyTo.ID)
			text += fmt.Sprintf("**[ʀᴇᴘʟɪᴇᴅ ᴜsᴇʀ ɪᴅ:](tg://user?id=%d)** `%d`\n\n", m.ReplyTo.Sender.ID, m.ReplyTo.Sender.ID)
		}

		if m.ReplyTo.OriginalChat != nil {
			text += fmt.Sprintf("ᴛʜᴇ ғᴏʀᴡᴀʀᴅᴇᴅ ᴄʜᴀɴɴᴇʟ, %s, ʜᴀs ᴀɴ ɪᴅ ᴏғ `%d`\n\n", m.ReplyTo.OriginalChat.Title, m.ReplyTo.OriginalChat.ID)
		}

		// SenderChat logic applies to anonymous admins / channels
		if m.ReplyTo.Sender != nil && m.ReplyTo.Sender.IsBot && m.ReplyTo.Sender.ID == m.Chat.ID {
			text += fmt.Sprintf("ɪᴅ ᴏғ ᴛʜᴇ ʀᴇᴘʟɪᴇᴅ ᴄʜᴀᴛ/ᴄʜᴀɴɴᴇʟ, ɪs `%d`", m.Chat.ID)
		}
	}

	b.Send(m.Chat, text, tb.ModeMarkdown, &tb.SendOptions{DisableWebPagePreview: true})
}
