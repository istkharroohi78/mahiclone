package cplugin

import (
	"fmt"
	"strings"
	"time"

	"ANJALI/config"
	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// TransferOwner handles /transfer
func TransferOwner(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Access Denied:** Only the Bot Owner can transfer ownership.", tb.ModeMarkdown)
		return
	}

	var newOwner *tb.User
	if m.ReplyTo != nil {
		newOwner = m.ReplyTo.Sender
	} else {
		args := strings.Split(m.Text, " ")
		if len(args) > 1 {
			// Simulating getting user by ID/username
			chat, err := b.ChatByID(args[1])
			if err == nil {
				newOwner = &tb.User{ID: int(chat.ID), FirstName: chat.FirstName, Username: chat.Username}
			}
		}
	}

	if newOwner == nil {
		b.Send(m.Chat, "💡 **Usage:**\nReply to a user or type `/transfer @username`.", tb.ModeMarkdown)
		return
	}
	if newOwner.IsBot {
		b.Send(m.Chat, "❌ You cannot make a bot the owner.")
		return
	}

	utils.SetOwnerIDFromDB(b.Me.ID, int64(newOwner.ID))
	b.Send(m.Chat, fmt.Sprintf("✅ **Ownership Transferred!**\n👑 New Owner: [%s](tg://user?id=%d)", newOwner.FirstName, newOwner.ID), tb.ModeMarkdown)
}

// StartCommand handles /start
func StartCommand(b *tb.Bot, m *tb.Message) {
	botID := b.Me.ID

	if m.Private() {
		utils.AddServedUserClone(int64(m.Sender.ID), botID)

		// Send Loading Animation
		loading, _ := b.Send(m.Chat, "<b>⌛️ ʟᴏᴀᴅɪɴɢ</b>", tb.ModeHTML)
		frames := []string{
			"<b>⌛️ 🇩 🇮 🇳 🇬  🇩 🇴 🇳 🇬 </b>",
			"<b>🌀 🇸 🇹 🇦 🇷 🇹 🇮 🇳 🇬  🇧 🇦 🇧 🇾 💋 </b>",
			fmt.Sprintf("<b>🌀 %s </b>", b.Me.FirstName),
		}
		for _, frame := range frames {
			time.Sleep(300 * time.Millisecond)
			b.Edit(loading, frame, tb.ModeHTML)
		}

		// Fetch configurations
		supportChat := utils.GetCloneSupport(botID)
		if supportChat == "" {
			supportChat = config.LoadConfig().SupportChat
		}

		botNameUpper := strings.ToUpper(b.Me.FirstName)
		caption := fmt.Sprintf(`<b>───[ ˹ %s ˼ 🎵 ]───</b>

<b>Hᴏʟᴏᴏ - !! <span class="tg-spoiler"><a href="tg://user?id=%d">%s</a></span></b>

<b>I ᴀᴍ ᴛʜᴇ ғᴀsᴛ ᴀɴᴅ ᴘᴏᴡᴇʀғᴜʟ ᴍᴜsɪᴄ ᴘʟᴀʏᴇʀ ʙᴏᴛ ᴡɪᴛʜ sᴏᴍᴇ ᴀᴡᴇsᴏᴍᴇ ғᴇᴀᴛᴜʀᴇs.</b>

<blockquote><b>🎶 ʜɪɢʜ-ǫᴜᴧʟɪᴛʏ ᴍᴜꜱɪᴄ ᴘʟᴧʏєʀ ʙσᴛ</b>
<b>ғσʀ ᴛєʟєɢʀᴧϻ ɢʀσᴜᴘꜱ & ᴄʜᴧηηєʟꜱ</b>

<b>🔥 ɪηꜱᴛᴧηᴛ ꜱᴛʀєᴧϻɪηɢ</b>
<b>❤️ ꜱϻσσᴛʜ ᴘʟᴧʏʙᴧᴄᴋ</b>
<b>🎧 ᴄʀʏꜱᴛᴧʟ ꜱσᴜηᴅ | ησ ʟᴧɢ</b></blockquote>

<b>Cʟɪᴄᴋ ᴏɴ ᴛʜᴇ ʜᴇʟᴘ ʙᴜᴛᴛᴏɴ ᴛᴏ ɢᴇᴛ ɪɴғᴏʀᴍᴀᴛɪᴏɴ ᴀʙᴏᴜᴛ ᴍʏ ᴍᴏᴅᴜʟᴇs ᴀɴᴅ ᴄᴏᴍᴍᴀɴᴅs.</b>`, botNameUpper, m.Sender.ID, m.Sender.FirstName)

		markup := utils.MakeStartPanel(b.Me.Username, supportChat, int64(m.Sender.ID))
		photo := &tb.Photo{File: tb.FromURL(GetRandomStartImg()), Caption: caption}

		b.Delete(loading)
		b.Send(m.Chat, photo, markup, tb.ModeHTML)

	} else {
		// Group Start
		utils.AddServedChatClone(m.Chat.ID, botID)
		uptime := utils.GetReadableTime(int(time.Since(bootTime).Seconds()))
		caption := fmt.Sprintf("🔥 **%s is alive!**\n\n**Uptime:** %s", b.Me.FirstName, uptime)

		markup := utils.MakeGPPanel(b.Me.Username, config.LoadConfig().SupportChat)
		photo := &tb.Photo{File: tb.FromURL(GetRandomStartImg()), Caption: caption}
		b.Send(m.Chat, photo, markup, tb.ModeMarkdown)
	}
}

// Media Setters (e.g., /setstartimg, /delstartimg)
func SetStartMediaCommands(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		return
	}

	cmd := strings.ToLower(strings.Split(m.Text, " ")[0])

	switch cmd {
	case "/setstartimg", "/addstartimg":
		if m.ReplyTo != nil && m.ReplyTo.Photo != nil {
			utils.SetCloneSearchType(b.Me.ID, "start_image", m.ReplyTo.Photo.FileID)
			b.Send(m.Chat, "✅ Start Image Added!")
		} else {
			b.Send(m.Chat, "💡 Reply to a photo.")
		}
	case "/delstartimg", "/resetstartimg":
		utils.DeleteCloneSearchType(b.Me.ID)
		b.Send(m.Chat, "✅ Start Images Reset!")
	case "/setstarteffect", "/addstarteffect":
		args := strings.Split(m.Text, " ")
		if len(args) > 1 {
			utils.SetCloneSearchType(b.Me.ID, "start_effect", args[1])
			b.Send(m.Chat, "✅ Effect Added!")
		}
	}
}
