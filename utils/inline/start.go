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
			CreateBtn(m, "✙ ᴧᴅᴅ ϻє ᴛσ ʏσᴜʀ ɢʀσυᴘ ✙", "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), sMap[2], 0, false),
			CreateBtn(m, "sυᴘᴘσʀᴛ", "", supportChat, sMap[2], 0, false),
		),
	)
	return m
}

func PrivatePanel(botUsername, supportChat, supportChannel, githubURL string, ownerID int64) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(CreateBtn(m, "ᴧᴅᴅ ϻє ᴛσ ʏσᴜʀ ɢʀσυᴘ", "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), sMap[1], 0, false)),
		m.Row(
			CreateBtn(m, "σᴡηєʀ", "", fmt.Sprintf("tg://user?id=%d", ownerID), sMap[2], 0, false),
			CreateBtn(m, "ᴄʟσηє", "clone_page", "", sMap[2], 0, false),
		),
		m.Row(
			CreateBtn(m, "sυᴘᴘσʀᴛ", "support_page", "", sMap[2], 0, false),
			CreateBtn(m, "sσᴜʀᴄє", "gib_source", "", sMap[2], 0, false),
		),
		m.Row(CreateBtn(m, "ʜєʟᴘ ᴧηᴅ ᴄσϻϻᴧηᴅs", "settings_back_helper", "", sMap[1], 0, false)),
	)
	return m
}

func SupportPanel(supportChat, supportChannel string) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(
			CreateBtn(m, "sυᴘᴘσʀᴛ", "", supportChat, sMap[2], 0, false),
			CreateBtn(m, "υᴘᴅᴧᴛєs", "", supportChannel, sMap[2], 0, false),
		),
		m.Row(CreateBtn(m, "ʙᴧᴄᴋ", "settingsback_helper", "", sMap[1], 0, false)),
	)
	return m
}

func AboutPanel(supportChat, supportChannel, githubURL string, ownerID int64) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(
			CreateBtn(m, "σᴡηєʀ", "", fmt.Sprintf("tg://user?id=%d", ownerID), sMap[2], 0, false),
			CreateBtn(m, "ɢɪᴛʜυʙ", "", githubURL, sMap[2], 0, false),
		),
		m.Row(
			CreateBtn(m, "υᴘᴅᴧᴛєs", "", supportChannel, sMap[2], 0, false),
			CreateBtn(m, "sυᴘᴘσʀᴛ", "", supportChat, sMap[2], 0, false),
		),
		m.Row(CreateBtn(m, "ʙᴧᴄᴋ", "settingsback_helper", "", sMap[1], 0, false)),
	)
	return m
}
