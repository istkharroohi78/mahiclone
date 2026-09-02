package decorators

import (
	"fmt"

	"ANJALI/config"
	"ANJALI/utils"
	"ANJALI/utils/database"

	tb "gopkg.in/tucnak/telebot.v2"
)

// AdminRightsCheck wraps commands to ensure only authorized users or admins execute them
func AdminRightsCheck(b *tb.Bot, m *tb.Message, next func(b *tb.Bot, m *tb.Message, chatID int64)) {
	if database.IsMaintenance() && !utils.IsSudoer(int64(m.Sender.ID)) {
		b.Send(m.Chat, fmt.Sprintf("%s ɪs ᴜɴᴅᴇʀ ᴍᴀɪɴᴛᴇɴᴀɴᴄᴇ, ᴠɪsɪᴛ [sᴜᴘᴘᴏʀᴛ ᴄʜᴀᴛ](%s) ғᴏʀ ᴋɴᴏᴡɪɴɢ ᴛʜᴇ ʀᴇᴀsᴏɴ.", b.Me.FirstName, config.LoadConfig().SupportChat), tb.ModeMarkdown)
		return
	}

	b.Delete(m)

	var chatID int64 = m.Chat.ID
	// Channel Play check logic
	if len(m.Text) > 1 && m.Text[1] == 'c' {
		linkedChat := database.GetCMode(m.Chat.ID)
		if linkedChat == 0 {
			b.Send(m.Chat, "⚠️ Channel play is not configured.")
			return
		}
		chatID = linkedChat
	}

	// Active Chat Check (using the GetActiveChats array from database)
	isActive := false
	for _, id := range database.GetActiveChats() {
		if id == chatID {
			isActive = true
			break
		}
	}
	if !isActive {
		b.Send(m.Chat, "⚠️ No active stream in this chat.")
		return
	}

	if !database.IsNonAdminChat(m.Chat.ID) && !utils.IsSudoer(int64(m.Sender.ID)) {
		adminList := database.GetAdminCache(m.Chat.ID)
		isAdmin := false
		for _, adminID := range adminList {
			if adminID == int64(m.Sender.ID) {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			if database.IsSkipMode(m.Chat.ID) {
				votesNeeded := database.GetUpvoteCount(chatID)
				text := fmt.Sprintf(`<b>ᴀᴅᴍɪɴ ʀɪɢʜᴛs ɴᴇᴇᴅᴇᴅ</b>
				
ʀᴇғʀᴇsʜ ᴀᴅᴍɪɴ ᴄᴀᴄʜᴇ ᴠɪᴀ : /reload

» %d ᴠᴏᴛᴇs ɴᴇᴇᴅᴇᴅ ғᴏʀ ᴘᴇʀғᴏʀᴍɪɴɢ ᴛʜɪs ᴀᴄᴛɪᴏɴ.`, votesNeeded)

				markup := &tb.ReplyMarkup{}
				markup.Inline(markup.Row(markup.Data("ᴠᴏᴛᴇ", "ADMIN_UpVote", fmt.Sprintf("%d", chatID))))
				b.Send(m.Chat, text, markup, tb.ModeHTML)
				return
			}
			b.Send(m.Chat, "❌ Admin rights needed.")
			return
		}
	}

	// Execute next function if authorized
	next(b, m, chatID)
}

// ActualAdminCB checks if the callback query is executed by an admin
func ActualAdminCB(b *tb.Bot, c *tb.Callback, next func(b *tb.Bot, c *tb.Callback)) {
	if database.IsMaintenance() && !utils.IsSudoer(int64(c.Sender.ID)) {
		b.Respond(c, &tb.CallbackResponse{Text: "Bot is under maintenance.", ShowAlert: true})
		return
	}

	if c.Message.Chat.Type != tb.ChatPrivate && !database.IsNonAdminChat(c.Message.Chat.ID) {
		member, err := b.ChatMemberOf(c.Message.Chat, c.Sender)
		if err != nil || (member.Role != tb.Administrator && member.Role != tb.Creator) {
			if !utils.IsSudoer(int64(c.Sender.ID)) {
				b.Respond(c, &tb.CallbackResponse{Text: "❌ Admin rights needed.", ShowAlert: true})
				return
			}
		}
	}

	next(b, c)
}
