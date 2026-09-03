package inline

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Constants for button styles
const (
	Primary = iota
	Danger
	Success
)

// --- MISSING FUNCTIONS ADDED HERE ---

// CreateBtn creates a Telebot inline button dynamically (Fixes restart, stats, auth errors)
func CreateBtn(markup *tb.ReplyMarkup, text, data, url string, btnType int, query string, isCurrent bool) tb.Btn {
	if url != "" {
		return markup.URL(text, url)
	}
	return markup.Data(text, data)
}

// HelpPannel returns the main help menu markup (Fixes help.go and start.go errors)
func HelpPannel() *tb.ReplyMarkup {
	markup := &tb.ReplyMarkup{}
	markup.Inline(markup.Row(markup.Data("Back", "help_back")))
	return markup
}

// StyledButton helper used in your stream markups
func StyledButton(text, data, url string, btnType int) tb.Btn {
	m := &tb.ReplyMarkup{}
	if url != "" {
		return m.URL(text, url)
	}
	return m.Data(text, data)
}

// --- YOUR ORIGINAL CODE ---

// Option 1: Static Buttons
func GetStaticButtons(botUsername string) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}

	btnResume := StyledButton("▷", "resume_cb", "", Success)
	btnPause := StyledButton("II", "pause_cb", "", Danger)
	btnSkip := StyledButton("‣‣I", "skip_cb", "", Primary)
	btnEnd := StyledButton("▢", "end_cb", "", Danger)

	// "Clone Me" hatakar "Add Me" set kiya gaya hai same premium font me
	btnAdd := StyledButton("『 ✦ 𝐀ᴅᴅ 𝐌є ✦ 』", "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), Success)

	m.Inline(
		m.Row(btnResume, btnPause, btnSkip, btnEnd),
		m.Row(btnAdd),
	)
	return m
}

// Close Key
func GetCloseKey() *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	btnClose := StyledButton("✯ CLOSE ✯", "close", "", Danger)

	m.Inline(
		m.Row(btnClose),
	)
	return m
}

// Option 2: Dynamic Stream Markup (Recommended)
func StreamMarkup(chatId int64, botUsername string) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}

	// Top Row: Basic Controls
	btnResume := StyledButton("▷", fmt.Sprintf("ADMIN Resume|%d", chatId), "", Success)
	btnPause := StyledButton("II", fmt.Sprintf("ADMIN Pause|%d", chatId), "", Danger)
	btnSkip := StyledButton("‣‣I", fmt.Sprintf("ADMIN Skip|%d", chatId), "", Primary)
	btnEnd := StyledButton("▢", fmt.Sprintf("ADMIN Stop|%d", chatId), "", Danger)

	// Middle Row: Seek Buttons
	btnSeekBack := StyledButton("<- 20s", fmt.Sprintf("ADMIN SeekBack|%d", chatId), "", Primary)
	btnSeekFwd := StyledButton("20s + ->", fmt.Sprintf("ADMIN SeekForward|%d", chatId), "", Primary)

	// Bottom Row: Add Me & Close merged
	btnAdd := StyledButton("『 𝐀ᴅᴅ 𝐌є 』", "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), Success)
	btnClose := StyledButton("✯ CLOSE ✯", "close", "", Danger)

	m.Inline(
		m.Row(btnResume, btnPause, btnSkip, btnEnd),
		m.Row(btnSeekBack, btnSeekFwd),
		m.Row(btnAdd, btnClose), // Space bachane ke liye merged
	)
	return m
}
