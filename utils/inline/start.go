package inline

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func StartPanel(botUsername, supportChat string) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(
			createBtn(m, "✙ ᴧᴅᴅ ϻє ᴛσ ʏσᴜʀ ɢʀσυᴘ ✙", "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), sMap[2], "", false),
			createBtn(m, "sυᴘᴘσʀᴛ", "", supportChat, sMap[2], "", false),
		),
	)
	return m
}

func PrivatePanel(botUsername, supportChat, supportChannel, githubURL string, ownerID int64) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(createBtn(m, "ᴧᴅᴅ ϻє ᴛσ ʏσᴜʀ ɢʀσυᴘ", "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), sMap[1], "", false)),
		m.Row(
			createBtn(m, "σᴡηєʀ", "", fmt.Sprintf("tg://user?id=%d", ownerID), sMap[2], "", false),
			createBtn(m, "ᴄʟσηє", "clone_page", "", sMap[2], "", false),
		),
		m.Row(
			createBtn(m, "sυᴘᴘσʀᴛ", "support_page", "", sMap[2], "", false),
			createBtn(m, "sσᴜʀᴄє", "gib_source", "", sMap[2], "", false),
		),
		m.Row(createBtn(m, "ʜєʟᴘ ᴧηᴅ ᴄσϻϻᴧηᴅs", "settings_back_helper", "", sMap[1], "", false)),
	)
	return m
}

func SupportPanel(supportChat, supportChannel string) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(
			createBtn(m, "sυᴘᴘσʀᴛ", "", supportChat, sMap[2], "", false),
			createBtn(m, "υᴘᴅᴧᴛєs", "", supportChannel, sMap[2], "", false),
		),
		m.Row(createBtn(m, "ʙᴧᴄᴋ", "settingsback_helper", "", sMap[1], "", false)),
	)
	return m
}

func AboutPanel(supportChat, supportChannel, githubURL string, ownerID int64) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(
			createBtn(m, "σᴡηєʀ", "", fmt.Sprintf("tg://user?id=%d", ownerID), sMap[2], "", false),
			createBtn(m, "ɢɪᴛʜυʙ", "", githubURL, sMap[2], "", false),
		),
		m.Row(
			createBtn(m, "υᴘᴅᴧᴛєs", "", supportChannel, sMap[2], "", false),
			createBtn(m, "sυᴘᴘσʀᴛ", "", supportChat, sMap[2], "", false),
		),
		m.Row(createBtn(m, "ʙᴧᴄᴋ", "settingsback_helper", "", sMap[1], "", false)),
	)
	return m
}
