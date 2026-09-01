package sudo

import (
	"fmt"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	// Apne project module ka naam 'ANJALI' lagayein
	mystrings "ANJALI/i18n"
	"ANJALI/utils/database"
	"ANJALI/config" 
)

// isSudo checks if the user ID exists in the Sudoers list
func isSudo(userID int64) bool {
	for _, id := range config.Sudoers {
		if id == userID {
			return true
		}
	}
	return false
}

// contains checks if a chatID is in the blacklisted array
func contains(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func RegisterBlChatHandlers(b *tb.Bot) {

	// 1. Blacklist Chat Command
	blChatHandler := func(m *tb.Message) {
		if !isSudo(m.Sender.ID) {
			return // Ignore command if user is not SUDO
		}

		langData := mystrings.GetString(database.GetLang(m.Chat.ID))
		args := strings.Split(strings.TrimSpace(m.Text), " ")

		if len(args) != 2 {
			b.Send(m.Chat, langData["black_1"])
			return
		}

		chatID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			b.Send(m.Chat, "⚠️ Invalid Chat ID format. Must be an integer.\n\n— the shiv")
			return
		}

		blacklistedChats := database.GetBlacklistedChats()
		if contains(blacklistedChats, chatID) {
			b.Send(m.Chat, langData["black_2"])
			return
		}

		success := database.BlacklistChat(chatID)
		if success {
			b.Send(m.Chat, langData["black_3"])
		} else {
			b.Send(m.Chat, langData["black_9"])
		}

		// Bot attempts to leave the blacklisted chat
		targetChat := &tb.Chat{ID: chatID}
		b.Leave(targetChat)
	}

	b.Handle("/blchat", blChatHandler)
	b.Handle("/blacklistchat", blChatHandler)


	// 2. Whitelist Chat Command
	wlChatHandler := func(m *tb.Message) {
		if !isSudo(m.Sender.ID) {
			return
		}

		langData := mystrings.GetString(database.GetLang(m.Chat.ID))
		args := strings.Split(strings.TrimSpace(m.Text), " ")

		if len(args) != 2 {
			b.Send(m.Chat, langData["black_4"])
			return
		}

		chatID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			b.Send(m.Chat, "⚠️ Invalid Chat ID format.\n\n— the shiv")
			return
		}

		blacklistedChats := database.GetBlacklistedChats()
		if !contains(blacklistedChats, chatID) {
			b.Send(m.Chat, langData["black_5"])
			return
		}

		success := database.WhitelistChat(chatID)
		if success {
			b.Send(m.Chat, langData["black_6"])
		} else {
			b.Send(m.Chat, langData["black_9"])
		}
	}

	b.Handle("/whitelistchat", wlChatHandler)
	b.Handle("/unblacklistchat", wlChatHandler)
	b.Handle("/unblchat", wlChatHandler)


	// 3. List All Blacklisted Chats
	blChatsHandler := func(m *tb.Message) {
		langData := mystrings.GetString(database.GetLang(m.Chat.ID))
		blacklistedChats := database.GetBlacklistedChats()

		if len(blacklistedChats) == 0 {
			// Formatting the string to insert the bot's username (app.mention equivalent)
			formattedMsg := strings.Replace(langData["black_8"], "{0}", b.Me.Username, -1)
			b.Send(m.Chat, formattedMsg)
			return
		}

		text := langData["black_7"] + "\n"
		for count, chatID := range blacklistedChats {
			title := "ᴘʀɪᴠᴀᴛᴇ ᴄʜᴀᴛ"
			
			// Trying to fetch the chat title using Telegram API
			chat, err := b.ChatByID(fmt.Sprintf("%d", chatID))
			if err == nil && chat.Title != "" {
				title = chat.Title
			}

			text += fmt.Sprintf("%d. %s [<code>%d</code>]\n", count+1, title, chatID)
		}

		b.Send(m.Chat, text, &tb.SendOptions{ParseMode: tb.ModeHTML})
	}

	b.Handle("/blchats", blChatsHandler)
	b.Handle("/blacklistedchats", blChatsHandler)
}
