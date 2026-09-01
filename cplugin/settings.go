package cplugin

import (
	"fmt"
	"strings"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// SettingsCommand handles /settings and /setting
func SettingsCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	botName := b.Me.FirstName
	text := fmt.Sprintf("⚙️ **Settings for %s**\n\nGroup: %s (`%d`)", botName, m.Chat.Title, m.Chat.ID)

	markup := utils.SettingMarkup()
	b.Send(m.Chat, text, markup, tb.ModeMarkdown)
}

// SettingsCallback handles all settings-related inline buttons
func SettingsCallback(b *tb.Bot, c *tb.Callback) {
	data := strings.TrimSpace(c.Data)

	// Settings Back Helper
	if data == "settings_helper" || data == "settingsback_helper" {
		if c.Message.Chat.Type == tb.ChatPrivate {
			b.Respond(c, &tb.CallbackResponse{})
			markup := utils.PrivatePanel()
			b.Edit(c.Message, fmt.Sprintf("Welcome back to Settings, %s!", c.Sender.FirstName), markup)
			return
		}

		markup := utils.SettingMarkup()
		b.Edit(c.Message, fmt.Sprintf("⚙️ **Settings for %s**\n\nGroup: %s", b.Me.FirstName, c.Message.Chat.Title), markup)
		b.Respond(c, &tb.CallbackResponse{Text: "Opening Settings..."})
		return
	}

	// Informational Callbacks (No changes made, just alerts)
	infoResponses := map[string]string{
		"SEARCHANSWER":   "This changes the search mode (Inline/Direct).",
		"PLAYMODEANSWER": "This changes the play mode.",
		"PLAYTYPEANSWER": "This restricts who can play music (Admins/Everyone).",
		"AUTHANSWER":     "Authorized Users Management.",
		"VOTEANSWER":     "Vote Mode settings.",
	}

	if responseText, exists := infoResponses[data]; exists {
		b.Respond(c, &tb.CallbackResponse{Text: responseText, ShowAlert: true})
		return
	}

	// Active State Changers
	if !isAdminCheck(b, c.Message.Chat, int64(c.Sender.ID)) {
		b.Respond(c, &tb.CallbackResponse{Text: "❌ Admin rights needed!", ShowAlert: true})
		return
	}

	chatID := c.Message.Chat.ID
	var newMarkup *tb.ReplyMarkup

	switch data {
	case "MODECHANGE":
		current := utils.GetPlaymode(chatID)
		if current == "Direct" {
			utils.SetPlaymode(chatID, "Inline")
		} else {
			utils.SetPlaymode(chatID, "Direct")
		}
		newMarkup = utils.PlaymodeUsersMarkup(chatID)
		b.Respond(c, &tb.CallbackResponse{Text: "Playmode Updated!", ShowAlert: true})

	case "PLAYTYPECHANGE":
		current := utils.GetPlaytype(chatID)
		if current == "Everyone" {
			utils.SetPlaytype(chatID, "Admin")
		} else {
			utils.SetPlaytype(chatID, "Everyone")
		}
		newMarkup = utils.PlaymodeUsersMarkup(chatID)
		b.Respond(c, &tb.CallbackResponse{Text: "Playtype Updated!", ShowAlert: true})

	case "VOMODECHANGE":
		if utils.IsSkipmode(chatID) {
			utils.SkipOff(chatID)
		} else {
			utils.SkipOn(chatID)
		}
		currentVotes := utils.GetUpvoteCount(chatID)
		newMarkup = utils.VoteModeMarkup(currentVotes, utils.IsSkipmode(chatID))
		b.Respond(c, &tb.CallbackResponse{Text: "Vote Mode Updated!", ShowAlert: true})

	case "AUTH":
		if utils.IsNonAdminChat(chatID) {
			utils.RemoveNonAdminChat(chatID)
		} else {
			utils.AddNonAdminChat(chatID)
		}
		newMarkup = utils.AuthUsersMarkup(chatID)
		b.Respond(c, &tb.CallbackResponse{Text: "Auth Settings Updated!", ShowAlert: true})
	}

	if newMarkup != nil {
		b.Edit(c.Message, newMarkup)
	}
}
