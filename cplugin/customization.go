package cplugin

import (
	"fmt"
	"strings"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// IsCloneOwner checks if the user is the bot owner
func IsCloneOwner(botID int64, userID int64) bool {
	ownerID := utils.GetOwnerIDFromDB(botID)
	return userID == ownerID
}

// AddToRandomList appends a new item to the database list using a separator
func AddToRandomList(botID int64, typeKey string, newValue string) bool {
	currentData := utils.GetCloneSearchType(botID, typeKey)

	if currentData != "" {
		if strings.Contains(currentData, newValue) {
			return false
		}
		finalValue := fmt.Sprintf("%s|||%s", currentData, newValue)
		utils.SetCloneSearchType(botID, typeKey, finalValue)
	} else {
		utils.SetCloneSearchType(botID, typeKey, newValue)
	}
	return true
}

// SetPlayText handles /setplaytext and /addplaytext
func SetPlayText(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Only the Bot Owner can change these settings.**", tb.ModeMarkdown)
		return
	}

	args := strings.SplitN(m.Text, " ", 2)
	if len(args) < 2 {
		b.Send(m.Chat, "Usage: `/setplaytext <Text/Emoji>`", tb.ModeMarkdown)
		return
	}

	text := args[1]
	AddToRandomList(b.Me.ID, "text", text)
	b.Send(m.Chat, fmt.Sprintf("✅ **Added to Random List:**\n\n%s", text), tb.ModeMarkdown)
}

// SetPlayMedia handles setting Stickers, Animations (GIFs), Videos, and Photos
func SetPlayMedia(b *tb.Bot, m *tb.Message, mediaType string) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Only the Bot Owner can change these settings.**", tb.ModeMarkdown)
		return
	}

	if m.ReplyTo == nil {
		b.Send(m.Chat, fmt.Sprintf("Usage: Reply to a %s with the command.", strings.Title(mediaType)), tb.ModeMarkdown)
		return
	}

	var fileID string
	switch mediaType {
	case "sticker":
		if m.ReplyTo.Sticker != nil {
			fileID = m.ReplyTo.Sticker.FileID
		}
	case "animation":
		if m.ReplyTo.Animation != nil {
			fileID = m.ReplyTo.Animation.FileID
		}
	case "video":
		if m.ReplyTo.Video != nil {
			fileID = m.ReplyTo.Video.FileID
		}
	case "photo":
		if m.ReplyTo.Photo != nil {
			fileID = m.ReplyTo.Photo.FileID
		}
	}

	if fileID == "" {
		b.Send(m.Chat, fmt.Sprintf("❌ Please reply to a valid %s.", mediaType), tb.ModeMarkdown)
		return
	}

	AddToRandomList(b.Me.ID, mediaType, fileID)

	// Auto-Hide Text if adding visual media
	if mediaType != "sticker" {
		currentText := utils.GetCloneSearchType(b.Me.ID, "text")
		if currentText == "" {
			utils.SetCloneSearchType(b.Me.ID, "text", "⠀")
		}
	}

	b.Send(m.Chat, fmt.Sprintf("✅ **%s Added to Random List!**", strings.Title(mediaType)), tb.ModeMarkdown)
}

// SetStreamText handles /setstreamtext
func SetStreamText(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Only the Bot Owner can change these settings.**", tb.ModeMarkdown)
		return
	}

	args := strings.SplitN(m.Text, " ", 2)
	if len(args) < 2 {
		helpText := "**Usage:** `/setstreamtext <Your Caption>`\n\n**Available Variables:**\n`{1}` : Song Name\n`{2}` : Duration\n`{3}` : Requested By\n\n**Example:**\n`/setstreamtext 🎸 Playing: {1} | ⏳ Time: {2}`"
		b.Send(m.Chat, helpText, tb.ModeMarkdown)
		return
	}

	text := args[1]
	utils.SetCloneStreamCaption(b.Me.ID, text)
	b.Send(m.Chat, fmt.Sprintf("✅ **Stream Caption Updated:**\n\n%s", text), tb.ModeMarkdown)
}

// ResetPlayMode handles /delplay, /resetplay, /delplaymode
func ResetPlayMode(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Only the Bot Owner can change these settings.**", tb.ModeMarkdown)
		return
	}
	utils.DeleteCloneSearchType(b.Me.ID)
	b.Send(m.Chat, "🗑️ **Search Mode Reset!**\nAll saved random lists cleared.", tb.ModeMarkdown)
}

// ResetStreamText handles /delstreamtext, /resetstreamtext
func ResetStreamText(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Only the Bot Owner can change these settings.**", tb.ModeMarkdown)
		return
	}
	utils.DeleteCloneStreamCaption(b.Me.ID)
	b.Send(m.Chat, "🗑️ **Stream Caption Reset!**", tb.ModeMarkdown)
}
