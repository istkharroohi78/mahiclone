package inline

import (
	"fmt"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// timeToSeconds converts a time string like "03:15" to 195 seconds
func timeToSeconds(t string) int {
	parts := strings.Split(t, ":")
	sec := 0
	mult := 1
	for i := len(parts) - 1; i >= 0; i-- {
		val, _ := strconv.Atoi(parts[i])
		sec += val * mult
		mult *= 60
	}
	return sec
}

// createBtn generates a Telebot inline button. 
// It splits Pyrogram-style callbacks (e.g., "ADMIN Pause|123") into Telebot's unique ID and payload.
func cloneButton(menu *tb.ReplyMarkup) tb.Btn {
	return createBtn(menu, "ᴄʟᴏɴᴇ-ᴍᴇ", "", "https://t.me/clone_MUSICrobot", 0, "", false)
}

// TrackMarkup
func TrackMarkup(langData map[string]string, videoid string, userID int64, channel string, fplay string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	uid := fmt.Sprintf("%d", userID)

	menu.Inline(
		menu.Row(
			createBtn(menu, langData["P_B_1"], fmt.Sprintf("MusicStream %s|%s|a|%s|%s", videoid, uid, channel, fplay), "", 0, "", false),
			createBtn(menu, langData["P_B_2"], fmt.Sprintf("MusicStream %s|%s|v|%s|%s", videoid, uid, channel, fplay), "", 0, "", false),
		),
		menu.Row(
			cloneButton(menu),
			createBtn(menu, langData["CLOSE_BUTTON"], fmt.Sprintf("forceclose %s|%s", videoid, uid), "", 0, "", false),
		),
	)
	return menu
}

// StreamMarkupTimer
func StreamMarkupTimer(langData map[string]string, chatID int64, played string, dur string) *tb.ReplyMarkup {
	playedSec := timeToSeconds(played)
	durationSec := 0
	
	if strings.ToLower(dur) != "live" && strings.ToLower(dur) != "unknown" && dur != "0" {
		durationSec = timeToSeconds(dur)
	}

	totalBlocks := 10
	filledBlocks := 0
	if durationSec > 0 {
		filledBlocks = int((float64(playedSec) / float64(durationSec)) * float64(totalBlocks))
	}
	
	if filledBlocks < 0 {
		filledBlocks = 0
	} else if filledBlocks >= totalBlocks {
		filledBlocks = totalBlocks - 1
	}

	bar := strings.Repeat("▰", filledBlocks) + "🎵" + strings.Repeat("▱", totalBlocks-filledBlocks-1)
	timerText := fmt.Sprintf("%s %s %s", played, bar, dur)
	cid := fmt.Sprintf("%d", chatID)

	menu := &tb.ReplyMarkup{}
	menu.Inline(
		menu.Row(createBtn(menu, timerText, "GetTimer", "", 0, "", false)),
		menu.Row(
			createBtn(menu, PlayEmoji, "ADMIN Resume|"+cid, "", 0, "", false),
			createBtn(menu, PauseEmoji, "ADMIN Pause|"+cid, "", 0, "", false),
			createBtn(menu, ReplayEmoji, "ADMIN Replay|"+cid, "", 0, "", false),
			createBtn(menu, SkipEmoji, "ADMIN Skip|"+cid, "", 0, "", false),
			createBtn(menu, StopEmoji, "ADMIN Stop|"+cid, "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "ᴀᴜᴛᴏ-ᴘʟᴀʏ", "ADMIN Autoplay|"+cid, "", 0, "", false),
			cloneButton(menu),
		),
		menu.Row(createBtn(menu, langData["CLOSE_BUTTON"], "close", "", 0, "", false)),
	)
	return menu
}

// StreamMarkup
func StreamMarkup(langData map[string]string, chatID int64) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	cid := fmt.Sprintf("%d", chatID)

	menu.Inline(
		menu.Row(
			createBtn(menu, PlayEmoji, "ADMIN Resume|"+cid, "", 0, "", false),
			createBtn(menu, PauseEmoji, "ADMIN Pause|"+cid, "", 0, "", false),
			createBtn(menu, ReplayEmoji, "ADMIN Replay|"+cid, "", 0, "", false),
			createBtn(menu, SkipEmoji, "ADMIN Skip|"+cid, "", 0, "", false),
			createBtn(menu, StopEmoji, "ADMIN Stop|"+cid, "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "ᴀᴜᴛᴏ-ᴘʟᴀʏ", "ADMIN Autoplay|"+cid, "", 0, "", false),
			cloneButton(menu),
		),
		menu.Row(createBtn(menu, langData["CLOSE_BUTTON"], "close", "", 0, "", false)),
	)
	return menu
}

// PlaylistMarkup
func PlaylistMarkup(langData map[string]string, videoid string, userID int64, ptype string, channel string, fplay string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	uid := fmt.Sprintf("%d", userID)

	menu.Inline(
		menu.Row(
			createBtn(menu, langData["P_B_1"], fmt.Sprintf("LuckyPlaylists %s|%s|%s|a|%s|%s", videoid, uid, ptype, channel, fplay), "", 0, "", false),
			createBtn(menu, langData["P_B_2"], fmt.Sprintf("LuckyPlaylists %s|%s|%s|v|%s|%s", videoid, uid, ptype, channel, fplay), "", 0, "", false),
		),
		menu.Row(
			cloneButton(menu),
			createBtn(menu, langData["CLOSE_BUTTON"], fmt.Sprintf("forceclose %s|%s", videoid, uid), "", 0, "", false),
		),
	)
	return menu
}

// LivestreamMarkup
func LivestreamMarkup(langData map[string]string, videoid string, userID int64, mode string, channel string, fplay string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	uid := fmt.Sprintf("%d", userID)

	menu.Inline(
		menu.Row(createBtn(menu, langData["P_B_3"], fmt.Sprintf("LiveStream %s|%s|%s|%s|%s", videoid, uid, mode, channel, fplay), "", 0, "", false)),
		menu.Row(
			cloneButton(menu),
			createBtn(menu, langData["CLOSE_BUTTON"], fmt.Sprintf("forceclose %s|%s", videoid, uid), "", 0, "", false),
		),
	)
	return menu
}

// SliderMarkup
func SliderMarkup(langData map[string]string, videoid string, userID int64, query string, queryType string, channel string, fplay string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	uid := fmt.Sprintf("%d", userID)
	
	if len(query) > 20 {
		query = query[:20]
	}

	menu.Inline(
		menu.Row(
			createBtn(menu, langData["P_B_1"], fmt.Sprintf("MusicStream %s|%s|a|%s|%s", videoid, uid, channel, fplay), "", 0, "", false),
			createBtn(menu, langData["P_B_2"], fmt.Sprintf("MusicStream %s|%s|v|%s|%s", videoid, uid, channel, fplay), "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "ʙᴀᴄᴋ", fmt.Sprintf("slider B|%s|%s|%s|%s|%s", queryType, query, uid, channel, fplay), "", 0, "", false),
			createBtn(menu, langData["CLOSE_BUTTON"], fmt.Sprintf("forceclose %s|%s", query, uid), "", 0, "", false),
			createBtn(menu, "ɴᴇxᴛ", fmt.Sprintf("slider F|%s|%s|%s|%s|%s", queryType, query, uid, channel, fplay), "", 0, "", false),
		),
		menu.Row(cloneButton(menu)),
	)
	return menu
}

// TelegramMarkup
func TelegramMarkup(langData map[string]string, chatID int64) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	cid := fmt.Sprintf("%d", chatID)

	menu.Inline(
		menu.Row(
			createBtn(menu, "ɴᴇxᴛ", "PanelMarkup None|"+cid, "", 0, "", false),
			createBtn(menu, langData["CLOSEMENU_BUTTON"], "close", "", 0, "", false),
		),
	)
	return menu
}

// QueueMarkup
func QueueMarkup(langData map[string]string, videoid string, chatID int64, botUsername string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	cid := fmt.Sprintf("%d", chatID)

	menu.Inline(
		menu.Row(createBtn(menu, langData["S_B_3"], "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), 0, "", false)),
		menu.Row(
			createBtn(menu, PlayEmoji, "ADMIN Resume|"+cid, "", 0, "", false),
			createBtn(menu, PauseEmoji, "ADMIN Pause|"+cid, "", 0, "", false),
			createBtn(menu, ReplayEmoji, "ADMIN Replay|"+cid, "", 0, "", false),
			createBtn(menu, SkipEmoji, "ADMIN Skip|"+cid, "", 0, "", false),
			createBtn(menu, StopEmoji, "ADMIN Stop|"+cid, "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "ᴀᴜᴛᴏ-ᴘʟᴀʏ", "ADMIN Autoplay|"+cid, "", 0, "", false),
			cloneButton(menu),
		),
		menu.Row(createBtn(menu, "ᴍᴏʀᴇ", "PanelMarkup None|"+cid, "", 0, "", false)),
	)
	return menu
}

// PanelMarkup1
func PanelMarkup1(langData map[string]string, videoid string, chatID int64, botUsername string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	cid := fmt.Sprintf("%d", chatID)

	menu.Inline(
		menu.Row(createBtn(menu, langData["S_B_3"], "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), 0, "", false)),
		menu.Row(
			createBtn(menu, "sʜᴜғғʟᴇ", "ADMIN Shuffle|"+cid, "", 0, "", false),
			createBtn(menu, "ʟᴏᴏᴘ", "ADMIN Loop|"+cid, "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "-10 sᴇᴄ", "ADMIN 1|"+cid, "", 0, "", false),
			createBtn(menu, "+10 sᴇᴄ", "ADMIN 2|"+cid, "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "ᴀᴜᴛᴏ-ᴘʟᴀʏ", "ADMIN Autoplay|"+cid, "", 0, "", false),
			cloneButton(menu),
		),
		menu.Row(
			createBtn(menu, "ʜᴏᴍᴇ", fmt.Sprintf("Pages Back|2|%s|%s", videoid, cid), "", 0, "", false),
			createBtn(menu, "ɴᴇxᴛ", fmt.Sprintf("Pages Forw|2|%s|%s", videoid, cid), "", 0, "", false),
		),
	)
	return menu
}

// PanelMarkup2
func PanelMarkup2(langData map[string]string, videoid string, chatID int64, botUsername string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	cid := fmt.Sprintf("%d", chatID)

	menu.Inline(
		menu.Row(createBtn(menu, langData["S_B_3"], "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), 0, "", false)),
		menu.Row(
			createBtn(menu, "0.5x", "SpeedUP "+cid+"|0.5", "", 0, "", false),
			createBtn(menu, "0.75x", "SpeedUP "+cid+"|0.75", "", 0, "", false),
			createBtn(menu, "1.0x", "SpeedUP "+cid+"|1.0", "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "1.5x", "SpeedUP "+cid+"|1.5", "", 0, "", false),
			createBtn(menu, "2.0x", "SpeedUP "+cid+"|2.0", "", 0, "", false),
		),
		menu.Row(
			cloneButton(menu),
			createBtn(menu, "ʙᴀᴄᴋ", fmt.Sprintf("Pages Back|1|%s|%s", videoid, cid), "", 0, "", false),
		),
	)
	return menu
}

// PanelMarkup3
func PanelMarkup3(langData map[string]string, videoid string, chatID int64) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	cid := fmt.Sprintf("%d", chatID)

	menu.Inline(
		menu.Row(
			createBtn(menu, "0.5x", "SpeedUP "+cid+"|0.5", "", 0, "", false),
			createBtn(menu, "0.75x", "SpeedUP "+cid+"|0.75", "", 0, "", false),
			createBtn(menu, "1.0x", "SpeedUP "+cid+"|1.0", "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "1.5x", "SpeedUP "+cid+"|1.5", "", 0, "", false),
			createBtn(menu, "2.0x", "SpeedUP "+cid+"|2.0", "", 0, "", false),
		),
		menu.Row(
			cloneButton(menu),
			createBtn(menu, "ʙᴀᴄᴋ", fmt.Sprintf("Pages Back|2|%s|%s", videoid, cid), "", 0, "", false),
		),
	)
	return menu
}

// PanelMarkup4
func PanelMarkup4(langData map[string]string, videoid string, chatID int64, played string, dur string) *tb.ReplyMarkup {
	playedSec := timeToSeconds(played)
	durationSec := 0
	
	if strings.ToLower(dur) != "live" && strings.ToLower(dur) != "unknown" && dur != "0" {
		durationSec = timeToSeconds(dur)
	}

	totalBlocks := 10
	filledBlocks := 0
	if durationSec > 0 {
		filledBlocks = int((float64(playedSec) / float64(durationSec)) * float64(totalBlocks))
	}
	
	if filledBlocks < 0 {
		filledBlocks = 0
	} else if filledBlocks >= totalBlocks {
		filledBlocks = totalBlocks - 1
	}

	bar := strings.Repeat("▰", filledBlocks) + "🎵" + strings.Repeat("▱", totalBlocks-filledBlocks-1)
	timerText := fmt.Sprintf("%s %s %s", played, bar, dur)
	cid := fmt.Sprintf("%d", chatID)

	menu := &tb.ReplyMarkup{}
	menu.Inline(
		menu.Row(createBtn(menu, timerText, "GetTimer", "", 0, "", false)),
		menu.Row(
			createBtn(menu, PlayEmoji, "ADMIN Resume|"+cid, "", 0, "", false),
			createBtn(menu, PauseEmoji, "ADMIN Pause|"+cid, "", 0, "", false),
			createBtn(menu, ReplayEmoji, "ADMIN Replay|"+cid, "", 0, "", false),
			createBtn(menu, SkipEmoji, "ADMIN Skip|"+cid, "", 0, "", false),
			createBtn(menu, StopEmoji, "ADMIN Stop|"+cid, "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "ᴀᴜᴛᴏ-ᴘʟᴀʏ", "ADMIN Autoplay|"+cid, "", 0, "", false),
			cloneButton(menu),
		),
		menu.Row(createBtn(menu, langData["CLOSE_BUTTON"], "close", "", 0, "", false)),
		menu.Row(createBtn(menu, "ʜᴏᴍᴇ", fmt.Sprintf("MainMarkup %s|%d", videoid, chatID), "", 0, "", false)),
	)
	return menu
}

// PanelMarkup5
func PanelMarkup5(langData map[string]string, videoid string, chatID int64, botUsername string) *tb.ReplyMarkup {
	menu := &tb.ReplyMarkup{}
	cid := fmt.Sprintf("%d", chatID)

	menu.Inline(
		menu.Row(createBtn(menu, langData["S_B_3"], "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), 0, "", false)),
		menu.Row(
			createBtn(menu, PlayEmoji, "ADMIN Resume|"+cid, "", 0, "", false),
			createBtn(menu, PauseEmoji, "ADMIN Pause|"+cid, "", 0, "", false),
			createBtn(menu, ReplayEmoji, "ADMIN Replay|"+cid, "", 0, "", false),
			createBtn(menu, SkipEmoji, "ADMIN Skip|"+cid, "", 0, "", false),
			createBtn(menu, StopEmoji, "ADMIN Stop|"+cid, "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "ᴀᴜᴛᴏ-ᴘʟᴀʏ", "ADMIN Autoplay|"+cid, "", 0, "", false),
			cloneButton(menu),
		),
		menu.Row(
			createBtn(menu, "ʜᴏᴍᴇ", fmt.Sprintf("MainMarkup %s|%s", videoid, cid), "", 0, "", false),
			createBtn(menu, "ɴᴇxᴛ", fmt.Sprintf("Pages Forw|1|%s|%s", videoid, cid), "", 0, "", false),
		),
	)
	return menu
}

// PanelMarkupClone
func PanelMarkupClone(langData map[string]string, videoid string, chatID int64, played string, dur string) *tb.ReplyMarkup {
	playedSec := timeToSeconds(played)
	durationSec := 0
	
	if strings.ToLower(dur) != "live" && strings.ToLower(dur) != "unknown" && dur != "0" {
		durationSec = timeToSeconds(dur)
	}

	totalBlocks := 10
	filledBlocks := 0
	if durationSec > 0 {
		filledBlocks = int((float64(playedSec) / float64(durationSec)) * float64(totalBlocks))
	}
	
	if filledBlocks < 0 {
		filledBlocks = 0
	} else if filledBlocks >= totalBlocks {
		filledBlocks = totalBlocks - 1
	}

	bar := strings.Repeat("▰", filledBlocks) + "🎵" + strings.Repeat("▱", totalBlocks-filledBlocks-1)
	timerText := fmt.Sprintf("%s %s %s", played, bar, dur)
	cid := fmt.Sprintf("%d", chatID)

	menu := &tb.ReplyMarkup{}
	menu.Inline(
		menu.Row(createBtn(menu, timerText, "GetTimer", "", 0, "", false)),
		menu.Row(
			createBtn(menu, PlayEmoji, "ADMIN Resume|"+cid, "", 0, "", false),
			createBtn(menu, PauseEmoji, "ADMIN Pause|"+cid, "", 0, "", false),
			createBtn(menu, ReplayEmoji, "ADMIN Replay|"+cid, "", 0, "", false),
			createBtn(menu, SkipEmoji, "ADMIN Skip|"+cid, "", 0, "", false),
			createBtn(menu, StopEmoji, "ADMIN Stop|"+cid, "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "-20s", "ADMIN SeekBack|"+cid, "", 0, "", false),
			createBtn(menu, "+20s", "ADMIN SeekForward|"+cid, "", 0, "", false),
		),
		menu.Row(
			createBtn(menu, "ᴀᴜᴛᴏ-ᴘʟᴀʏ", "ADMIN Autoplay|"+cid, "", 0, "", false),
			cloneButton(menu),
		),
		menu.Row(createBtn(menu, langData["CLOSE_BUTTON"], "close", "", 0, "", false)),
	)
	return menu
}
