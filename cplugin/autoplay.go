package cplugin

import (
	"fmt"
	"strconv"
	"strings"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

const AutoplayBanner = "https://files.catbox.moe/6r97s4.jpg"

// PremiumEmojis list (stored for reference, telebot v2 uses standard UI buttons)
var PremiumEmojis = []string{
	"5258362837411045098",
	"6102938383456146362",
	"5463274047771000031",
	"6100397162976252509",
	"5373310679241466020",
	"5408916593780470262",
	"5776182936638329359",
	"5258389041006518073",
	"6280269890821558384",
	"5936143551854285132",
	"6172332822892647766",
	"5891211339170326418",
	"5409368076447657845",
	"6172312314423808834",
	"6082387600599944892",
	"6271537028307881531",
}

func AutoplayCaption(enabled bool) string {
	status := "🔴 𝐃ɪsᴀʙʟᴇᴅ"
	if enabled {
		status = "🟢 𝐄ɴᴀʙʟᴇᴅ"
	}

	return fmt.Sprintf(`**🎵 𝐀ᴜᴛᴏ 𝐏ʟᴀʏ 𝐒ᴇᴛᴛɪɴɢ𝐬**

➻ 𝐌ᴀɴᴀɢᴇ 𝐀ᴜᴛᴏ 𝐏ʟᴀʏ ғᴇᴀᴛᴜʀᴇ ғᴏʀ ᴛʜɪs ɢʀᴏᴜᴘ.

**✦ 𝐂ᴜʀʀᴇɴᴛ 𝐒ᴛᴀᴛᴜ𝐬**
%s

➻ 𝐖ʜᴇɴ 𝐀ᴜᴛᴏ 𝐏ʟᴀʏ ɪ𝐬 𝐄ɴᴀʙʟᴇᴅ, ᴛʜᴇ ʙᴏᴛ ᴡɪʟʟ
ᴀᴜᴛᴏᴍᴀᴛɪᴄᴀʟʟʏ ᴘʟᴀʏ ʀᴇᴄᴏᴍᴍᴇɴᴅᴇᴅ ᴛʀᴀᴄᴋ𝐬
ᴡʜᴇɴ ᴛʜᴇ ǫᴜᴇᴜᴇ ᴇɴᴅ𝐬.

━━━━━━━━━━━━━━━
⚡ 𝐏ᴏᴡᴇʀᴇᴅ ʙʏ ➛ 𝐁𝐞ᴛᴀ𝐁ᴏᴛ𝐬`, status)
}

func AutoplayPanelMarkup(chatID int64, enabled bool) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}

	status := "🔴 𝐃ɪsᴀʙʟᴇᴅ"
	if enabled {
		status = "🟢 𝐄ɴᴀʙʟᴇᴅ"
	}

	btnEnable := m.Data("𝐀ᴜᴛᴏ 𝐏ʟᴀʏ 𝐄ɴᴀʙʟᴇ", fmt.Sprintf("AP_EN|%d", chatID))
	btnDisable := m.Data("𝐀ᴜᴛᴏ 𝐏ʟᴀʏ 𝐃ɪsᴀʙʟᴇ", fmt.Sprintf("AP_DIS|%d", chatID))
	btnStatus := m.Data(fmt.Sprintf("𝐀ᴜᴛᴏ 𝐏ʟᴀʏ : %s", status), "AP_STAT")

	m.Inline(
		m.Row(btnEnable, btnDisable),
		m.Row(btnStatus),
	)
	return m
}

// AutoplayCommand handles the /autoplay command
func AutoplayCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Sirf Admins he is command ko use kar sakte hain!**", tb.ModeMarkdown)
		return
	}

	payload := strings.ToLower(strings.TrimSpace(m.Payload))

	if payload == "enable" || payload == "on" {
		utils.AutoplayOn(m.Chat.ID)
		text := fmt.Sprintf("<blockquote><b>🟢 🎧 Ʌυᴛσᴘʟᴧʏ sʏsᴛєϻ</b>\n\n<b>Ʌυᴛσᴘʟᴧʏ ғσʀ ᴛʜɪs ɢʀσυᴘ ɪs ησᴡ єηᴧʙʟєᴅ 🟢.</b>\n└ <b>ʙʏ :</b> [%s](tg://user?id=%d)</blockquote>", m.Sender.FirstName, m.Sender.ID)
		b.Send(m.Chat, text, tb.ModeHTML)
		return
	} else if payload == "disable" || payload == "off" {
		utils.AutoplayOff(m.Chat.ID)
		text := fmt.Sprintf("<blockquote><b>🔴 🎧 Ʌυᴛσᴘʟᴧʏ sʏsᴛєϻ</b>\n\n<b>Ʌυᴛσᴘʟᴧʏ ғσʀ ᴛʜɪs ɢʀσυᴘ ɪs ησᴡ ᴅɪsᴧʙʟєᴅ 🔴.</b>\n└ <b>ʙʏ :</b> [%s](tg://user?id=%d)</blockquote>", m.Sender.FirstName, m.Sender.ID)
		b.Send(m.Chat, text, tb.ModeHTML)
		return
	}

	enabled := utils.IsAutoplayOn(m.Chat.ID)
	photo := &tb.Photo{File: tb.FromURL(AutoplayBanner), Caption: AutoplayCaption(enabled)}

	b.Send(m.Chat, photo, AutoplayPanelMarkup(m.Chat.ID, enabled), tb.ModeMarkdown)
}

// AutoplayCallback handles inline button toggles for Autoplay
func AutoplayCallback(b *tb.Bot, c *tb.Callback) {
	if !isAdminCheck(b, c.Message.Chat, int64(c.Sender.ID)) {
		b.Respond(c, &tb.CallbackResponse{
			Text:      "❌ You must be an admin to change this setting!",
			ShowAlert: true,
		})
		return
	}

	data := strings.Split(c.Data, "|")
	if len(data) < 2 {
		return
	}

	action := data[0]
	chatID, _ := strconv.ParseInt(data[1], 10, 64)

	var enabled bool
	if strings.Contains(action, "AP_EN") {
		utils.AutoplayOn(chatID)
		enabled = true
		b.Respond(c, &tb.CallbackResponse{Text: "🟢 𝐀ᴜᴛᴏ 𝐏ʟᴀʏ 𝐄ɴᴀʙʟᴇᴅ"})
	} else if strings.Contains(action, "AP_DIS") {
		utils.AutoplayOff(chatID)
		enabled = false
		b.Respond(c, &tb.CallbackResponse{Text: "🔴 𝐀ᴜᴛᴏ 𝐏ʟᴀʏ 𝐃ɪsᴀʙʟᴇᴅ"})
	}

	photo := &tb.Photo{File: tb.FromURL(AutoplayBanner), Caption: AutoplayCaption(enabled)}
	b.Edit(c.Message, photo, AutoplayPanelMarkup(chatID, enabled), tb.ModeMarkdown)
}

// AutoplayStatusCallback handles the status button click
func AutoplayStatusCallback(b *tb.Bot, c *tb.Callback) {
	b.Respond(c, &tb.CallbackResponse{Text: "⚡ 𝐀ᴜᴛᴏ 𝐏ʟᴀʏ 𝐒ᴛᴀᴛᴜ𝐬"})
}
