package admins

import (
	"fmt"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	// Apne project module ka naam 'ANJALI' lagayein
	"ANJALI/utils/database"
)

const photoURL = "https://files.catbox.moe/6r97s4.jpg"

// getPanel generates the caption and inline keyboard for the Autoplay settings
func getPanel(chatID int64, enabled bool) (string, *tb.ReplyMarkup) {
	status := "🔴 𝐃ɪsᴀʙʟᴇᴅ"
	if enabled {
		status = "🟢 𝐄ɴᴀʙʟᴇᴅ"
	}

	caption := fmt.Sprintf(`**🎵 𝐀ᴜᴛᴏ 𝐏ʟᴀʏ 𝐒ᴇᴛᴛɪɴɢ𝐬**

➻ 𝐌ᴀɴᴀɢᴇ 𝐀ᴜᴛᴏ 𝐏ʟᴀʏ ғᴇᴀᴛᴜʀᴇ ғᴏʀ ᴛʜɪs ɢʀᴏᴜᴘ.

**✦ 𝐂ᴜʀʀᴇɴᴛ 𝐒ᴛᴀᴛᴜ𝐬**
%s

━━━━━━━━━━━━━━━
⚡ 𝐏ᴏᴡᴇʀᴇᴅ ʙʏ ➛ the shiv`, status)

	menu := &tb.ReplyMarkup{}

	// Telebot buttons with payloads
	btnEnable := menu.Data("✅ 𝐄ɴᴀʙʟᴇ", "AUTOPLAY_ENABLE", fmt.Sprintf("%d", chatID))
	btnDisable := menu.Data("❌ 𝐃ɪsᴀʙʟᴇ", "AUTOPLAY_DISABLE", fmt.Sprintf("%d", chatID))
	btnStatus := menu.Data(fmt.Sprintf("𝐀ᴜᴛᴏ 𝐏ʟᴀʏ : %s", status), "AUTOPLAY_STATUS")

	menu.Inline(
		menu.Row(btnEnable, btnDisable),
		menu.Row(btnStatus),
	)

	return caption, menu
}

// RegisterAutoplayHandlers binds the /autoplay command and its callbacks
func RegisterAutoplayHandlers(b *tb.Bot) {

	// 1. Command Handler
	b.Handle("/autoplay", func(m *tb.Message) {
		if !m.FromGroup() {
			return
		}
		
		// Note: Ensure isAdmin is accessible (defined in auth.go within the same 'admins' package)
		if !isAdmin(b, m) {
			b.Send(m.Chat, "❌ You must be an admin to change this setting!")
			return
		}

		chatID := m.Chat.ID
		args := strings.Split(m.Text, " ")

		if len(args) > 1 {
			query := strings.ToLower(args[1])
			if query == "enable" || query == "on" {
				database.AutoplayOn(chatID)
				caption, menu := getPanel(chatID, true)
				
				photo := &tb.Photo{File: tb.FromURL(photoURL), Caption: caption + "\n**✅ 𝐀ᴜᴛᴏ 𝐏ʟᴀʏ 𝐄ɴᴀʙʟᴇᴅ 𝐒ᴜᴄᴄᴇssғᴜʟʟʏ!**"}
				b.Send(m.Chat, photo, &tb.SendOptions{ReplyMarkup: menu, ParseMode: tb.ModeMarkdown})
				return
			} else if query == "disable" || query == "off" {
				database.AutoplayOff(chatID)
				caption, menu := getPanel(chatID, false)
				
				photo := &tb.Photo{File: tb.FromURL(photoURL), Caption: caption + "\n**❌ 𝐀ᴜᴛᴏ 𝐏ʟᴀʏ 𝐃ɪsᴀʙʟᴇᴅ 𝐒ᴜᴄᴄᴇssғᴜʟʟʏ!**"}
				b.Send(m.Chat, photo, &tb.SendOptions{ReplyMarkup: menu, ParseMode: tb.ModeMarkdown})
				return
			}
		}

		enabled := database.IsAutoplayOn(chatID)
		caption, menu := getPanel(chatID, enabled)

		photo := &tb.Photo{File: tb.FromURL(photoURL), Caption: caption}
		b.Send(m.Chat, photo, &tb.SendOptions{ReplyMarkup: menu, ParseMode: tb.ModeMarkdown})
	})

	// 2. Enable Callback
	b.Handle("\fAUTOPLAY_ENABLE", func(c *tb.Callback) {
		chatID, _ := strconv.ParseInt(c.Data, 10, 64)

		member, err := b.ChatMemberOf(&tb.Chat{ID: chatID}, c.Sender)
		if err != nil || (member.Role != tb.Administrator && member.Role != tb.Creator) {
			b.Respond(c, &tb.CallbackResponse{Text: "❌ You must be an admin to change this setting!", ShowAlert: true})
			return
		}

		database.AutoplayOn(chatID)
		caption, menu := getPanel(chatID, true)

		photo := &tb.Photo{File: tb.FromURL(photoURL), Caption: caption}
		b.Edit(c.Message, photo, &tb.SendOptions{ReplyMarkup: menu, ParseMode: tb.ModeMarkdown})
		b.Respond(c, &tb.CallbackResponse{Text: "Auto Play Enabled ✅"})
	})

	// 3. Disable Callback
	b.Handle("\fAUTOPLAY_DISABLE", func(c *tb.Callback) {
		chatID, _ := strconv.ParseInt(c.Data, 10, 64)

		member, err := b.ChatMemberOf(&tb.Chat{ID: chatID}, c.Sender)
		if err != nil || (member.Role != tb.Administrator && member.Role != tb.Creator) {
			b.Respond(c, &tb.CallbackResponse{Text: "❌ You must be an admin to change this setting!", ShowAlert: true})
			return
		}

		database.AutoplayOff(chatID)
		caption, menu := getPanel(chatID, false)

		photo := &tb.Photo{File: tb.FromURL(photoURL), Caption: caption}
		b.Edit(c.Message, photo, &tb.SendOptions{ReplyMarkup: menu, ParseMode: tb.ModeMarkdown})
		b.Respond(c, &tb.CallbackResponse{Text: "Auto Play Disabled ❌"})
	})

	// 4. Status Panel Button Callback
	b.Handle("\fAUTOPLAY_STATUS", func(c *tb.Callback) {
		b.Respond(c, &tb.CallbackResponse{Text: "Auto Play Status Panel", ShowAlert: false})
	})
}
