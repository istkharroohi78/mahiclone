package cplugin

import (
	"fmt"
	"strconv"
	"strings"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// LiveStreamCallback handles 'LiveStream' inline buttons
func LiveStreamCallback(b *tb.Bot, c *tb.Callback) {
	data := strings.Split(strings.TrimSpace(c.Data), " ")
	if len(data) < 2 {
		return
	}

	args := strings.Split(data[1], "|")
	if len(args) < 5 {
		return
	}

	vidID, userIDStr, mode, cplay, fplay := args[0], args[1], args[2], args[3], args[4]
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	if int64(c.Sender.ID) != userID {
		b.Respond(c, &tb.CallbackResponse{Text: "❌ This is not for you!", ShowAlert: true})
		return
	}

	chat, channelName := utils.GetChannelPlayCB(cplay, c.Message.Chat.ID)

	video := false
	if mode == "v" {
		video = true
	}

	b.Delete(c.Message)
	b.Respond(c, &tb.CallbackResponse{})

	msgText := "🔄 Processing live stream..."
	if channelName != "" {
		msgText = fmt.Sprintf("🔄 Processing live stream for %s...", channelName)
	}
	mystic, _ := b.Send(c.Message.Chat, msgText)

	// Fetch track details
	details, err := utils.GetYouTubeTrack(vidID, true)
	if err != nil {
		b.Edit(mystic, "❌ Error extracting track details.")
		return
	}

	ffplay := false
	if fplay == "f" {
		ffplay = true
	}

	if details.DurationMin == "0" || details.DurationMin == "" {
		err := utils.StreamLive(mystic, userID, details, chat, c.Sender.FirstName, c.Message.Chat.ID, video, "live", ffplay)
		if err != nil {
			b.Edit(mystic, fmt.Sprintf("❌ Stream Error: %v", err))
		}
	} else {
		b.Edit(mystic, "» ɴᴏᴛ ᴀ ʟɪᴠᴇ sᴛʀᴇᴀᴍ.")
	}
	b.Delete(mystic)
}
