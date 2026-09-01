package cplugin

import (
	"math/rand"
	"time"

	"ANJALI/config"

	tb "gopkg.in/tucnak/telebot.v2"
)

const BotLink = "https://t.me/clone_MUSICrobot"
const FallbackImage = "https://d.uguu.se/AnYfJXGx.jpg"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetRandomStartImg safely fetches a random start image from the config
func GetRandomStartImg() string {
	cfg := config.LoadConfig()
	if len(cfg.StartImgURLs) > 0 {
		return cfg.StartImgURLs[rand.Intn(len(cfg.StartImgURLs))]
	}
	return FallbackImage
}

// CloneCommand handles the /clone command
func CloneCommand(b *tb.Bot, m *tb.Message) {
	styleMap := GetStyleMap()

	// Create the markup with dynamic color from styleMap
	markup := &tb.ReplyMarkup{}
	btnClone := CreateBtn(markup, "ɢᴏ ᴀɴᴅ ᴄʟᴏɴᴇ", "", BotLink, styleMap[1], false)
	markup.Inline(markup.Row(btnClone))

	photo := &tb.Photo{
		File:    tb.FromURL(GetRandomStartImg()),
		Caption: "🌟 **Clone Bot System**\n\nCreate your own music bot easily! Click the button below to get started.",
	}

	_, err := b.Send(m.Chat, photo, markup, tb.ModeMarkdown)
	if err != nil {
		b.Send(m.Chat, "❌ Failed to send clone menu.")
	}
}
