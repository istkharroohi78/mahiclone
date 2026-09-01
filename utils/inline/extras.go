package inline

import (
	tb "gopkg.in/tucnak/telebot.v2"

	// Apne project module ka naam 'ANJALI' lagayein
	"ANJALI/config"
)

// BotPlaylistMarkup creates a row with Support and Close buttons
func BotPlaylistMarkup(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	suppBtn := menu.URL(langData["S_B_9"], config.SupportChat)
	closeBtn := menu.Data(langData["CLOSE_BUTTON"], "close")

	menu.Inline(
		menu.Row(suppBtn, closeBtn),
	)
	return menu
}

// CloseMarkup creates a single close button menu
func CloseMarkup(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	closeBtn := menu.Data(langData["CLOSE_BUTTON"], "close")

	menu.Inline(
		menu.Row(closeBtn),
	)
	return menu
}

// SuppMarkup creates a single support chat button menu
func SuppMarkup(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	suppBtn := menu.URL(langData["S_B_9"], config.SupportChat)

	menu.Inline(
		menu.Row(suppBtn),
	)
	return menu
}
