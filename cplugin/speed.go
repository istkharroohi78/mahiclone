package cplugin

import (
	"fmt"
	"strings"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// SpeedCommand handles /speed, /playback
func SpeedCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}
	chatID := m.Chat.ID

	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Admin rights required!**", tb.ModeMarkdown)
		return
	}

	playingData := utils.GetPlayingData(chatID)
	if playingData == nil || playingData.Seconds == 0 {
		b.Send(m.Chat, "⚠️ Nothing is playing or you can't change speed for this stream.", tb.ModeMarkdown)
		return
	}

	markup := utils.SpeedMarkup(chatID)
	b.Send(m.Chat, fmt.Sprintf("⚡ **Playback Speed Control**\nSelect speed for %s:", b.Me.FirstName), markup, tb.ModeMarkdown)
}

// SpeedCallback handles SpeedUP callback
func SpeedCallback(b *tb.Bot, c *tb.Callback) {
	if !isAdminCheck(b, c.Message.Chat, int64(c.Sender.ID)) {
		b.Respond(c, &tb.CallbackResponse{Text: "❌ Admin rights needed!", ShowAlert: true})
		return
	}

	data := strings.Split(c.Data, " ")
	if len(data) < 2 {
		return
	}
	args := strings.Split(data[1], "|")
	if len(args) < 2 {
		return
	}
	chatID := c.Message.Chat.ID
	speed := args[1]

	playingData := utils.GetPlayingData(chatID)
	if playingData == nil {
		b.Respond(c, &tb.CallbackResponse{Text: "⚠️ Queue is empty!", ShowAlert: true})
		return
	}

	b.Respond(c, &tb.CallbackResponse{Text: "Changing speed..."})
	mystic, _ := b.Send(c.Message.Chat, "🔄 Updating playback speed...")

	err := utils.SpeedUpStream(chatID, speed)
	if err != nil {
		b.Edit(mystic, "❌ Failed to change speed.")
		return
	}

	b.Edit(mystic, fmt.Sprintf("⚡ **Speed set to %sx by [%s](tg://user?id=%d)**", speed, c.Sender.FirstName, c.Sender.ID), utils.CloseMarkup(), tb.ModeMarkdown)
}
