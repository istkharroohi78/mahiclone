package inline

import (
	"fmt"
	"math/rand"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// getRandomEmoji returns a random custom emoji ID
func getRandomEmoji() string {
	rand.Seed(time.Now().UnixNano())
	return PremiumEmojis[rand.Intn(len(PremiumEmojis))]
}

func HelpPanel(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	first := []tb.Row{
		menu.Row(
			createBtn(menu, langData["H_B_1"], "help_callback_hb1", "", Primary, "", true),
			createBtn(menu, langData["H_B_2"], "help_callback_hb2", "", Primary, "", true),
			createBtn(menu, langData["H_B_3"], "help_callback_hb3", "", Primary, "", true),
		),
		menu.Row(
			createBtn(menu, langData["H_B_4"], "help_callback_hb4", "", Primary, "", true),
			createBtn(menu, langData["H_B_5"], "help_callback_hb5", "", Primary, "", true),
			createBtn(menu, langData["H_B_6"], "help_callback_hb6", "", Primary, "", true),
		),
		menu.Row(
			createBtn(menu, langData["H_B_7"], "help_callback_hb7", "", Primary, "", true),
			createBtn(menu, langData["H_B_8"], "help_callback_hb8", "", Primary, "", true),
			createBtn(menu, langData["H_B_9"], "help_callback_hb9", "", Primary, "", true),
		),
		menu.Row(
			createBtn(menu, langData["H_B_10"], "help_callback_hb10", "", Primary, "", true),
			createBtn(menu, langData["H_B_11"], "help_callback_hb11", "", Primary, "", true),
			createBtn(menu, langData["H_B_12"], "help_callback_hb12", "", Primary, "", true),
		),
		menu.Row(
			createBtn(menu, langData["H_B_13"], "help_callback_hb13", "", Primary, "", true),
			createBtn(menu, langData["H_B_14"], "help_callback_hb14", "", Primary, "", true),
			createBtn(menu, langData["H_B_15"], "help_callback_hb15", "", Primary, "", true),
		),
		menu.Row(
			createBtn(menu, langData["BACK_BUTTON"], "settingsback_helper", "", Danger, "", true),
		),
	}

	menu.Inline(first...)
	return menu
}

func FirstPage(langData map[string]string, isOwner bool) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	first := []tb.Row{
		menu.Row(
			createBtn(menu, langData["H_B_1"], "help_callback_hb1", "", Primary, "", true),
			createBtn(menu, langData["H_B_2"], "help_callback_hb2", "", Primary, "", true),
			createBtn(menu, langData["H_B_3"], "help_callback_hb3", "", Primary, "", true),
		),
		menu.Row(
			createBtn(menu, langData["H_B_11"], "help_callback_hb11", "", Primary, "", true),
			createBtn(menu, langData["H_B_8"], "help_callback_hb8", "", Primary, "", true),
			createBtn(menu, langData["H_B_6"], "help_callback_hb6", "", Primary, "", true),
		),
		menu.Row(
			createBtn(menu, langData["H_B_13"], "help_callback_hb13", "", Primary, "", true),
			createBtn(menu, langData["H_B_12"], "help_callback_hb12", "", Primary, "", true),
			createBtn(menu, langData["H_B_9"], "help_callback_cloghelp", "", Primary, "", true),
		),
		menu.Row(
			createBtn(menu, langData["H_B_10"], "help_callback_hb10", "", Primary, "", true),
			createBtn(menu, langData["H_B_14"], "help_callback_hb14", "", Primary, "", true),
			createBtn(menu, langData["H_B_15"], "help_callback_hb15", "", Primary, "", true),
		),
	}

	if isOwner {
		first = append(first, menu.Row(createBtn(menu, "🛠 ᴄʟᴏɴᴇ ғᴇᴀᴛᴜʀᴇ", "help_callback_chelp", "", Success, "", true)))
	}

	first = append(first, menu.Row(createBtn(menu, langData["BACK_BUTTON"], "settingsback_home", "", Danger, "", true)))
	
	menu.Inline(first...)
	return menu
}

func CloneHelpPanel(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	buttons := []tb.Row{
		menu.Row(createBtn(menu, "ᴍᴀɴᴀɢᴇ", "help_callback_clone_manage", "", Primary, "", true)),
		menu.Row(
			createBtn(menu, "sᴛᴀʀᴛ", "help_callback_clone_start", "", Primary, "", true),
			createBtn(menu, "ᴘɪɴɢ", "help_callback_clone_ping", "", Primary, "", true),
		),
		menu.Row(
			createBtn(menu, "ᴘʟᴀʏ ᴍᴏᴅᴇ", "help_callback_clone_play", "", Primary, "", true),
			createBtn(menu, "ʟᴏɢɢᴇʀ", "help_callback_clone_logger", "", Primary, "", true),
		),
		menu.Row(createBtn(menu, "ʙᴜᴛᴛᴏɴs & ʀᴇɴᴀᴍᴇ", "help_callback_clone_buttons", "", Primary, "", true)),
		menu.Row(createBtn(menu, langData["BACK_BUTTON"], "settings_back_helper", "", Danger, "", true)),
	}
	
	menu.Inline(buttons...)
	return menu
}

func CloneBackMarkup(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	menu.Inline(menu.Row(createBtn(menu, langData["BACK_BUTTON"], "help_callback_chelp", "", Danger, "", true)))
	return menu
}

func HelpBackMarkup(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	menu.Inline(menu.Row(createBtn(menu, langData["BACK_BUTTON"], "settings_back_helper", "", Danger, "", true)))
	return menu
}

func PrivateHelpPanel(langData map[string]string, botUsername string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	menu.Inline(menu.Row(createBtn(menu, langData["S_B_4"], "", fmt.Sprintf("https://t.me/%s?start=help", botUsername), Primary, "", true)))
	return menu
}
