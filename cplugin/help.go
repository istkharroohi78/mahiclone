package cplugin

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"ANJALI/config"
	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

const FallbackHelpImg = "https://files.catbox.moe/10zwqs.jpg"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetRandomHelpImg() string {
	cfg := config.LoadConfig()
	if len(cfg.HelpImgURLs) > 0 {
		return cfg.HelpImgURLs[rand.Intn(len(cfg.HelpImgURLs))]
	}
	return FallbackHelpImg
}

// HelpCommand handles /help in private and groups
func HelpCommand(b *tb.Bot, m *tb.Message) {
	botID := b.Me.ID
	cloneData := utils.GetCloneBotData(botID)
	cfg := config.LoadConfig()

	isOwner := false
	if cloneData != nil && int64(m.Sender.ID) == cloneData.UserID {
		isOwner = true
	}

	supportChat := cfg.SupportChat
	if cloneData != nil && cloneData.SupportChat != "" {
		supportChat = fmt.Sprintf("https://t.me/%s", cloneData.SupportChat)
	}

	helpText := fmt.Sprintf("Welcome to the Help Menu! For more assistance, visit our [Support Chat](%s).", supportChat)
	keyboard := utils.FirstPageMarkup(isOwner)

	if m.Private() {
		photo := &tb.Photo{File: tb.FromURL(GetRandomHelpImg()), Caption: helpText}
		b.Send(m.Chat, photo, keyboard, tb.ModeMarkdown)
	} else {
		// Group
		keyboard = utils.PrivateHelpPanel()
		b.Send(m.Chat, "Contact me in PM for help.", keyboard, tb.ModeMarkdown)
	}
}

// HelpCallback handles all help navigation callbacks
func HelpCallback(b *tb.Bot, c *tb.Callback) {
	cb := strings.Split(c.Data, " ")
	if len(cb) < 2 {
		return
	}
	action := cb[1]

	keyboard := utils.HelpBackMarkup()
	cloneBackKb := utils.CloneBackMarkup()

	var newText string
	var newMarkup *tb.ReplyMarkup

	switch action {
	case "hb1":
		newText, newMarkup = utils.HelpText1, keyboard
	case "hb2":
		newText, newMarkup = utils.HelpText2, keyboard
	// Add cases hb3 to hb15 mapped to their respective strings
	case "chelp":
		newText, newMarkup = utils.CloneHelpMenu, utils.CloneHelpPanel()
	case "clone_manage":
		newText, newMarkup = utils.CloneManageText, cloneBackKb
	case "clone_start":
		newText, newMarkup = utils.CloneStartText, cloneBackKb
	case "clone_ping":
		newText, newMarkup = utils.ClonePingText, cloneBackKb
	case "clone_buttons":
		newText, newMarkup = utils.CloneButtonsText, cloneBackKb
	case "clone_play":
		newText, newMarkup = utils.ClonePlayModeText, cloneBackKb
	case "clone_logger":
		newText, newMarkup = utils.CloneLoggerText, cloneBackKb
	default:
		return
	}

	b.Edit(c.Message, newText, newMarkup, tb.ModeMarkdown)
	b.Respond(c, &tb.CallbackResponse{})
}
