package cplugin

import (
	"fmt"
	"math/rand"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Button styles (Mapped for custom Telebot forks or visual reference)
const (
	Primary = iota
	Success
	Danger
	Transparent
)

// PremiumEmojis list updated with requested IDs
var PremiumEmojis = []string{
	"5258362837411045098",
	"6102938383456146362",
	"5463274047771000031",
	"6100397162976252509",
	"5373310679241466020",
	"5408916593780470262",
	"5776182936638329359",
	"5258389041006518073",
	"6280269890821558384",
	"5936143551854285132",
	"6172332822892647766",
	"5891211339170326418",
	"5409368076447657845",
	"6172312314423808834",
	"6082387600599944892",
	"6271537028307881531",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetStyleMap creates a dynamic color generator for buttons
func GetStyleMap() map[int]int {
	styles := []int{Primary, Success, Danger}
	rand.Shuffle(len(styles), func(i, j int) {
		styles[i], styles[j] = styles[j], styles[i]
	})
	return map[int]int{
		1: styles[0],
		2: styles[1],
		3: styles[2],
		4: styles[0],
		5: styles[1],
	}
}

// CreateBtn builds an inline button, appending premium emojis if supported
func CreateBtn(m *tb.ReplyMarkup, text, data, url string, style int, noEmoji bool) tb.Btn {
	// Note: Custom emoji ID injection would require a modified tb.Btn struct.
	// We handle standard Telebot inline button construction here.
	if url != "" {
		return m.URL(text, url)
	}
	return m.Data(text, data)
}

// AddMeButton generates the standard "Add Me" button
func AddMeButton(m *tb.ReplyMarkup, botUsername string, style int) tb.Btn {
	return CreateBtn(m, "『 𝐀ᴅᴅ 𝐌є 𝐁ᴀʙʏ 』", "", fmt.Sprintf("https://t.me/%s?startgroup=true", botUsername), style, false)
}

// GetBar generates a textual progress bar for tracks
func GetBar(played, dur string) string {
	playedSec := TimeToSeconds(played) // Assuming TimeToSeconds is in the same package (formatters.go)
	durSec := TimeToSeconds(dur)

	totalBlocks := 10
	filledBlocks := 0
	if durSec > 0 {
		filledBlocks = int((float64(playedSec) / float64(durSec)) * float64(totalBlocks))
	}

	bar := ""
	for i := 0; i < filledBlocks; i++ {
		bar += "▰"
	}
	for i := 0; i < totalBlocks-filledBlocks; i++ {
		bar += "▱"
	}
	return fmt.Sprintf("%s %s %s", played, bar, dur)
}

// --- MARKUPS ---

func TrackMarkup(vidID string, userID int64, channel, fPlay string) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(
			CreateBtn(m, "Audio", fmt.Sprintf("MusicStream %s|%d|a|%s|%s", vidID, userID, channel, fPlay), "", sMap[2], false),
			CreateBtn(m, "Video", fmt.Sprintf("MusicStream %s|%d|v|%s|%s", vidID, userID, channel, fPlay), "", sMap[2], false),
		),
		m.Row(CreateBtn(m, "✯ CLOSE ✯", fmt.Sprintf("forceclose %s|%d", vidID, userID), "", sMap[1], false)),
	)
	return m
}

func StreamMarkupTimer(chatID int64, played, dur string) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(CreateBtn(m, GetBar(played, dur), "GetTimer", "", Transparent, true)),
		m.Row(
			CreateBtn(m, "▷", fmt.Sprintf("ADMIN Resume|%d", chatID), "", sMap[3], true),
			CreateBtn(m, "II", fmt.Sprintf("ADMIN Pause|%d", chatID), "", sMap[3], true),
			CreateBtn(m, "‣‣I", fmt.Sprintf("ADMIN Skip|%d", chatID), "", sMap[3], true),
		),
		m.Row(CreateBtn(m, "❖ 𝐀ᴜᴛᴏ𝐏ʟᴀʏ ❖", fmt.Sprintf("ADMIN Autoplay|%d", chatID), "", sMap[1], false)),
		m.Row(CreateBtn(m, "✯ CLOSE ✯", "close", "", sMap[1], false)),
	)
	return m
}

func StreamMarkup(chatID int64) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(
			CreateBtn(m, "▷", fmt.Sprintf("ADMIN Resume|%d", chatID), "", sMap[3], true),
			CreateBtn(m, "II", fmt.Sprintf("ADMIN Pause|%d", chatID), "", sMap[3], true),
			CreateBtn(m, "‣‣I", fmt.Sprintf("ADMIN Skip|%d", chatID), "", sMap[3], true),
		),
		m.Row(CreateBtn(m, "❖ 𝐀ᴜᴛᴏ𝐏ʟᴀʏ ❖", fmt.Sprintf("ADMIN Autoplay|%d", chatID), "", sMap[1], false)),
		m.Row(CreateBtn(m, "✯ CLOSE ✯", "close", "", sMap[1], false)),
	)
	return m
}

func QueueMarkup(vidID string, chatID int64, botUsername string) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(AddMeButton(m, botUsername, sMap[1])),
		m.Row(
			CreateBtn(m, "II ᴘᴀᴜsᴇ", fmt.Sprintf("ADMIN Pause|%d", chatID), "", sMap[3], true),
			CreateBtn(m, "▢ sᴛᴏᴘ", fmt.Sprintf("ADMIN Stop|%d", chatID), "", sMap[3], true),
			CreateBtn(m, "sᴋɪᴘ ‣‣I", fmt.Sprintf("ADMIN Skip|%d", chatID), "", sMap[3], true),
		),
		m.Row(
			CreateBtn(m, "▷ ʀᴇsᴜᴍᴇ", fmt.Sprintf("ADMIN Resume|%d", chatID), "", sMap[2], true),
			CreateBtn(m, "ʀᴇᴘʟᴀʏ ↺", fmt.Sprintf("ADMIN Replay|%d", chatID), "", sMap[2], true),
		),
		m.Row(CreateBtn(m, "❖ 𝐀ᴜᴛᴏ𝐏ʟᴀʏ ❖", fmt.Sprintf("ADMIN Autoplay|%d", chatID), "", sMap[1], false)),
		m.Row(CreateBtn(m, "ᴍᴏʀᴇ", fmt.Sprintf("PanelMarkup None|%d", chatID), "", sMap[1], false)),
	)
	return m
}

func PanelMarkup1(vidID string, chatID int64, botUsername string) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(AddMeButton(m, botUsername, sMap[1])),
		m.Row(
			CreateBtn(m, "sᴜғғʟᴇ", fmt.Sprintf("ADMIN Shuffle|%d", chatID), "", sMap[2], true),
			CreateBtn(m, "ʟᴏᴏᴘ ↺", fmt.Sprintf("ADMIN Loop|%d", chatID), "", sMap[2], true),
		),
		m.Row(
			CreateBtn(m, "◁ 10 sᴇᴄ", fmt.Sprintf("ADMIN 1|%d", chatID), "", sMap[2], true),
			CreateBtn(m, "10 sᴇᴄ ▷", fmt.Sprintf("ADMIN 2|%d", chatID), "", sMap[2], true),
		),
		m.Row(CreateBtn(m, "❖ 𝐀ᴜᴛᴏ𝐏ʟᴀʏ ❖", fmt.Sprintf("ADMIN Autoplay|%d", chatID), "", sMap[1], false)),
		m.Row(
			CreateBtn(m, "ʜᴏᴍᴇ", fmt.Sprintf("Pages Back|2|%s|%d", vidID, chatID), "", sMap[2], true),
			CreateBtn(m, "ɴᴇxᴛ", fmt.Sprintf("Pages Forw|2|%s|%d", vidID, chatID), "", sMap[2], true),
		),
	)
	return m
}

func PanelMarkupClone(vidID string, chatID int64) *tb.ReplyMarkup {
	m := &tb.ReplyMarkup{}
	sMap := GetStyleMap()

	m.Inline(
		m.Row(
			CreateBtn(m, "▷", fmt.Sprintf("ADMIN Resume|%d", chatID), "", sMap[3], true),
			CreateBtn(m, "II", fmt.Sprintf("ADMIN Pause|%d", chatID), "", sMap[3], true),
			CreateBtn(m, "‣‣I", fmt.Sprintf("ADMIN Skip|%d", chatID), "", sMap[3], true),
		),
		m.Row(
			CreateBtn(m, "<- 20s", fmt.Sprintf("ADMIN SeekBack|%d", chatID), "", sMap[4], true),
			CreateBtn(m, "🔁", fmt.Sprintf("ADMIN Loop|%d", chatID), "", sMap[4], true),
			CreateBtn(m, "🔀", fmt.Sprintf("ADMIN Shuffle|%d", chatID), "", sMap[4], true),
			CreateBtn(m, "20s + ->", fmt.Sprintf("ADMIN SeekForward|%d", chatID), "", sMap[4], true),
		),
		m.Row(CreateBtn(m, "✯ CLOSE ✯", "close", "", sMap[1], true)),
	)
	return m
}
