package stream

import (
	"fmt"
	"log"

	"ANJALI/utils/database"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Stream processes the track and triggers playback
func Stream(b *tb.Bot, m *tb.Message, chatID, origChatID int64, result map[string]string, video bool, streamType string, forcePlay bool) {
	if len(result) == 0 {
		return
	}

	if forcePlay {
		// Stop current stream via PyTgCalls Wrapper
	}

	title := result["title"]
	duration := result["duration_min"]
	vidID := result["vidid"]
	filePath := result["filepath"] // Or fetch via Download func
	userName := m.Sender.FirstName

	if streamType == "youtube" {
		// Mock PyTgCalls JoinCall
		log.Printf("Joining Call: %d with file %s", chatID, filePath)

		PutQueue(chatID, origChatID, filePath, title, duration, userName, vidID, int64(m.Sender.ID), "video", forcePlay, "youtube")
		database.AddActiveChat(chatID)

		caption := fmt.Sprintf(`<blockquote><b>▶️ sᴛʀᴇᴀᴍɪɴɢ sᴛᴀʀᴛᴇᴅ</b>

<b>• ᴛɪᴛʟᴇ :</b> %s
<b>• ᴅᴜʀᴀᴛɪᴏɴ :</b> %s
<b>• ʀᴇǫᴜᴇsᴛᴇᴅ ʙʏ :</b> [%s](tg://user?id=%d)</blockquote>`, title, duration, userName, m.Sender.ID)

		// Create premium markup (imported from your buttons helper in prod)
		markup := &tb.ReplyMarkup{}
		markup.Inline(markup.Row(markup.Data("Pᴀᴜsᴇ", "ADMIN_Pause")))

		thumbURL := result["thumb"]
		if thumbURL == "" {
			thumbURL = "https://files.catbox.moe/6r97s4.jpg"
		}

		b.Send(m.Chat, &tb.Photo{File: tb.FromURL(thumbURL), Caption: caption}, markup, tb.ModeHTML)

	} else if streamType == "telegram" {
		PutQueue(chatID, origChatID, filePath, title, duration, userName, "", int64(m.Sender.ID), "video", forcePlay, "telegram")
		database.AddActiveChat(chatID)

		caption := fmt.Sprintf("▶️ **Started Playing Local Media**\n\n**Requested by:** %s", userName)
		b.Send(m.Chat, &tb.Photo{File: tb.FromURL("https://telegra.ph/file/2e3d368e77c449c287430.jpg"), Caption: caption}, tb.ModeMarkdown)
	}
}

func IncrementPlayed() {
	// Logic to count played songs
}
