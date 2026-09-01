package play

import (
	"ANJALI/utils/database"
	"ANJALI/utils/decorators"
	"ANJALI/utils/inline"

	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterPlaymodeHandlers(b *tb.Bot) {
	b.Handle("/playmode", func(m *tb.Message) {
		decorators.Language(b, m, func(b *tb.Bot, m *tb.Message, lang string) {
			loc := map[string]string{ /* Load locale here */ }
			chatID := m.Chat.ID

			isDir := database.GetPlaymode(chatID) == "Direct"
			isGrp := !database.IsNonAdminChat(chatID)
			isPty := database.GetPlaytype(chatID) != "Everyone"

			markup := inline.PlaymodeUsersMarkup(loc, isDir, isGrp, isPty)
			b.Send(m.Chat, "⚙️ **Play Mode Settings**", markup, tb.ModeMarkdown)
		})
	})

	// Alias
	b.Handle("/mode", func(m *tb.Message) { b.Trigger("/playmode", m) })
}
