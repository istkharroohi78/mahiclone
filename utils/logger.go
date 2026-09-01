package utils

import (
	"fmt"
	"log"
	"strings"

	"ANJALI/config"

	tb "gopkg.in/tucnak/telebot.v2"
)

var PremiumEmojis = []string{
	"5422831825178206894", "5368324170673489600",
	"5206607081334906820", "5206380668048496464",
}

func GetOwner(b *tb.Bot, chatID int64) string {
	chat, err := b.ChatByID(fmt.Sprintf("%d", chatID))
	if err != nil {
		return "Unknown"
	}
	admins, err := b.AdminsOf(chat)
	if err != nil {
		return "Unknown"
	}
	for _, admin := range admins {
		if admin.Role == tb.Creator {
			return fmt.Sprintf("[%s](tg://user?id=%d)", admin.User.FirstName, admin.User.ID)
		}
	}
	return "Unknown"
}

func PlayLogs(b *tb.Bot, m *tb.Message, streamType string) {
	cfg := config.LoadConfig()
	if cfg.LoggerID == 0 || m.Chat.ID == cfg.LoggerID {
		return
	}

	query := "Link/File or Reply"
	args := strings.SplitN(m.Text, " ", 2)
	if len(args) > 1 {
		query = args[1]
	}

	memberCount, _ := b.Len(m.Chat)
	owner := GetOwner(b, m.Chat.ID)

	chatLink := ""
	if m.Chat.Username != "" {
		chatLink = "https://t.me/" + m.Chat.Username
	}

	text := fmt.Sprintf(`<blockquote><b>%s ᴘʟᴀʏ ʟᴏɢ</b>

<b>• ʀᴇǫᴜᴇsᴛ ʙʏ : [%s](tg://user?id=%d)</b>
<b>• ǫᴜᴇʀʏ : %s</b>
<b>• ᴄʜᴀᴛ : %s [<code>%d</code>]</b>
<b>• ᴏᴡɴᴇʀ : %s</b>
<b>• ᴍᴇᴍʙᴇʀs : %d</b></blockquote>`,
		b.Me.FirstName, m.Sender.FirstName, m.Sender.ID, query, m.Chat.Title, m.Chat.ID, owner, memberCount)

	markup := &tb.ReplyMarkup{}
	var buttons []tb.Btn
	if chatLink != "" {
		buttons = append(buttons, markup.URL("ɢʀᴏᴜᴘ ʟɪɴᴋ", chatLink))
	}
	buttons = append(buttons, markup.URL("sᴜᴘᴘᴏʀᴛ", "https://t.me/betabot_support"))
	markup.Inline(markup.Row(buttons...))

	logChat, err := b.ChatByID(fmt.Sprintf("%d", cfg.LoggerID))
	if err == nil {
		b.Send(logChat, text, tb.ModeHTML, markup, &tb.SendOptions{DisableWebPagePreview: true})
	}
}

func CloneBotLogs(b *tb.Bot, m *tb.Message, cloneLoggerID int64, streamType string) {
	query := "Link/File or Reply"
	args := strings.SplitN(m.Text, " ", 2)
	if len(args) > 1 {
		query = args[1]
	}

	ownerLogText := fmt.Sprintf(`<b><a href="https://t.me/%s">%s</a> ᴘʟᴀʏ ʟᴏɢ</b>

<b>• ʀᴇǫᴜᴇsᴛ ʙʏ :</b> <a href="tg://user?id=%d">%s</a>
<b>• ǫᴜᴇʀʏ :</b> %s
<b>• ᴄʜᴀᴛ :</b> %s [<code>%d</code>]`,
		b.Me.Username, b.Me.FirstName, m.Sender.ID, m.Sender.FirstName, query, m.Chat.Title, m.Chat.ID)

	markup := &tb.ReplyMarkup{}
	markup.Inline(markup.Row(markup.URL("sᴜᴘᴘᴏʀᴛ", "https://t.me/betabot_support")))

	if cloneLoggerID != 0 && m.Chat.ID != cloneLoggerID {
		logChat, err := b.ChatByID(fmt.Sprintf("%d", cloneLoggerID))
		if err == nil {
			b.Send(logChat, ownerLogText, tb.ModeHTML, markup, &tb.SendOptions{DisableWebPagePreview: true})
		} else {
			log.Printf("[ERROR] Sending to Clone Owner Log Failed: %v", err)
		}
	}
}
