package inline

import (
	"math/rand"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Button Styles (Assuming custom implementation constants)
const (
	Primary = 1
	Success = 2
	Danger  = 3
)

var PremiumEmojis = []int64{
	5258362837411045098, 6102938383456146362, 5463274047771000031, 6100397162976252509,
	5373310679241466020, 5408916593780470262, 5776182936638329359, 5258389041006518073,
	6280269890821558384, 5936143551854285132, 6172332822892647766, 5891211339170326418,
	5409368076447657845, 6172312314423808834, 6082387600599944892, 6271537028307881531,
}

const (
	PlayEmoji   int64 = 6158973722255429425 // ▶️
	PauseEmoji  int64 = 4970176665062736422 // ⏸️
	ReplayEmoji int64 = 5258419835922030550 // 🔁
	SkipEmoji   int64 = 4969851488793788974 // ⏭️
	StopEmoji   int64 = 6129486856212979482 // 🛑
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetStyleMap returns a randomized map of button styles
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
	}
}

// CreateBtn builds a smart inline button
func CreateBtn(markup *tb.ReplyMarkup, text, cb, url string, style int, emojiID int64, noEmoji bool) tb.Btn {
	var btn tb.Btn
	if url != "" {
		btn = markup.URL(text, url)
	} else {
		btn = markup.Data(text, cb)
	}

	// Simulated injection of custom fields (Requires a patched telebot version to serialize)
	// btn.Style = style
	if !noEmoji {
		if emojiID != 0 {
			// btn.IconCustomEmojiID = emojiID
		} else {
			// btn.IconCustomEmojiID = PremiumEmojis[rand.Intn(len(PremiumEmojis))]
		}
	}
	return btn
}

// CloneButton returns the standard clone-me button
func CloneButton(markup *tb.ReplyMarkup, style int) tb.Btn {
	return CreateBtn(markup, "ᴄʟᴏɴᴇ-ᴍᴇ", "", "https://t.me/clone_MUSICrobot", style, 0, false)
}
