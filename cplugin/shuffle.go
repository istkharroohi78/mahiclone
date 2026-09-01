package cplugin

import (
	"fmt"
	"math/rand"
	"time"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ShuffleCommand handles /shuffle and /cshuffle
func ShuffleCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	chatID := m.Chat.ID

	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Sirf Admins he is command ko use kar sakte hain!**", tb.ModeMarkdown)
		return
	}

	queue := utils.GetQueue(chatID)
	if len(queue) == 0 {
		b.Send(m.Chat, "⚠️ **Queue is empty!**", tb.ModeMarkdown)
		return
	}

	if len(queue) == 1 {
		b.Send(m.Chat, "⚠️ **Not enough tracks in queue to shuffle!**", utils.CloseMarkup(), tb.ModeMarkdown)
		return
	}

	// Extract the currently playing track
	playingTrack := queue[0]
	remainingQueue := queue[1:]

	// Shuffle the rest of the queue
	rand.Shuffle(len(remainingQueue), func(i, j int) {
		remainingQueue[i], remainingQueue[j] = remainingQueue[j], remainingQueue[i]
	})

	// Reassemble queue
	newQueue := append([]utils.QueueItem{playingTrack}, remainingQueue...)
	utils.SetQueue(chatID, newQueue)

	text := fmt.Sprintf("🔀 **Queue Shuffled**\n\n**Admin:** [%s](tg://user?id=%d)", m.Sender.FirstName, m.Sender.ID)
	b.Send(m.Chat, text, utils.CloseMarkup(), tb.ModeMarkdown)
}
