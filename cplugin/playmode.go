package cplugin

import (
	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

func PlaymodeCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	playmode := utils.GetPlaymode(m.Chat.ID)
	isDirect := playmode == "Direct"
	isGroup := !utils.IsNonAdminChat(m.Chat.ID)

	playType := utils.GetPlaytype(m.Chat.ID)
	isPlaytypeEveryone := playType == "Everyone"

	// Fetch dynamic markup from utils/buttons.go
	markup := utils.PlaymodeUsersMarkup(isDirect, isGroup, isPlaytypeEveryone)

	b.Send(m.Chat, "✨ **Playmode Settings**\nSelect the configuration for your group.", markup, tb.ModeMarkdown)
}
