package inline

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func SpeedMarkup(loc map[string]string, chatID int64) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	m.Inline(
		m.Row(
			m.Data("🕒 0.5x", fmt.Sprintf("SpeedUP %d|0.5", chatID)),
			m.Data("🕓 0.75x", fmt.Sprintf("SpeedUP %d|0.75", chatID)),
		),
		m.Row(m.Data(loc["P_B_4"], fmt.Sprintf("SpeedUP %d|1.0", chatID))),
		m.Row(
			m.Data("🕤 1.5x", fmt.Sprintf("SpeedUP %d|1.5", chatID)),
			m.Data("🕛 2.0x", fmt.Sprintf("SpeedUP %d|2.0", chatID)),
		),
		m.Row(m.Data(loc["CLOSE_BUTTON"], "close")),
	)
	return m
}
