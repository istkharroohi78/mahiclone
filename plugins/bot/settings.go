package bot

import (
	"ANJALI/utils/database"
	"ANJALI/utils/decorators"
	"ANJALI/utils/inline"

	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterSettingsHandlers(b *tb.Bot) {
	// Command: /settings
	b.Handle("/settings", func(m *tb.Message) {
		decorators.Language(b, m, func(b *tb.Bot, m *tb.Message, lang string) {
			loc := map[string]string{ /* Load locale here */ }
			markup := inline.SettingMarkup(loc)
			b.Send(m.Chat, "⚙️ **Settings Menu**", markup, tb.ModeMarkdown)
		})
	})

	// Callbacks
	b.Handle("\fsettings_helper", func(c *tb.Callback) {
		decorators.LanguageCB(b, c, func(b *tb.Bot, c *tb.Callback, lang string) {
			loc := map[string]string{ /* Load locale here */ }
			markup := inline.SettingMarkup(loc)
			b.Edit(c.Message, "⚙️ **Settings Menu**", markup, tb.ModeMarkdown)
		})
	})

	b.Handle("\fMODECHANGE", func(c *tb.Callback) {
		decorators.ActualAdminCB(b, c, func(b *tb.Bot, c *tb.Callback) {
			chatID := c.Message.Chat.ID
			playMode := database.GetPlaymode(chatID)

			if playMode == "Direct" {
				database.SetPlaymode(chatID, "Inline")
			} else {
				database.SetPlaymode(chatID, "Direct")
			}

			b.Respond(c, &tb.CallbackResponse{Text: "Playmode Changed!"})

			// Refresh UI
			loc := map[string]string{ /* Load locale here */ }
			isDir := database.GetPlaymode(chatID) == "Direct"
			isGrp := !database.IsNonAdminChat(chatID)
			isPty := database.GetPlaytype(chatID) != "Everyone"

			markup := inline.PlaymodeUsersMarkup(loc, isDir, isGrp, isPty)
			b.EditReplyMarkup(c.Message, markup)
		})
	})
}
