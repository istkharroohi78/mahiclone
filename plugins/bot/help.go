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

// Note: Ensure this function is DELETED from start.go to prevent the "redeclared" error.
func getRandomStartImg() string {
	images := config.LoadConfig().StartImgURL
	if len(images) > 0 {
		return images[rand.Intn(len(images))]
	}
	return "https://telegra.ph/file/2e3d368e77c449c287430.jpg"
}

func RegisterHelpHandlers(b *tb.Bot) {
	// Command: /help (Combined Group and Private logic)
	b.Handle("/help", func(m *tb.Message) {
		decorators.LanguageStart(b, m, func(b *tb.Bot, m *tb.Message, lang string) {
			loc := map[string]string{ /* Load locale here */ }

			if m.Chat.Type == tb.ChatPrivate {
				// Private Chat Logic
				b.Delete(m)
				
				// Fix: Changed "HelpPannel" to "HelpPanel"
				markup := inline.HelpPanel(loc) 
				
				photo := &tb.Photo{
					File:    tb.FromURL(getRandomStartImg()),
					Caption: fmt.Sprintf("Hi! Click the buttons below for help.\n\nSupport: %s", config.LoadConfig().SupportChat),
				}
				b.Send(m.Chat, photo, markup, tb.ModeHTML)
				
			} else {
				// Group Chat Logic
				markup := inline.PrivateHelpPanel(loc, b.Me.Username)
				b.Send(m.Chat, "Contact me in PM for help.", markup, tb.ModeHTML)
			}
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
