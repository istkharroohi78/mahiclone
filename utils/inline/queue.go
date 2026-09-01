package inline

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// QueueMarkupTimer generates the queue timer or simple queue markup
func QueueMarkupTimer(langData map[string]string, duration string, cplay string, videoid string, played string, dur string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	btnQueued := createBtn(menu, langData["QU_B_1"], fmt.Sprintf("GetQueued %s|%s", cplay, videoid), "", 0, "", false)
	btnClose := createBtn(menu, langData["CLOSE_BUTTON"], "close", "", 0, "", false)

	if duration == "Unknown" {
		menu.Inline(menu.Row(btnQueued, btnClose))
	} else {
		// Attempting to handle Python-style .format(played, dur) natively in Go
		timerText := langData["QU_B_2"]
		if strings.Contains(timerText, "{") {
			timerText = strings.ReplaceAll(timerText, "{0}", played)
			timerText = strings.ReplaceAll(timerText, "{1}", dur)
			timerText = strings.ReplaceAll(timerText, "{}", played+" | "+dur) 
		} else {
			timerText = fmt.Sprintf(timerText, played, dur)
		}
		
		btnTimer := createBtn(menu, timerText, "GetTimer", "", 0, "", false)
		menu.Inline(
			menu.Row(btnTimer),
			menu.Row(btnQueued, btnClose),
		)
	}
	return menu
}

// QueueBackMarkup generates the back button for the queue list
func QueueBackMarkup(langData map[string]string, cplay string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}

	btnBack := createBtn(menu, langData["BACK_BUTTON"], fmt.Sprintf("queue_back_timer %s", cplay), "", 0, "", false)
	btnClose := createBtn(menu, langData["CLOSE_BUTTON"], "close", "", 0, "", false)

	menu.Inline(
		menu.Row(btnBack, btnClose),
	)
	return menu
}

// AQMarkup generates the compact admin controls for the active queue
func AQMarkup(langData map[string]string, chatID int64) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	cid := fmt.Sprintf("%d", chatID)

	menu.Inline(
		menu.Row(
			createBtn(menu, "▷", "ADMIN Resume|"+cid, "", 0, "", false),
			createBtn(menu, "II", "ADMIN Pause|"+cid, "", 0, "", false),
			createBtn(menu, "↻", "ADMIN Replay|"+cid, "", 0, "", false),
			createBtn(menu, "‣‣I", "ADMIN Skip|"+cid, "", 0, "", false),
			createBtn(menu, "▢", "ADMIN Stop|"+cid, "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, langData["CLOSE_BUTTON"], "close", "", 0, "", false),
		),
	)
	return menu
}

// QueueAdvancedMarkup handles the expanded queue controls 
// (Renamed slightly to avoid conflict with the queue markup in stream_controls.go)
func QueueAdvancedMarkup(langData map[string]string, vidid string, chatID int64, botUsername string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	cid := fmt.Sprintf("%d", chatID)

	menu.Inline(
		menu.Row(createBtn(menu, langData["S_B_5"], "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), 0, "", false)),
		menu.Row(
			createBtn(menu, "ᴘᴀᴜsᴇ", "ADMIN Pause|"+cid, "", 0, "", false),
			createBtn(menu, "sᴛᴏᴘ", "ADMIN Stop|"+cid, "", 0, "", false),
			createBtn(menu, "sᴋɪᴘ", "ADMIN Skip|"+cid, "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "ʀᴇsᴜᴍᴇ", "ADMIN Resume|"+cid, "", 0, "", false),
			createBtn(menu, "ʀᴇᴘʟᴀʏ", "ADMIN Replay|"+cid, "", 0, "", false),
		),
		menu.Row(createBtn(menu, "ᴍᴏʀᴇ", "", "https://t.me/betabot_hub", 0, "", false)),
	)
	return menu
}
