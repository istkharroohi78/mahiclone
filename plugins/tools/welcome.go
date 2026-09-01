package tools

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type CustomWelcome struct {
	Type   string
	Text   string
	FileID string
	Markup *tb.ReplyMarkup
}

var (
	welcomeState   = make(map[int64]bool)
	customWelcomes = make(map[int64]CustomWelcome)
	welTimeState   = make(map[int64]int)
	welMutex       sync.RWMutex
)

// Parses buttons formatted as: [Button Name](buttonurl:https://t.me/link)
func parseWelcomeButtons(text string) (string, *tb.ReplyMarkup) {
	re := regexp.MustCompile(`\[(.*?)\]\(buttonurl:(.*?)\)`)
	matches := re.FindAllStringSubmatch(text, -1)

	if len(matches) == 0 {
		return text, nil
	}

	markup := &tb.ReplyMarkup{}
	var btns []tb.Btn
	for _, match := range matches {
		btns = append(btns, markup.URL(match[1], match[2]))
	}

	markup.Inline(markup.Row(btns...))
	cleanText := re.ReplaceAllString(text, "")
	return strings.TrimSpace(cleanText), markup
}

func RegisterWelcomeHandlers(b *tb.Bot) {
	// /welcome [on|off]
	b.Handle("/welcome", func(m *tb.Message) {
		args := strings.Split(m.Text, " ")
		if len(args) < 2 {
			b.Send(m.Chat, "Usage: `/welcome [on|off]`", tb.ModeMarkdown)
			return
		}

		state := strings.ToLower(args[1])
		welMutex.Lock()
		defer welMutex.Unlock()

		if state == "on" {
			welcomeState[m.Chat.ID] = true
			b.Send(m.Chat, "✅ **Welcome notifications enabled for this group.**", tb.ModeMarkdown)
		} else {
			welcomeState[m.Chat.ID] = false
			b.Send(m.Chat, "❌ **Welcome notifications disabled for this group.**", tb.ModeMarkdown)
		}
	})

	// /set_welcome [Custom text/caption]
	b.Handle("/set_welcome", func(m *tb.Message) {
		rawText := strings.TrimPrefix(m.Text, "/set_welcome")
		if rawText == "" && m.ReplyTo != nil {
			rawText = m.ReplyTo.Text
		}

		if strings.TrimSpace(rawText) == "" {
			b.Send(m.Chat, `⚠️ **Usage Example:**
`+"`/set_welcome Welcome {mention} to {title}! Members: {count}\n[Support](buttonurl:https://t.me/channel)`", tb.ModeMarkdown)
			return
		}

		cleanText, markup := parseWelcomeButtons(rawText)

		welMutex.Lock()
		customWelcomes[m.Chat.ID] = CustomWelcome{
			Type:   "text",
			Text:   cleanText,
			Markup: markup,
		}
		welMutex.Unlock()

		b.Send(m.Chat, "✅ **Custom Welcome message saved successfully!**", tb.ModeMarkdown)
	})

	// /weltime [minutes]
	b.Handle("/weltime", func(m *tb.Message) {
		args := strings.Split(m.Text, " ")
		if len(args) < 2 {
			b.Send(m.Chat, "Usage: `/weltime [minutes]` (e.g. `/weltime 5`)", tb.ModeMarkdown)
			return
		}

		mins, err := strconv.Atoi(args[1])
		if err != nil || mins < 1 {
			b.Send(m.Chat, "❌ Invalid number of minutes provided.", tb.ModeMarkdown)
			return
		}

		welMutex.Lock()
		welTimeState[m.Chat.ID] = mins * 60
		welMutex.Unlock()

		b.Send(m.Chat, fmt.Sprintf("⏱️ **Welcome messages will now auto-delete after %d minutes.**", mins), tb.ModeMarkdown)
	})

	// Member Join Listener
	b.Handle(tb.OnUserJoined, func(m *tb.Message) {
		welMutex.RLock()
		enabled, exists := welcomeState[m.Chat.ID]
		custom, hasCustom := customWelcomes[m.Chat.ID]
		autoDelSec, hasDel := welTimeState[m.Chat.ID]
		welMutex.RUnlock()

		if exists && !enabled {
			return
		}

		if !hasDel {
			autoDelSec = 300 // 5 min default
		}

		mention := fmt.Sprintf("[%s](tg://user?id=%d)", m.UserJoined.FirstName, m.UserJoined.ID)
		welcomeText := fmt.Sprintf("👋 **Hey %s, Welcome to %s!**", mention, m.Chat.Title)

		var sentMsg *tb.Message
		if hasCustom {
			text := strings.ReplaceAll(custom.Text, "{mention}", mention)
			text = strings.ReplaceAll(text, "{title}", m.Chat.Title)
			text = strings.ReplaceAll(text, "{id}", fmt.Sprintf("%d", m.UserJoined.ID))

			if custom.Markup != nil {
				sentMsg, _ = b.Send(m.Chat, text, custom.Markup, tb.ModeMarkdown)
			} else {
				sentMsg, _ = b.Send(m.Chat, text, tb.ModeMarkdown)
			}
		} else {
			sentMsg, _ = b.Send(m.Chat, welcomeText, tb.ModeMarkdown)
		}

		if sentMsg != nil && autoDelSec > 0 {
			go func(msg *tb.Message, delay int) {
				time.Sleep(time.Duration(delay) * time.Second)
				b.Delete(msg)
			}(sentMsg, autoDelSec)
		}
	})
}
