package play

import (
	"fmt"
	"strings"

	"ANJALI/utils/database"
	"ANJALI/utils/decorators"
	"ANJALI/utils/stream"

	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterLiveStreamHandlers(b *tb.Bot) {
	b.Handle("\fLiveStream", func(c *tb.Callback) {
		decorators.LanguageCB(b, c, func(b *tb.Bot, c *tb.Callback, lang string) {
			data := strings.Split(c.Data, " ")
			if len(data) < 2 {
				return
			}

			args := strings.Split(data[1], "|")
			if len(args) < 5 {
				return
			}

			vidID := args[0]
			userIDStr := args[1]
			mode := args[2]
			cplay := args[3]
			fplay := args[4]

			if fmt.Sprintf("%d", c.Sender.ID) != userIDStr {
				b.Respond(c, &tb.CallbackResponse{Text: "❌ This is not for you!", ShowAlert: true})
				return
			}

			var chatID = c.Message.Chat.ID
			channel := ""
			if cplay == "c" {
				linked := database.GetCMode(c.Message.Chat.ID)
				if linked == 0 {
					return
				}
				chatID = linked
				chat, _ := b.ChatByID(fmt.Sprintf("%d", linked))
				channel = chat.Title
			}

			b.Delete(c.Message)
			b.Respond(c)

			targetName := channel
			if channel == "" {
				targetName = "this group"
			}
			msg, _ := b.Send(c.Message.Chat, fmt.Sprintf("⏳ **Processing Live Stream for %s...**", targetName), tb.ModeMarkdown)

			// Execute actual Live Stream logic
			isVideo := false
			if mode == "v" {
				isVideo = true
			}

			forcePlay := false
			if fplay == "f" {
				forcePlay = true
			}

			// Mock payload for Stream function
			details := map[string]string{
				"vidid":        vidID,
				"title":        "Live Stream",
				"duration_min": "Live",
			}

			// FIXED: Removed 'err :=' because stream.Stream does not return any value
			stream.Stream(b, c.Message, chatID, c.Message.Chat.ID, details, isVideo, "live", forcePlay)
			b.Delete(msg)
		})
	})
}
