package cplugin

import (
	"fmt"
	"strconv"
	"strings"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// SeekCommand handles /seek and /seekback
func SeekCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}
	chatID := m.Chat.ID

	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Sirf Admins he is command ko use kar sakte hain!**", tb.ModeMarkdown)
		return
	}

	args := strings.Split(m.Text, " ")
	if len(args) == 1 {
		b.Send(m.Chat, "⚠️ Please provide duration to seek in seconds. Example: `/seek 10`", tb.ModeMarkdown)
		return
	}

	durationToSkip, err := strconv.Atoi(args[1])
	if err != nil {
		b.Send(m.Chat, "⚠️ Please provide a valid number.", tb.ModeMarkdown)
		return
	}

	playingData := utils.GetPlayingData(chatID)
	if playingData == nil {
		b.Send(m.Chat, "⚠️ Nothing is playing right now.", tb.ModeMarkdown)
		return
	}

	durationSeconds := playingData.Seconds
	durationPlayed := playingData.Played

	if durationSeconds == 0 {
		b.Send(m.Chat, "⚠️ Cannot seek live streams.", tb.ModeMarkdown)
		return
	}

	isBackward := strings.Contains(strings.ToLower(args[0]), "back")
	var toSeek int

	if isBackward {
		toSeek = durationPlayed - durationToSkip
		if toSeek < 0 {
			toSeek = 0
		}
	} else {
		toSeek = durationPlayed + durationToSkip
		if toSeek > durationSeconds {
			toSeek = durationSeconds - 5
		}
	}

	mystic, _ := b.Send(m.Chat, "🔄 Seeking stream...", tb.ModeMarkdown)

	err = utils.SeekStream(chatID, toSeek)
	if err != nil {
		b.Edit(mystic, "❌ Error seeking stream.")
		return
	}

	utils.UpdatePlayedTime(chatID, toSeek)

	b.Edit(mystic, fmt.Sprintf("✅ **Seeked stream to %s**\n\n**Admin:** [%s](tg://user?id=%d)", utils.SecondsToMin(toSeek), m.Sender.FirstName, m.Sender.ID), tb.ModeMarkdown)
}
