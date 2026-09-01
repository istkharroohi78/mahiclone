package inline

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

// SongMarkup generates the audio/video download selection panel
func SongMarkup(langData map[string]string, vidid string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	menu.Inline(
		menu.Row(
			menu.Data(langData["SG_B_2"], "song_helper", "audio|"+vidid),
			menu.Data(langData["SG_B_3"], "song_helper", "video|"+vidid),
		),
		menu.Row(
			menu.Data(langData["CLOSE_BUTTON"], "close"),
		),
	)
	return menu
}
