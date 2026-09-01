package bot

import (
	"fmt"
	"math/rand"
	"strings"

	"ANJALI/config"
	"ANJALI/utils/decorators"
	"ANJALI/utils/inline"
	"ANJALI/utils/stuffs"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Fallback logic for Random Start Image
func getRandomStartImg() string {
	images := config.LoadConfig().StartImgURL
	if len(images) > 0 {
		return images[rand.Intn(len(images))]
	}
	return "https://telegra.ph/file/2e3d368e77c449c287430.jpg"
}

func RegisterHelpHandlers(b *tb.Bot) {
	// Command: /help (Private)
	b.Handle("/help", func(m *tb.Message) {
		if m.Chat.Type != tb.ChatPrivate {
			return
		}
		decorators.LanguageStart(b, m, func(b *tb.Bot, m *tb.Message, lang string) {
			b.Delete(m)
			loc := map[string]string{ /* Load locale here */ }

			markup := inline.HelpPannel(loc)
			photo := &tb.Photo{
				File:    tb.FromURL(getRandomStartImg()),
				Caption: fmt.Sprintf("Hi! Click the buttons below for help.\n\nSupport: %s", config.LoadConfig().SupportChat),
			}

			// Simulate has_spoiler for photo by sending as document/animation or relying on custom clients
			b.Send(m.Chat, photo, markup, tb.ModeHTML)
		})
	})

	// Command: /help (Group)
	b.Handle("/help", func(m *tb.Message) {
		if m.Chat.Type == tb.ChatPrivate {
			return
		}
		decorators.LanguageStart(b, m, func(b *tb.Bot, m *tb.Message, lang string) {
			loc := map[string]string{ /* Load locale here */ }
			markup := inline.PrivateHelpPanel(loc, b.Me.Username)
			b.Send(m.Chat, "Contact me in PM for help.", markup, tb.ModeHTML)
		})
	})

	// Callbacks
	b.Handle("\fhelp_callback", func(c *tb.Callback) {
		decorators.LanguageCB(b, c, func(b *tb.Bot, c *tb.Callback, lang string) {
			data := strings.Split(c.Data, " ")
			if len(data) < 2 {
				return
			}
			cb := data[1]
			loc := map[string]string{ /* Load locale here */ }
			markup := inline.HelpBackMarkup(loc)

			var helpText string
			switch cb {
			case "hb1":
				helpText = stuffs.HelpChatGPT
			case "hb2":
				helpText = stuffs.HelpSticker
			case "hb3":
				helpText = stuffs.HelpTagAll
			case "hb4":
				helpText = stuffs.HelpInfo
			case "hb5":
				helpText = stuffs.HelpGroup
			case "hb6":
				helpText = stuffs.HelpExtra
			case "hb7":
				helpText = stuffs.HelpImage
			case "hb8":
				helpText = stuffs.HelpAction
			case "hb9":
				helpText = stuffs.HelpSearch
			case "hb10":
				helpText = stuffs.HelpFont
			case "hb11":
				helpText = stuffs.HelpGame
			case "hb12":
				helpText = stuffs.HelpTG
			case "hb13":
				helpText = stuffs.HelpImposter
			case "hb14":
				helpText = stuffs.HelpTD
			case "hb15":
				helpText = stuffs.HelpFun
			}

			if helpText != "" {
				b.Edit(c.Message, helpText, markup, tb.ModeHTML)
			}
		})
	})
}
