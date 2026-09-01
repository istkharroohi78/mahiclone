package cplugin

import (
	"fmt"
	"log"
	"strconv"

	// Aapke utils package ka import jahan database functions hain
	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// ActiveChatsCommand handles /ac, /activevc, /activevoice
func ActiveChatsCommand(b *tb.Bot, m *tb.Message) {
	botID := b.Me.ID
	userID := m.Sender.ID

	// 1. Fetch Owner from Database
	cloneData := utils.GetCloneBotData(botID)
	if cloneData == nil {
		return
	}

	ownerID := cloneData.UserID

	// STRICT OWNER CHECK
	if int64(userID) != ownerID {
		b.Send(m.Chat, "❌ **Only the Bot Owner can view these stats.**", tb.ModeMarkdown)
		return
	}

	waitingMsg, err := b.Send(m.Chat, "🔄 **Checking active groups...**")
	if err != nil {
		log.Println("Error sending wait message:", err)
		return
	}

	// DATA FETCHING
	globalAudio := utils.GetActiveChats()
	globalVideo := utils.GetActiveVideoChats()

	// Fetch Served Chats for this specific bot
	cloneServedChats := utils.GetServedChatsClone(botID)

	// Create a map for faster lookup (O(1) instead of O(N))
	myChatIDs := make(map[int64]bool)
	for _, chat := range cloneServedChats {
		myChatIDs[chat.ChatID] = true
	}

	// FILTERING AUDIO
	myAudioCount := 0
	for _, chatIDStr := range globalAudio {
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err == nil && myChatIDs[chatID] {
			myAudioCount++
		}
	}

	// FILTERING VIDEO
	myVideoCount := 0
	for _, chatIDStr := range globalVideo {
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err == nil && myChatIDs[chatID] {
			myVideoCount++
		}
	}

	// RESULT
	text := fmt.Sprintf("📊 **Bot Activity Status**\n\n"+
		"👤 **Owner:** [%s](tg://user?id=%d)\n"+
		"🤖 **Bot:** @%s\n\n"+
		"🏢 **Total Groups:** `%d`\n"+
		"🎧 **Active Audio:** `%d`\n"+
		"📹 **Active Video:** `%d`\n",
		m.Sender.FirstName, m.Sender.ID, b.Me.Username, len(myChatIDs), myAudioCount, myVideoCount)

	markup := &tb.ReplyMarkup{}
	closeBtn := markup.Data("✯ Close ✯", "close")
	markup.Inline(markup.Row(closeBtn))

	// Edit waiting message
	_, err = b.Edit(waitingMsg, text, markup, tb.ModeMarkdown)
	if err != nil {
		log.Println("Failed to edit activity message:", err)
	}
}
