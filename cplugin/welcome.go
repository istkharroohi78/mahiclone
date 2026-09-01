package cplugin

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	welcomeState   = make(map[int64]bool)
	lastWelcomeMsg = make(map[int64]*tb.Message)
	customWelcomes = make(map[int64]map[string]interface{})
	weltimeState   = make(map[int64]int)
)

// ParseButtons extracts inline buttons and cleans text
func ParseButtons(text string) (string, *tb.ReplyMarkup) {
	if text == "" {
		return "", nil
	}

	var cleanText []string
	markup := &tb.ReplyMarkup{}
	var allRows []tb.Row

	colorMap := map[string]int{
		"red":   Danger,
		"green": Success,
		"blue":  Primary,
	}
	availableStyles := []int{Primary, Success, Danger}

	lines := strings.Split(text, "\n")
	btnRegex := regexp.MustCompile(`\[(.+?)\]\(buttonurl:([^\s\)]+)(?:\s+color:(red|green|blue))?\)`)
	fallbackRegex := regexp.MustCompile(`\[(.+?)\]\(buttonurl:(.+?)\)`)

	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "(buttonurl:") {
			var row []tb.Btn
			parts := strings.Split(line, "|")

			for _, part := range parts {
				matches := btnRegex.FindStringSubmatch(part)
				if len(matches) > 0 {
					btnName := strings.TrimSpace(matches[1])
					btnURL := strings.TrimSpace(matches[2])
					colorStr := ""
					if len(matches) > 3 {
						colorStr = strings.ToLower(matches[3])
					}

					style := availableStyles[rand.Intn(len(availableStyles))]
					if val, ok := colorMap[colorStr]; ok {
						style = val
					}
					row = append(row, CreateBtn(markup, btnName, "", btnURL, style, false))
				} else {
					fallbackMatches := fallbackRegex.FindStringSubmatch(part)
					if len(fallbackMatches) > 0 {
						btnName := strings.TrimSpace(fallbackMatches[1])
						btnURL := strings.TrimSpace(fallbackMatches[2])
						style := availableStyles[rand.Intn(len(availableStyles))]
						row = append(row, CreateBtn(markup, btnName, "", btnURL, style, false))
					}
				}
			}
			if len(row) > 0 {
				allRows = append(allRows, markup.Row(row...))
			}
		} else {
			cleanText = append(cleanText, line)
		}
	}

	if len(allRows) > 0 {
		markup.Inline(allRows...)
		return strings.Join(cleanText, "\n"), markup
	}
	return strings.Join(cleanText, "\n"), nil
}

// AutoDeleteMessage deletes a message after a delay
func AutoDeleteMessage(b *tb.Bot, m *tb.Message, delaySeconds int) {
	if delaySeconds <= 0 {
		return
	}
	go func() {
		time.Sleep(time.Duration(delaySeconds) * time.Second)
		b.Delete(m)
	}()
}

// WelcomeToggleCommand handles /welcome
func WelcomeToggleCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}
	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "**sᴏʀʀʏ ᴏɴʟʏ ᴀᴅᴍɪɴs ᴄᴀɴ ᴇɴᴀʙʟᴇ ᴡᴇʟᴄᴏᴍᴇ ɴᴏᴛɪғɪᴄᴀᴛɪᴏɴ!**", tb.ModeMarkdown)
		return
	}

	args := strings.Split(m.Text, " ")
	if len(args) != 2 || (strings.ToLower(args[1]) != "on" && strings.ToLower(args[1]) != "off") {
		b.Send(m.Chat, "**ᴜsᴀɢᴇ:**\n**⦿ /welcome [on|off]**", tb.ModeMarkdown)
		return
	}

	state := strings.ToLower(args[1])
	if state == "on" {
		welcomeState[m.Chat.ID] = true
		b.Send(m.Chat, fmt.Sprintf("**ᴇɴᴀʙʟᴇᴅ ᴡᴇʟᴄᴏᴍᴇ ɴᴏᴛɪғɪᴄᴀᴛɪᴏɴ ɪɴ %s**", m.Chat.Title), tb.ModeMarkdown)
	} else {
		welcomeState[m.Chat.ID] = false
		b.Send(m.Chat, fmt.Sprintf("**ᴅɪsᴀʙʟᴇᴅ ᴡᴇʟᴄᴏᴍᴇ ɴᴏᴛɪғɪᴄᴀᴛɪᴏɴ ɪɴ %s**", m.Chat.Title), tb.ModeMarkdown)
	}
}

// SetCustomWelcomeCommand handles /set_welcome
func SetCustomWelcomeCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}
	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "**sᴏʀʀʏ ᴏɴʟʏ ᴀᴅᴍɪɴs ᴄᴀɴ ᴜsᴇ ᴛʜɪs ᴄᴏᴍᴍᴀɴᴅ!**", tb.ModeMarkdown)
		return
	}

	args := strings.SplitN(m.Text, " ", 2)
	cmdText := ""
	if len(args) > 1 {
		cmdText = args[1]
	}

	if m.ReplyTo == nil && cmdText == "" {
		exampleText := `**⚠️ ᴡᴇʟᴄᴏᴍᴇ sᴇᴛ ᴋᴀʀɴᴇ ᴋᴇ ʟɪʏᴇ ᴋɪsɪ ᴍᴇssᴀɢᴇ (Pʜᴏᴛᴏ/Vɪᴅᴇᴏ/Gɪғ) ᴘᴀʀ ʀᴇᴘʟʏ ᴋᴀʀᴇɪɴ ʏᴀ ᴄᴏᴍᴍᴀɴᴅ ᴋᴇ sᴀᴀᴛʜ ᴛᴇxᴛ ʟɪᴋʜᴇɪɴ!**

**👇 ᴄᴏᴘʏ-ᴘᴀsᴛᴇ ᴇxᴀᴍᴘʟᴇ:**
/set_welcome ❖ Hᴇʟʟᴏ {mention}!
❖ Wᴇʟᴄᴏᴍᴇ ᴛᴏ ᴏᴜʀ ɢʀᴏᴜᴘ.
❖ Yᴏᴜ ᴀʀᴇ ᴏᴜʀ {count}ᴛʜ ᴍᴇᴍʙᴇʀ.

[❤ Dᴇᴠᴇʟᴏᴘᴇʀ](buttonurl:https://t.me/THE_SHIV color:red) | [✅ Uᴘᴅᴀᴛᴇs](buttonurl:https://t.me/Channel color:green)
[🛠 Sᴜᴘᴘᴏʀᴛ](buttonurl:https://t.me/Support color:blue)

**💡 Tɪᴘ:** Uᴘᴀʀ ᴡᴀʟᴇ ᴄᴏᴅᴇ ᴋᴏ ᴄᴏᴘʏ ᴋᴀʀᴋᴇ ʙʜᴇᴊ ᴅᴇɪɴ, ᴀᴀᴘᴋᴀ ᴡᴇʟᴄᴏᴍᴇ sᴇᴛ ʜᴏ ᴊᴀʏᴇɢᴀ!
⏱️ **Aᴜᴛᴏ-Dᴇʟᴇᴛᴇ:** Wᴇʟᴄᴏᴍᴇ ᴋᴏ ᴀᴜᴛᴏ ᴅᴇʟᴇᴛᴇ ᴋᴀʀɴᴇ ᴋᴇ ʟɪʏᴇ /weltime 5 sᴇᴛ ᴋᴀʀᴇɪɴ.`
		b.Send(m.Chat, exampleText, tb.ModeMarkdown)
		return
	}

	msgType := "text"
	var fileID string
	rawText := cmdText

	if m.ReplyTo != nil {
		if rawText == "" {
			rawText = m.ReplyTo.Text
			if rawText == "" {
				rawText = m.ReplyTo.Caption
			}
		}

		if m.ReplyTo.Photo != nil {
			msgType, fileID = "photo", m.ReplyTo.Photo.FileID
		} else if m.ReplyTo.Video != nil {
			msgType, fileID = "video", m.ReplyTo.Video.FileID
		} else if m.ReplyTo.Animation != nil {
			msgType, fileID = "animation", m.ReplyTo.Animation.FileID
		} else if m.ReplyTo.Sticker != nil {
			msgType, fileID = "sticker", m.ReplyTo.Sticker.FileID
		}
	}

	cleanText, customMarkup := ParseButtons(rawText)

	customWelcomes[m.Chat.ID] = map[string]interface{}{
		"type":    msgType,
		"file_id": fileID,
		"text":    cleanText,
		"markup":  customMarkup,
	}

	b.Send(m.Chat, "**✅ ᴄᴜsᴛᴏᴍ ᴡᴇʟᴄᴏᴍᴇ ᴡɪᴛʜ ʙᴜᴛᴛᴏɴs sᴇᴛ sᴜᴄᴄᴇssғᴜʟʟʏ!**", tb.ModeMarkdown)
}

// ClearWelcomeCommand handles /cwelcome
func ClearWelcomeCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}
	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "**sᴏʀʀʏ ᴏɴʟʏ ᴀᴅᴍɪɴs ᴄᴀɴ ᴜsᴇ ᴛʜɪs ᴄᴏᴍᴍᴀɴᴅ!**", tb.ModeMarkdown)
		return
	}

	if _, exists := customWelcomes[m.Chat.ID]; exists {
		delete(customWelcomes, m.Chat.ID)
		b.Send(m.Chat, "**✅ ᴄᴜsᴛᴏᴍ ᴡᴇʟᴄᴏᴍᴇ ᴄʟᴇᴀʀᴇᴅ!\n\nɴᴏᴡ ᴅᴇғᴀᴜʟᴛ ɪᴍᴀɢᴇ ᴡᴇʟᴄᴏᴍᴇ ᴡɪʟʟ ʙᴇ ᴜsᴇᴅ.**", tb.ModeMarkdown)
	} else {
		b.Send(m.Chat, "**⚠️ ɴᴏ ᴄᴜsᴛᴏᴍ ᴡᴇʟᴄᴏᴍᴇ ɪs sᴇᴛ ғᴏʀ ᴛʜɪs ɢʀᴏᴜᴘ.**", tb.ModeMarkdown)
	}
}

// WeltimeCommand handles /weltime
func WeltimeCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}
	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "**sᴏʀʀʏ ᴏɴʟʏ ᴀᴅᴍɪɴs ᴄᴀɴ ᴜsᴇ ᴛʜɪs ᴄᴏᴍᴍᴀɴᴅ!**", tb.ModeMarkdown)
		return
	}

	args := strings.Split(m.Text, " ")
	if len(args) != 2 {
		b.Send(m.Chat, "**ᴜsᴀɢᴇ:**\n**⦿ /weltime [minutes|off]** (e.g., /weltime 5)", tb.ModeMarkdown)
		return
	}

	val := strings.ToLower(args[1])
	if val == "off" {
		weltimeState[m.Chat.ID] = 0
		b.Send(m.Chat, "**✅ ᴡᴇʟᴄᴏᴍᴇ ᴀᴜᴛᴏ-ᴅᴇʟᴇᴛᴇ ᴅɪsᴀʙʟᴇᴅ.**", tb.ModeMarkdown)
		return
	}

	var minutes int
	if _, err := fmt.Sscanf(val, "%d", &minutes); err == nil {
		weltimeState[m.Chat.ID] = minutes * 60
		b.Send(m.Chat, fmt.Sprintf("**✅ ᴡᴇʟᴄᴏᴍᴇ ᴍᴇssᴀɢᴇs ᴡɪʟʟ ɴᴏᴡ ʙᴇ ᴅᴇʟᴇᴛᴇᴅ ᴀғᴛᴇʀ %d ᴍɪɴᴜᴛᴇs.**", minutes), tb.ModeMarkdown)
	}
}

// ChatMemberHandler processes ChatMemberUpdated events
func ChatMemberHandler(b *tb.Bot, update *tb.ChatMemberUpdate) {
	if update.NewChatMember == nil || update.OldChatMember != nil {
		return
	}
	if update.NewChatMember.Role == tb.Kicked || update.NewChatMember.Role == tb.Left {
		return
	}

	chatID := update.Chat.ID
	if state, exists := welcomeState[chatID]; exists && !state {
		return
	}

	user := update.NewChatMember.User
	memberCount, _ := b.Len(update.Chat)

	if oldMsg, exists := lastWelcomeMsg[chatID]; exists {
		b.Delete(oldMsg)
	}

	mention := fmt.Sprintf("[%s](tg://user?id=%d)", user.FirstName, user.ID)
	var sentMsg *tb.Message

	if custom, ok := customWelcomes[chatID]; ok {
		text := custom["text"].(string)
		text = strings.ReplaceAll(text, "{mention}", mention)
		text = strings.ReplaceAll(text, "{id}", fmt.Sprintf("%d", user.ID))
		text = strings.ReplaceAll(text, "{username}", "@"+user.Username)
		text = strings.ReplaceAll(text, "{count}", fmt.Sprintf("%d", memberCount))

		markup, _ := custom["markup"].(*tb.ReplyMarkup)

		switch custom["type"].(string) {
		case "text":
			sentMsg, _ = b.Send(update.Chat, text, markup, tb.ModeMarkdown)
		case "photo":
			sentMsg, _ = b.Send(update.Chat, &tb.Photo{File: tb.File{FileID: custom["file_id"].(string)}, Caption: text}, markup, tb.ModeMarkdown)
		case "video":
			sentMsg, _ = b.Send(update.Chat, &tb.Video{File: tb.File{FileID: custom["file_id"].(string)}, Caption: text}, markup, tb.ModeMarkdown)
		case "animation":
			sentMsg, _ = b.Send(update.Chat, &tb.Animation{File: tb.File{FileID: custom["file_id"].(string)}, Caption: text}, markup, tb.ModeMarkdown)
		case "sticker":
			sentMsg, _ = b.Send(update.Chat, &tb.Sticker{File: tb.File{FileID: custom["file_id"].(string)}}, markup)
		}
	} else {
		// Default Welcome
		caption := fmt.Sprintf(`**⎊─────☵ ᴡᴇʟᴄᴏᴍᴇ ☵─────⎊**

**▬▭▬▭▬▭▬▭▬▭▬▭▬▭▬**

**☉ ɴᴀᴍᴇ ⧽** %s
**☉ ɪᴅ ⧽** %d
**☉ ᴜ_ɴᴀᴍᴇ ⧽** @%s
**☉ ᴛᴏᴛᴀʟ ᴍᴇᴍʙᴇʀs ⧽** %d

**▬▭▬▭▬▭▬▭▬▭▬▭▬▭▬**

**⎉──────▢✭ 侖 ✭▢──────⎉**`, mention, user.ID, user.Username, memberCount)

		markup := &tb.ReplyMarkup{}
		btnView := CreateBtn(markup, "๏ ᴠɪᴇᴡ ɴᴇᴡ ᴍᴇᴍʙᴇʀ ๏", "", fmt.Sprintf("tg://openmessage?user_id=%d", user.ID), Primary, false)
		btnKidnap := CreateBtn(markup, "✙ ᴋɪᴅɴᴀᴘ ᴍᴇ ✙", "", fmt.Sprintf("https://t.me/%s?startgroup=true", b.Me.Username), Success, false)

		markup.Inline(
			markup.Row(btnView),
			markup.Row(btnKidnap),
		)

		sentMsg, _ = b.Send(update.Chat, caption, markup, tb.ModeMarkdown)
	}

	if sentMsg != nil {
		lastWelcomeMsg[chatID] = sentMsg
		delay := 300
		if val, exists := weltimeState[chatID]; exists {
			delay = val
		}
		AutoDeleteMessage(b, sentMsg, delay)
	}
}
