package utils

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

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
