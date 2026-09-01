package admins

import (
	"fmt"
	"strings"

	"ANJALI/utils/database"
	"ANJALI/utils/decorators"
	"ANJALI/utils/inline"

	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterAdminCallbacks(b *tb.Bot) {
	// Support & Clone callbacks
	b.Handle("\fsupport_page", func(c *tb.Callback) {
		b.Respond(c)
		styleMap := inline.GetStyleMap()
		text := `<blockquote><b>✨ ᴡєʟᴄσϻє ᴛσ ᴛʜє sυᴘᴘσʀᴛ ϻєηυ ✨</b>

<b>ɪғ ʏσυ ηєєᴅ ᴧηʏ ʜєʟᴘ ʀєɢᴧʀᴅɪηɢ ᴛʜє ʙσᴛ σʀ ᴡᴧηᴛ ᴛσ ʀєᴘσʀᴛ ᴧ ʙυɢ, ᴊσɪη συʀ sυᴘᴘσʀᴛ ᴄʜᴧᴛ σʀ ᴄʜᴧηηєʟ ʙєʟσᴡ.</b></blockquote>`

		m := &tb.ReplyMarkup{}
		m.Inline(
			m.Row(
				inline.CreateBtn(m, "υᴘᴅᴧᴛєs", "", "https://t.me/betabot_hub", styleMap[1], 6039381989985882045, false),
				inline.CreateBtn(m, "sυᴘᴘσʀᴛ", "", "https://t.me/betabot_support", styleMap[2], 6021618194228187816, false),
			),
			m.Row(inline.CreateBtn(m, "ʙσᴛs", "", "https://t.me/betabot_hub/6701", styleMap[3], 5355051922862653659, false)),
			m.Row(inline.CreateBtn(m, "ʙᴧᴄᴋ", "settingsback_helper", "", styleMap[4], 5352759161945867747, false)),
		)

		b.Edit(c.Message, &tb.Photo{File: tb.FromURL("https://files.catbox.moe/4hl7n8.jpg"), Caption: text}, m, tb.ModeHTML)
	})

	// Master ADMIN Callback Router (Pause/Resume/Stop/Skip/Autoplay)
	b.Handle("\fADMIN", func(c *tb.Callback) {
		decorators.ActualAdminCB(b, c, func(b *tb.Bot, c *tb.Callback) {
			data := strings.Split(c.Data, "|")
			if len(data) < 2 {
				return
			}
			command := strings.TrimSpace(data[0])
			chatIDStr := strings.Split(data[1], "_")[0]
			var chatID int64
			fmt.Sscanf(chatIDStr, "%d", &chatID)

			switch command {
			case "Pause":
				if !database.IsMusicPlaying(chatID) {
					b.Respond(c, &tb.CallbackResponse{Text: "Stream is already paused!", ShowAlert: true})
					return
				}
				database.MusicOff(chatID)
				// Call Core Pause Logic here
				b.Send(c.Message.Chat, fmt.Sprintf("⏸️ **Stream Paused by** [%s](tg://user?id=%d)", c.Sender.FirstName, c.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
				b.Respond(c)

			case "Resume":
				if database.IsMusicPlaying(chatID) {
					b.Respond(c, &tb.CallbackResponse{Text: "Stream is already playing!", ShowAlert: true})
					return
				}
				database.MusicOn(chatID)
				// Call Core Resume Logic here
				b.Send(c.Message.Chat, fmt.Sprintf("▶️ **Stream Resumed by** [%s](tg://user?id=%d)", c.Sender.FirstName, c.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
				b.Respond(c)

			case "Stop", "End":
				database.SetLoop(chatID, 0)
				// Call Core Stop Logic here
				b.Send(c.Message.Chat, fmt.Sprintf("🛑 **Stream Ended by** [%s](tg://user?id=%d)", c.Sender.FirstName, c.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
				b.Delete(c.Message)
				b.Respond(c)

			case "Skip":
				// Pop from Queue & Play Next
				b.Edit(c.Message, fmt.Sprintf("<blockquote><b>▶️ ➻ sᴛʀєᴧϻ sᴋɪᴘᴘєᴅ 🥀</b>\n│ \n└<b>ʙʏ :</b> [%s](tg://user?id=%d)</blockquote>", c.Sender.FirstName, c.Sender.ID), tb.ModeHTML)
				b.Respond(c)
			}
		})
	})
}
