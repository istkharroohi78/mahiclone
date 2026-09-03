package play

import (
	"ANJALI/utils/database"
	"ANJALI/utils/decorators"
	"ANJALI/utils/inline"

	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterPlaymodeHandlers(b *tb.Bot) {
	// 1. Logic ko ek variable mein daal diya taaki b.Trigger ki zaroorat na pade
	playmodeFunc := func(m *tb.Message) {
		decorators.Language(b, m, func(b *tb.Bot, m *tb.Message, lang string) {
			loc := map[string]string{ /* Load locale here */ }
			chatID := m.Chat.ID

			// 2. Spelling theek kar di (GetPlaymode -> GetPlayMode)
			isDir := database.GetPlayMode(chatID) == "Direct"
			isGrp := !database.IsNonAdminChat(chatID)
			isPty := database.GetPlaytype(chatID) != "Everyone"

			markup := inline.PlaymodeUsersMarkup(loc, isDir, isGrp, isPty)
			b.Send(m.Chat, "⚙️ **Play Mode Settings**", markup, tb.ModeMarkdown)
		})
	}

	// Dono commands ko same logic assign kar diya
	b.Handle("/playmode", playmodeFunc)
	b.Handle("/mode", playmodeFunc)
}
