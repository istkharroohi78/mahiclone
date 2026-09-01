package inline

import (
	"fmt"
	"math/rand"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	// Replace 'ANJALI' with your actual module path if different
	"ANJALI/config"
)

// getRandomEmoji returns a random custom emoji ID
func getRandomEmoji() string {
	rand.Seed(time.Now().UnixNano())
	return PremiumEmojis[rand.Intn(len(PremiumEmojis))]
}

// Smart Button Creator
// In Telebot, we don't have ButtonStyle directly, but we can set text and payload
func createBtn(menu *tb.ReplyMarkup, text string, cb string, url string) tb.Btn {
	// Telebot doesn't natively support setting custom emoji IDs within the button struct easily
	// but you can append the actual emoji character to the text if you have it.
	// For now, we will construct the button.
	if url != "" {
		return menu.URL(text, url)
	}
	// "unique" prefix is required by telebot for callbacks
	return menu.Data(text, cb)
}

func HelpPanel(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	first := []tb.Row{
		menu.Row(
			createBtn(menu, langData["H_B_1"], "help_callback_hb1", ""),
			createBtn(menu, langData["H_B_2"], "help_callback_hb2", ""),
			createBtn(menu, langData["H_B_3"], "help_callback_hb3", ""),
		),
		menu.Row(
			createBtn(menu, langData["H_B_4"], "help_callback_hb4", ""),
			createBtn(menu, langData["H_B_5"], "help_callback_hb5", ""),
			createBtn(menu, langData["H_B_6"], "help_callback_hb6", ""),
		),
		menu.Row(
			createBtn(menu, langData["H_B_7"], "help_callback_hb7", ""),
			createBtn(menu, langData["H_B_8"], "help_callback_hb8", ""),
			createBtn(menu, langData["H_B_9"], "help_callback_hb9", ""),
		),
		menu.Row(
			createBtn(menu, langData["H_B_10"], "help_callback_hb10", ""),
			createBtn(menu, langData["H_B_11"], "help_callback_hb11", ""),
			createBtn(menu, langData["H_B_12"], "help_callback_hb12", ""),
		),
		menu.Row(
			createBtn(menu, langData["H_B_13"], "help_callback_hb13", ""),
			createBtn(menu, langData["H_B_14"], "help_callback_hb14", ""),
			createBtn(menu, langData["H_B_15"], "help_callback_hb15", ""),
		),
		menu.Row(
			createBtn(menu, langData["BACK_BUTTON"], "settingsback_helper", ""),
		),
	}

	menu.Inline(first...)
	return menu
}

func FirstPage(langData map[string]string, isOwner bool) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	first := []tb.Row{
		menu.Row(
			createBtn(menu, langData["H_B_1"], "help_callback_hb1", ""),
			createBtn(menu, langData["H_B_2"], "help_callback_hb2", ""),
			createBtn(menu, langData["H_B_3"], "help_callback_hb3", ""),
		),
		menu.Row(
			createBtn(menu, langData["H_B_11"], "help_callback_hb11", ""),
			createBtn(menu, langData["H_B_8"], "help_callback_hb8", ""),
			createBtn(menu, langData["H_B_6"], "help_callback_hb6", ""),
		),
		menu.Row(
			createBtn(menu, langData["H_B_13"], "help_callback_hb13", ""),
			createBtn(menu, langData["H_B_12"], "help_callback_hb12", ""),
			createBtn(menu, langData["H_B_9"], "help_callback_cloghelp", ""),
		),
		menu.Row(
			createBtn(menu, langData["H_B_10"], "help_callback_hb10", ""),
			createBtn(menu, langData["H_B_14"], "help_callback_hb14", ""),
			createBtn(menu, langData["H_B_15"], "help_callback_hb15", ""),
		),
	}

	if isOwner {
		first = append(first, menu.Row(createBtn(menu, "🛠 ᴄʟᴏɴᴇ ғᴇᴀᴛᴜʀᴇ", "help_callback_chelp", "")))
	}

	first = append(first, menu.Row(createBtn(menu, langData["BACK_BUTTON"], "settingsback_home", "")))
	
	menu.Inline(first...)
	return menu
}

func CloneHelpPanel(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	buttons := []tb.Row{
		menu.Row(createBtn(menu, "ᴍᴀɴᴀɢᴇ", "help_callback_clone_manage", "")),
		menu.Row(
			createBtn(menu, "sᴛᴀʀᴛ", "help_callback_clone_start", ""),
			createBtn(menu, "ᴘɪɴɢ", "help_callback_clone_ping", ""),
		),
		menu.Row(
			createBtn(menu, "ᴘʟᴀʏ ᴍᴏᴅᴇ", "help_callback_clone_play", ""),
			createBtn(menu, "ʟᴏɢɢᴇʀ", "help_callback_clone_logger", ""),
		),
		menu.Row(createBtn(menu, "ʙᴜᴛᴛᴏɴs & ʀᴇɴᴀᴍᴇ", "help_callback_clone_buttons", "")),
		menu.Row(createBtn(menu, langData["BACK_BUTTON"], "settings_back_helper", "")),
	}
	
	menu.Inline(buttons...)
	return menu
}

func CloneBackMarkup(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	menu.Inline(menu.Row(createBtn(menu, langData["BACK_BUTTON"], "help_callback_chelp", "")))
	return menu
}

func HelpBackMarkup(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	menu.Inline(menu.Row(createBtn(menu, langData["BACK_BUTTON"], "settings_back_helper", "")))
	return menu
}

func PrivateHelpPanel(langData map[string]string, botUsername string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	menu.Inline(menu.Row(createBtn(menu, langData["S_B_4"], "", fmt.Sprintf("https://t.me/%s?start=help", botUsername))))
	return menu
}
