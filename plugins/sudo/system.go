package sudo

import (
	"strings"

	"ANJALI/utils"
	"ANJALI/utils/database"

	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterSystemSettingsHandlers(b *tb.Bot) {
	// Logger Toggle
	b.Handle("/logger", func(m *tb.Message) {
		if !utils.IsSudoer(int64(m.Sender.ID)) {
			return
		}
		args := strings.Split(m.Text, " ")
		if len(args) != 2 {
			b.Send(m.Chat, "Usage: `/logger [enable|disable]`")
			return
		}

		state := strings.ToLower(args[1])
		if state == "enable" {
			database.AddOn(2) // 2 represents logger state in your DB config
			b.Send(m.Chat, "✅ **Logger Enabled.**")
		} else if state == "disable" {
			database.AddOff(2)
			b.Send(m.Chat, "❌ **Logger Disabled.**")
		}
	})

	// Maintenance Toggle
	b.Handle("/maintenance", func(m *tb.Message) {
		if !utils.IsSudoer(int64(m.Sender.ID)) {
			return
		}
		args := strings.Split(m.Text, " ")
		if len(args) != 2 {
			b.Send(m.Chat, "Usage: `/maintenance [enable|disable]`")
			return
		}

		state := strings.ToLower(args[1])
		if state == "enable" {
			database.MaintenanceOn()
			b.Send(m.Chat, "✅ **Maintenance Mode Enabled.** Bot will not process user queries.")
		} else if state == "disable" {
			database.MaintenanceOff()
			b.Send(m.Chat, "❌ **Maintenance Mode Disabled.** Bot is now active.")
		}
	})
}
