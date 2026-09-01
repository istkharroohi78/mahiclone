package admins

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

// AdminCheck verifies if the user is an admin or a Telegram Service Bot
func AdminCheck(b *tb.Bot, m *tb.Message) bool {
	if m.Sender == nil {
		return false
	}

	if m.Chat.Type != tb.ChatSuperGroup && m.Chat.Type != tb.ChatChannel {
		return false
	}

	// 777000 = Telegram Service Notifications
	// 1087968824 = GroupAnonymousBot
	if m.Sender.ID == 777000 || m.Sender.ID == 1087968824 {
		return true
	}

	member, err := b.ChatMemberOf(m.Chat, m.Sender)
	if err != nil {
		return false
	}

	if member.Role == tb.Administrator || member.Role == tb.Creator {
		return true
	}
	return false
}
