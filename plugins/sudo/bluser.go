package sudo

import (
	"fmt"
	"strconv"
	"strings"

	"ANJALI/utils"
	"ANJALI/utils/database"

	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterBlUserHandlers(b *tb.Bot) {
	b.Handle("/bluser", func(m *tb.Message) {
		if !utils.IsSudoer(int64(m.Sender.ID)) {
			return
		}

		args := strings.Split(m.Text, " ")
		if len(args) != 2 && m.ReplyTo == nil {
			b.Send(m.Chat, "Usage: `/bluser [User ID]` or reply to a user.", tb.ModeMarkdown)
			return
		}

		var targetID int64
		var err error

		if m.ReplyTo != nil {
			targetID = int64(m.ReplyTo.Sender.ID)
		} else {
			targetID, err = strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				b.Send(m.Chat, "❌ **Invalid User ID.**", tb.ModeMarkdown)
				return
			}
		}

		// Prevent banning the bot owner or oneself
		if targetID == int64(m.Sender.ID) || utils.IsSudoer(targetID) {
			b.Send(m.Chat, "❌ **You cannot blacklist yourself or another Sudoer.**", tb.ModeMarkdown)
			return
		}

		database.AddBannedUser(targetID)
		b.Send(m.Chat, fmt.Sprintf("✅ **User `%d` has been blacklisted!**", targetID), tb.ModeMarkdown)
	})

	b.Handle("/unbluser", func(m *tb.Message) {
		if !utils.IsSudoer(int64(m.Sender.ID)) {
			return
		}

		args := strings.Split(m.Text, " ")
		if len(args) != 2 && m.ReplyTo == nil {
			b.Send(m.Chat, "Usage: `/unbluser [User ID]` or reply to a user.", tb.ModeMarkdown)
			return
		}

		var targetID int64
		var err error

		if m.ReplyTo != nil {
			targetID = int64(m.ReplyTo.Sender.ID)
		} else {
			targetID, err = strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				b.Send(m.Chat, "❌ **Invalid User ID.**", tb.ModeMarkdown)
				return
			}
		}

		database.RemoveBannedUser(targetID)
		b.Send(m.Chat, fmt.Sprintf("✅ **User `%d` has been whitelisted!**", targetID), tb.ModeMarkdown)
	})
}
