package cplugin

import (
	"fmt"
	"strings"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

func cleanURLOrUsername(val string) string {
	val = strings.TrimSpace(val)
	val = strings.ReplaceAll(val, "https://", "")
	val = strings.ReplaceAll(val, "http://", "")
	val = strings.ReplaceAll(val, "t.me/", "")
	val = strings.ReplaceAll(val, "telegram.me/", "")
	val = strings.ReplaceAll(val, "@", "")
	return strings.Trim(val, "/")
}

// SetChannelCommand handles /setchannel
func SetChannelCommand(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Only the Bot Owner can use this command.**", tb.ModeMarkdown)
		return
	}

	args := strings.SplitN(m.Text, " ", 2)
	if len(args) != 2 {
		b.Send(m.Chat, "Usage: `/setchannel [username/link]`", tb.ModeMarkdown)
		return
	}

	finalVal := cleanURLOrUsername(args[1])
	utils.SetCloneChannel(b.Me.ID, finalVal)
	b.Send(m.Chat, fmt.Sprintf("✅ Update channel set to: %s", finalVal))
}

// SetSupportCommand handles /setsupport
func SetSupportCommand(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Only the Bot Owner can use this command.**", tb.ModeMarkdown)
		return
	}

	args := strings.SplitN(m.Text, " ", 2)
	if len(args) != 2 {
		b.Send(m.Chat, "Usage: `/setsupport [username/link]`", tb.ModeMarkdown)
		return
	}

	finalVal := cleanURLOrUsername(args[1])
	utils.SetCloneSupport(b.Me.ID, finalVal)
	b.Send(m.Chat, fmt.Sprintf("✅ Support chat set to: %s", finalVal))
}

// BotInfoCommand handles /botinfo
func BotInfoCommand(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Only the Bot Owner can use this command.**", tb.ModeMarkdown)
		return
	}

	cloneData := utils.GetCloneBotData(b.Me.ID)
	channel := "None"
	support := "None"
	if cloneData != nil {
		if cloneData.Channel != "" {
			channel = cloneData.Channel
		}
		if cloneData.SupportChat != "" {
			support = cloneData.SupportChat
		}
	}

	text := fmt.Sprintf("**ʙᴏᴛ ɪɴғᴏ:**\n➤ **ʙᴏᴛ ɪᴅ:** `%d`\n➤ **ᴄʜᴀɴɴᴇʟ:** %s\n➤ **sᴜᴘᴘᴏʀᴛ ᴄʜᴀᴛ:** %s", b.Me.ID, channel, support)
	b.Send(m.Chat, text, tb.ModeMarkdown)
}

// LoggerCommands handles /logger, /setlogger, /logstatus
func LoggerCommands(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Only the Bot Owner can use this command.**", tb.ModeMarkdown)
		return
	}

	args := strings.Split(m.Text, " ")
	cmd := strings.ToLower(args[0])

	switch cmd {
	case "/logstatus":
		status := utils.GetCloneLogStatus(b.Me.ID)
		logChannel := utils.GetCloneLogChannel(b.Me.ID)

		statusStr := "Disabled"
		if status {
			statusStr = "Enabled"
		}

		b.Send(m.Chat, fmt.Sprintf("**ʟᴏɢɢᴇʀ sᴛᴀᴛᴜs :**\n\n - sᴛᴀᴛᴜs : `%s`\n - ʟᴏɢɢᴇʀ ɪᴅ : `%s`", statusStr, logChannel), tb.ModeMarkdown)

	case "/logger":
		if len(args) != 2 || (strings.ToLower(args[1]) != "enable" && strings.ToLower(args[1]) != "disable") {
			b.Send(m.Chat, "**ᴇxᴀᴍᴘʟᴇ :** \n`/logger [ᴇɴᴀʙʟᴇ | ᴅɪsᴀʙʟᴇ]`", tb.ModeMarkdown)
			return
		}

		enable := strings.ToLower(args[1]) == "enable"
		utils.SetCloneLogStatus(b.Me.ID, enable)
		b.Send(m.Chat, fmt.Sprintf("✅ Logger status set to: %v", enable))

	case "/setlogger":
		if len(args) != 2 {
			b.Send(m.Chat, "**ᴇxᴀᴍᴘʟᴇ :** \n- `/setlogger -100xxxxxxxx`", tb.ModeMarkdown)
			return
		}
		utils.SetCloneLogChannel(b.Me.ID, args[1])
		b.Send(m.Chat, fmt.Sprintf("✅ Logger channel set to: `%s`", args[1]), tb.ModeMarkdown)
	}
}
