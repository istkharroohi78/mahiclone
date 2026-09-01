package cplugin

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	SpamThreshold     = 2
	SpamWindowSeconds = 5
)

// SpamMemory handles thread-safe spam protection
type SpamMemory struct {
	LastMessageTime map[int64]int64
	CommandCount    map[int64]int
	mu              sync.Mutex
}

var spamTracker = &SpamMemory{
	LastMessageTime: make(map[int64]int64),
	CommandCount:    make(map[int64]int),
}

// ChannelPlayCommand handles the /channelplay command
func ChannelPlayCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	// Admin verification (Using the helper defined in auth.go)
	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Sirf Admins he is command ko use kar sakte hain!**", tb.ModeMarkdown)
		return
	}

	userID := int64(m.Sender.ID)
	currentTime := time.Now().Unix()

	// 1. Spam Logic with Mutex Lock
	spamTracker.mu.Lock()
	lastMessageTime := spamTracker.LastMessageTime[userID]

	if currentTime-lastMessageTime < SpamWindowSeconds {
		spamTracker.LastMessageTime[userID] = currentTime
		spamTracker.CommandCount[userID]++

		if spamTracker.CommandCount[userID] > SpamThreshold {
			spamTracker.mu.Unlock()
			warning, _ := b.Send(m.Chat, fmt.Sprintf("**%s ᴘʟᴇᴀsᴇ ᴅᴏɴᴛ ᴅᴏ sᴘᴀᴍ, ᴀɴᴅ ᴛʀʏ ᴀɢᴀɪɴ ᴀғᴛᴇʀ 5 sᴇᴄ**", m.Sender.FirstName), tb.ModeMarkdown)
			time.Sleep(3 * time.Second)
			b.Delete(warning)
			return
		}
	} else {
		spamTracker.CommandCount[userID] = 1
		spamTracker.LastMessageTime[userID] = currentTime
	}
	spamTracker.mu.Unlock()

	// 2. Command Logic
	args := strings.SplitN(m.Text, " ", 3)
	if len(args) < 2 {
		b.Send(m.Chat, fmt.Sprintf("⚠️ **Usage:** `/channelplay [disable | linked | @channelusername]`"), tb.ModeMarkdown)
		return
	}

	query := strings.ToLower(strings.TrimSpace(args[1]))

	// Option 1: Disable
	if query == "disable" {
		utils.SetCMode(m.Chat.ID, 0)
		b.Send(m.Chat, "✅ **Channel Play Mode has been disabled.**", tb.ModeMarkdown)
		return
	}

	// Option 2: Linked Chat
	if query == "linked" {
		// Telebot doesn't natively expose linked chats easily without raw API calls,
		// Assuming utils.GetLinkedChatID handles the raw API fetch or bot configuration.
		linkedChatID := utils.GetLinkedChatID(b, m.Chat.ID)
		if linkedChatID != 0 {
			utils.SetCMode(m.Chat.ID, linkedChatID)
			b.Send(m.Chat, fmt.Sprintf("✅ **Channel Play Mode enabled for linked chat ID:** `%d`", linkedChatID), tb.ModeMarkdown)
		} else {
			b.Send(m.Chat, "❌ **No linked channel found for this group.**", tb.ModeMarkdown)
		}
		return
	}

	// Option 3: Custom Channel
	chat, err := b.ChatByID(query)
	if err != nil {
		b.Send(m.Chat, "❌ **Channel not found.** Make sure the bot is an admin in the channel.", tb.ModeMarkdown)
		return
	}

	if chat.Type != tb.Channel {
		b.Send(m.Chat, "❌ **The specified chat is not a channel.**", tb.ModeMarkdown)
		return
	}

	// Verify Ownership
	admins, err := b.AdminsOf(chat)
	if err != nil {
		b.Send(m.Chat, "❌ **Could not fetch channel admins.**", tb.ModeMarkdown)
		return
	}

	isOwner := false
	var ownerUsername string

	for _, admin := range admins {
		if admin.Role == tb.Creator { // Creator == Owner
			if admin.User.ID == m.Sender.ID {
				isOwner = true
			}
			ownerUsername = admin.User.Username
			break
		}
	}

	if !isOwner {
		b.Send(m.Chat, fmt.Sprintf("❌ **You are not the owner of %s.** Ask @%s to run this.", chat.Title, ownerUsername), tb.ModeMarkdown)
		return
	}

	// Success
	utils.SetCMode(m.Chat.ID, chat.ID)
	b.Send(m.Chat, fmt.Sprintf("✅ **Channel Play Mode successfully set for** %s (`%d`)", chat.Title, chat.ID), tb.ModeMarkdown)
}
