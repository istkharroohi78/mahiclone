package tools

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	// Apne project module ka naam 'ANJALI' lagayein
	mystrings "ANJALI/i18n"
	"ANJALI/utils/database"
)

// languagesKeyboard dynamically creates the language selection menu in 2 columns
func languagesKeyboard(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	var rows []tb.Row
	var currentRow []tb.Btn

	// Iterate over available languages loaded in strings map
	for code := range mystrings.Languages {
		// Callback data format: "languages:en"
		btn := menu.Data(strings.ToUpper(code), "languages", code)
		currentRow = append(currentRow, btn)

		// Create 2 columns (row_width=2)
		if len(currentRow) == 2 {
			rows = append(rows, menu.Row(currentRow...))
			currentRow = nil
		}
	}
	// Append remaining button if odd number of languages
	if len(currentRow) > 0 {
		rows = append(rows, menu.Row(currentRow...))
	}

	// Add Back and Close buttons
	backBtn := menu.Data(langData["BACK_BUTTON"], "settingsback_helper")
	closeBtn := menu.Data(langData["CLOSE_BUTTON"], "close")
	rows = append(rows, menu.Row(backBtn, closeBtn))

	menu.Inline(rows...)
	return menu
}

// RegisterLanguageHandlers binds the /lang commands and their callbacks
func RegisterLanguageHandlers(b *tb.Bot) {
	// 1. Command Handlers (/lang, /setlang, /language)
	langHandler := func(m *tb.Message) {
		chatID := m.Chat.ID
		oldLang := database.GetLang(chatID) // Fetch current chat language
		langData := mystrings.GetString(oldLang)

		menu := languagesKeyboard(langData)
		b.Send(m.Chat, langData["lang_1"], &tb.SendOptions{ReplyMarkup: menu})
	}

	b.Handle("/lang", langHandler)
	b.Handle("/setlang", langHandler)
	b.Handle("/language", langHandler)

	// 2. Universal Callback Handler for Pyrogram Regex equivalents
	b.Handle(tb.OnCallback, func(c *tb.Callback) {
		data := c.Data

		// Replicating filters.regex("LG")
		if strings.HasPrefix(data, "LG") {
			oldLang := database.GetLang(c.Message.Chat.ID)
			langData := mystrings.GetString(oldLang)
			
			menu := languagesKeyboard(langData)
			b.Edit(c.Message, langData["lang_1"], &tb.SendOptions{ReplyMarkup: menu})
			b.Respond(c, &tb.CallbackResponse{})
			return
		}

		// Replicating filters.regex(r"languages:(.*?)")
		// Note: Telebot button data comes as "uniqueData|payload". 
		// We handle both our custom "languages:code" and telebot's default pipe format.
		if strings.HasPrefix(data, "languages:") || strings.Contains(data, "\flanguages|") {
			
			// Extract language code
			var newLang string
			if strings.Contains(data, ":") {
				newLang = strings.Split(data, ":")[1]
			} else if strings.Contains(data, "|") {
				newLang = strings.Split(data, "|")[1]
			}

			if newLang == "" {
				return
			}

			oldLang := database.GetLang(c.Message.Chat.ID)
			oldLangData := mystrings.GetString(oldLang)

			// If language is already the same
			if oldLang == newLang {
				b.Respond(c, &tb.CallbackResponse{
					Text:      oldLangData["lang_4"],
					ShowAlert: true,
				})
				return
			}

			// Validate and set new language
			newLangData := mystrings.GetString(newLang)
			if len(newLangData) > 0 { // Successful fetch
				database.SetLang(c.Message.Chat.ID, newLang)
				
				b.Respond(c, &tb.CallbackResponse{
					Text:      newLangData["lang_2"],
					ShowAlert: true,
				})

				// Edit keyboard text to match the new language
				menu := languagesKeyboard(newLangData)
				b.Edit(c.Message, newLangData["lang_1"], &tb.SendOptions{ReplyMarkup: menu})
			} else {
				// Failed to fetch strings for the language
				b.Respond(c, &tb.CallbackResponse{
					Text:      oldLangData["lang_3"],
					ShowAlert: true,
				})
			}
		}
	})
}
