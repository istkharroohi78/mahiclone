package inline

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

// SettingMarkup generates the main settings panel
func SettingMarkup(langData map[string]string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	menu.Inline(
		menu.Row(
			menu.Data(langData["ST_B_1"], "AU"),
			menu.Data(langData["ST_B_3"], "LG"),
		),
		menu.Row(
			menu.Data(langData["ST_B_2"], "PM"),
		),
		menu.Row(
			menu.Data(langData["ST_B_4"], "VM"),
		),
		menu.Row(
			menu.Data(langData["CLOSE_BUTTON"], "close"),
		),
	)
	return menu
}

// VoteModeMarkup generates the voting mode configuration panel
func VoteModeMarkup(langData map[string]string, current string, mode bool) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	modeText := langData["ST_B_6"]
	if mode {
		modeText = langData["ST_B_5"]
	}

	menu.Inline(
		menu.Row(
			menu.Data("Vᴏᴛɪɴɢ ᴍᴏᴅᴇ ➜", "VOTEANSWER"),
			menu.Data(modeText, "VOMODECHANGE"),
		),
		menu.Row(
			menu.Data("-2", "FERRARIUDTI_M"),
			menu.Data(fmt.Sprintf("ᴄᴜʀʀᴇɴᴛ : %s", current), "ANSWERVOMODE"),
			menu.Data("+2", "FERRARIUDTI_A"),
		),
		menu.Row(
			menu.Data(langData["BACK_BUTTON"], "settings_helper"),
			menu.Data(langData["CLOSE_BUTTON"], "close"),
		),
	)
	return menu
}

// AuthUsersMarkup generates the authorized users configuration panel
func AuthUsersMarkup(langData map[string]string, status bool) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	statusText := langData["ST_B_9"]
	if status {
		statusText = langData["ST_B_8"]
	}

	menu.Inline(
		menu.Row(
			menu.Data(langData["ST_B_7"], "AUTHANSWER"),
			menu.Data(statusText, "AUTH"),
		),
		menu.Row(
			menu.Data(langData["ST_B_1"], "AUTHLIST"),
		),
		menu.Row(
			menu.Data(langData["BACK_BUTTON"], "settings_helper"),
			menu.Data(langData["CLOSE_BUTTON"], "close"),
		),
	)
	return menu
}

// PlaymodeUsersMarkup generates the playmode configuration panel
func PlaymodeUsersMarkup(langData map[string]string, direct bool, group bool, playtype bool) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	directText := langData["ST_B_12"]
	if direct {
		directText = langData["ST_B_11"]
	}

	groupText := langData["ST_B_9"]
	if group {
		groupText = langData["ST_B_8"]
	}

	playtypeText := langData["ST_B_9"]
	if playtype {
		playtypeText = langData["ST_B_8"]
	}

	menu.Inline(
		menu.Row(
			menu.Data(langData["ST_B_10"], "SEARCHANSWER"),
			menu.Data(directText, "MODECHANGE"),
		),
		menu.Row(
			menu.Data(langData["ST_B_13"], "AUTHANSWER"),
			menu.Data(groupText, "CHANNELMODECHANGE"),
		),
		menu.Row(
			menu.Data(langData["ST_B_14"], "PLAYTYPEANSWER"),
			menu.Data(playtypeText, "PLAYTYPECHANGE"),
		),
		menu.Row(
			menu.Data(langData["BACK_BUTTON"], "settings_helper"),
			menu.Data(langData["CLOSE_BUTTON"], "close"),
		),
	)
	return menu
}
