package decorators

import (
	"fmt"

	"ANJALI/config"
	"ANJALI/utils"
	"ANJALI/utils/database"

	tb "gopkg.in/tucnak/telebot.v2"
)

// LanguageStart middleware loads the language preference before executing the command
func LanguageStart(b *tb.Bot, m *tb.Message, next func(b *tb.Bot, m *tb.Message, lang string)) {
	langCode := database.GetLang(m.Chat.ID) // Assumes GetLang defaults to "en"
	next(b, m, langCode)
}

// LanguageCB middleware loads the language preference for callbacks
func LanguageCB(b *tb.Bot, c *tb.Callback, next func(b *tb.Bot, c *tb.Callback, lang string)) {
	if database.IsMaintenance() && !utils.IsSudoer(int64(c.Sender.ID)) {
		b.Respond(c, &tb.CallbackResponse{
			Text:      fmt.Sprintf("%s is under maintenance.", b.Me.FirstName),
			ShowAlert: true,
		})
		return
	}

	langCode := database.GetLang(c.Message.Chat.ID)
	next(b, c, langCode)
}

// Language middleware for standard group commands
func Language(b *tb.Bot, m *tb.Message, next func(b *tb.Bot, m *tb.Message, lang string)) {
	if database.IsMaintenance() && !utils.IsSudoer(int64(m.Sender.ID)) {
		b.Send(m.Chat, fmt.Sprintf("%s ɪs ᴜɴᴅᴇʀ ᴍᴀɪɴᴛᴇɴᴀɴᴄᴇ, ᴠɪsɪᴛ [sᴜᴘᴘᴏʀᴛ ᴄʜᴀᴛ](%s) ғᴏʀ ᴋɴᴏᴡɪɴɢ ᴛʜᴇ ʀᴇᴀsᴏɴ.", b.Me.FirstName, config.LoadConfig().SupportChat), tb.ModeMarkdown)
		return
	}

	b.Delete(m)
	langCode := database.GetLang(m.Chat.ID)
	next(b, m, langCode)
}
