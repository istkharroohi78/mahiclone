package cplugin

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"ANJALI/config"
	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// SettingsBackHelper handles the settings back button
func SettingsBackHelper(b *tb.Bot, c *tb.Callback) {
	b.Respond(c, &tb.CallbackResponse{})

	// Assuming GetRandomImg and PrivatePanel are defined in utils
	img := utils.GetRandomImg(config.StartImgURL)
	photo := &tb.Photo{File: tb.FromURL(img), Caption: fmt.Sprintf("Hey %s,\nWelcome back to Settings!", c.Sender.FirstName)}

	// Edit media logic in telebot (simplified for text/caption update if media update is complex)
	b.Edit(c.Message, photo, utils.PrivatePanel(), tb.ModeMarkdown)
}

// AdminCallback handles all ADMIN prefixed inline button actions
func AdminCallback(b *tb.Bot, c *tb.Callback) {
	data := strings.SplitN(strings.TrimSpace(c.Data), " ", 2)
	if len(data) < 2 {
		return
	}

	request := strings.Split(data[1], "|")
	if len(request) < 2 {
		return
	}

	command := request[0]
	chatIDStr := request[1]

	// Handle Upvote chat ID format (e.g., chatID_counter)
	if strings.Contains(chatIDStr, "_") {
		chatIDStr = strings.Split(chatIDStr, "_")[0]
	}

	chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

	if !utils.IsActiveChat(chatID) {
		b.Respond(c, &tb.CallbackResponse{Text: "⚠️ No active stream in this chat.", ShowAlert: true})
		return
	}

	// Admin Check
	if !isAdminCheck(b, c.Message.Chat, int64(c.Sender.ID)) {
		b.Respond(c, &tb.CallbackResponse{Text: "❌ Admin rights needed!", ShowAlert: true})
		return
	}

	mention := fmt.Sprintf("[%s](tg://user?id=%d)", c.Sender.FirstName, c.Sender.ID)

	switch command {
	case "Pause":
		if !utils.IsMusicPlaying(chatID) {
			b.Respond(c, &tb.CallbackResponse{Text: "⚠️ Stream is already paused.", ShowAlert: true})
			return
		}
		b.Respond(c, &tb.CallbackResponse{})
		utils.MusicOff(chatID)
		utils.PauseStream(chatID) // Call to Call Manager

		b.Send(c.Message.Chat, fmt.Sprintf("<blockquote><b>⏸ ➻ sᴛʀєᴧϻ ᴘᴧυsєᴅ 🥀</b>\n│ \n└<b>ʙʏ :</b> %s</blockquote>", mention), tb.ModeHTML)

	case "Resume":
		if utils.IsMusicPlaying(chatID) {
			b.Respond(c, &tb.CallbackResponse{Text: "⚠️ Stream is already playing.", ShowAlert: true})
			return
		}
		b.Respond(c, &tb.CallbackResponse{})
		utils.MusicOn(chatID)
		utils.ResumeStream(chatID)

		b.Send(c.Message.Chat, fmt.Sprintf("<blockquote><b>▶️ ➻ sᴛʀєᴧϻ ʀєsυϻєᴅ 🥀</b>\n│ \n└<b>ʙʏ :</b> %s</blockquote>", mention), tb.ModeHTML)

	case "Stop", "End":
		b.Respond(c, &tb.CallbackResponse{})
		utils.StopStream(chatID)
		utils.SetLoop(chatID, 0)

		b.Send(c.Message.Chat, fmt.Sprintf("<blockquote><b>⏹ ➻ sᴛʀєᴧϻ єηᴅєᴅ 🥀</b>\n│ \n└<b>ʙʏ :</b> %s</blockquote>", mention), tb.ModeHTML)
		b.Delete(c.Message)

	case "Skip":
		b.Respond(c, &tb.CallbackResponse{})
		utils.SkipStream(chatID) // Implement queue popping and playing next in utils
		b.Edit(c.Message, fmt.Sprintf("<blockquote><b>⏭ ➻ sᴛʀєᴧϻ sᴋɪᴘᴘєᴅ 🥀</b>\n│ \n└<b>ʙʏ :</b> %s</blockquote>", mention), tb.ModeHTML)

	case "Autoplay":
		state := utils.IsAutoplayOn(chatID)
		if state {
			utils.AutoplayOff(chatID)
			b.Respond(c, &tb.CallbackResponse{Text: "🔴 Autoplay Disabled!", ShowAlert: true})
			b.Send(c.Message.Chat, fmt.Sprintf("<blockquote><b>🔴 🎧 Ʌυᴛσᴘʟᴧʏ sʏsᴛєϻ</b>\n\n<b>Ʌυᴛσᴘʟᴧʏ ғσʀ ᴛʜɪs ɢʀσυᴘ ɪs ησᴡ ᴅɪsᴧʙʟєᴅ 🔴.</b>\n└ <b>ʙʏ :</b> %s</blockquote>", mention), tb.ModeHTML)
		} else {
			utils.AutoplayOn(chatID)
			b.Respond(c, &tb.CallbackResponse{Text: "🟢 Autoplay Enabled!", ShowAlert: true})
			b.Send(c.Message.Chat, fmt.Sprintf("<blockquote><b>🟢 🎧 Ʌυᴛσᴘʟᴧʏ sʏsᴛєϻ</b>\n\n<b>Ʌυᴛσᴘʟᴧʏ ғσʀ ᴛʜɪs ɢʀσυᴘ ɪs ησᴡ єηᴧʙʟєᴅ 🟢.</b>\n└ <b>ʙʏ :</b> %s</blockquote>", mention), tb.ModeHTML)
		}

	case "Thumb":
		// NEW: Thumbnail On/Off Feature Toggle
		state := utils.IsThumbnailEnabled(chatID)
		if state {
			utils.DisableThumbnail(chatID)
			b.Respond(c, &tb.CallbackResponse{Text: "🖼 Thumbnail Disabled!", ShowAlert: true})
			b.Send(c.Message.Chat, fmt.Sprintf("<blockquote><b>🖼 ➻ ᴛʜυϻʙηᴧɪʟ sʏsᴛєϻ</b>\n\n<b>ᴛʜυϻʙηᴧɪʟs ᴧʀє ησᴡ ᴅɪsᴧʙʟєᴅ ғσʀ ᴛʜɪs ɢʀσυᴘ.</b>\n└ <b>ʙʏ :</b> %s</blockquote>", mention), tb.ModeHTML)
		} else {
			utils.EnableThumbnail(chatID)
			b.Respond(c, &tb.CallbackResponse{Text: "🖼 Thumbnail Enabled!", ShowAlert: true})
			b.Send(c.Message.Chat, fmt.Sprintf("<blockquote><b>🖼 ➻ ᴛʜυϻʙηᴧɪʟ sʏsᴛєϻ</b>\n\n<b>ᴛʜυϻʙηᴧɪʟs ᴧʀє ησᴡ єηᴧʙʟєᴅ ғσʀ ᴛʜɪs ɢʀσυᴘ.</b>\n└ <b>ʙʏ :</b> %s</blockquote>", mention), tb.ModeHTML)
		}
	}
}

// MarkupTimer runs in the background to update progress bars every 5 minutes
func StartMarkupTimer(b *tb.Bot) {
	go func() {
		for {
			time.Sleep(300 * time.Second)
			activeChats := utils.GetActiveChats()

			for _, chatIDStr := range activeChats {
				chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

				if !utils.IsMusicPlaying(chatID) {
					continue
				}

				// Assuming utils provides a way to get the current playing message & stats
				playingData := utils.GetPlayingData(chatID)
				if playingData == nil || playingData.Seconds == 0 || playingData.MessageID == 0 {
					continue
				}

				// Update inline markup with new timer bar (Logic simplified for Go structure)
				chat, err := b.ChatByID(fmt.Sprintf("%d", chatID))
				if err == nil {
					msg := &tb.Message{ID: playingData.MessageID, Chat: chat}
					playedMin := utils.SecondsToMin(playingData.Played)

					newMarkup := StreamMarkupTimer(chatID, playedMin, playingData.Duration)
					b.Edit(msg, newMarkup)
				}
			}
		}
	}()
	log.Println("Background Markup Timer started.")
}
