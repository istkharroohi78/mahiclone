package bot

import (
	"fmt"
	"strings"
	"time"

	"ANJALI/config"
	"ANJALI/utils/database"
	"ANJALI/utils/decorators"
	"ANJALI/utils/inline"

	tb "gopkg.in/tucnak/telebot.v2"
)

var BootTime = time.Now()

// Helper function to prevent undefined error
func getRandomStartImg() string {
	return "https://files.catbox.moe/6r97s4.jpg" 
}

func RegisterStartHandlers(b *tb.Bot) {
	// Combined Command: /start (Handles both Private and Group)
	b.Handle("/start", func(m *tb.Message) {
		decorators.LanguageStart(b, m, func(b *tb.Bot, m *tb.Message, lang string) {
			
			if m.Chat.Type == tb.ChatPrivate {
				// --- PRIVATE CHAT START ---
				database.AddServedUser(int64(m.Sender.ID))

				args := strings.Split(m.Text, " ")
				if len(args) > 1 {
					// Handle Deep Links (e.g. /start help, /start info_vidid)
					query := args[1]
					if strings.HasPrefix(query, "help") {
						// Fixed HelpPannel argument error
						markup := inline.HelpPannel()
						photo := &tb.Photo{File: tb.FromURL(getRandomStartImg()), Caption: "Help Menu"}
						b.Send(m.Chat, photo, markup)
						return
					}
					// Other deep links like info, sudo, etc.
				} else {
					// Normal Start
					cfg := config.LoadConfig()
					// Removed unused 'loc' variable and fixed GithubURL
					markup := inline.PrivatePanel(b.Me.Username, cfg.SupportChat, cfg.SupportChannel, "https://github.com", cfg.OwnerID)

					photo := &tb.Photo{
						File:    tb.FromURL(getRandomStartImg()),
						Caption: fmt.Sprintf("Hello %s, I am %s!\n\nI can play music in voice chats.", m.Sender.FirstName, b.Me.FirstName),
					}
					b.Send(m.Chat, photo, markup, tb.ModeHTML)
				}
			} else {
				// --- GROUP CHAT START ---
				database.AddServedChat(m.Chat.ID)
				cfg := config.LoadConfig()
				markup := inline.StartPanel(b.Me.Username, cfg.SupportChat)

				photo := &tb.Photo{
					File:    tb.FromURL(getRandomStartImg()),
					Caption: fmt.Sprintf("Hello! I am %s.", b.Me.FirstName),
				}
				b.Send(m.Chat, photo, markup, tb.ModeHTML)
			}
		})
	})
}
