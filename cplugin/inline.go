package cplugin

import (
	"fmt"
	"strconv"
	"strings"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// InlineQueryHandler handles inline YT searches
func InlineQueryHandler(b *tb.Bot, q *tb.Query) {
	text := strings.TrimSpace(strings.ToLower(q.Text))

	if text == "" {
		// Return default cache/answer
		b.Answer(q, &tb.QueryResponse{
			Results:   tb.Results{},
			CacheTime: 10,
		})
		return
	}

	// Assuming utils.SearchYouTube interfaces with ytsearch API/library
	results := utils.SearchYouTube(text, 15)
	if len(results) == 0 {
		return
	}

	var tbResults tb.Results
	for i, res := range results {
		desc := fmt.Sprintf("%s | %s | %s | %s", res.Views, res.Duration, res.Channel, res.Published)
		caption := fmt.Sprintf(`❄ <b>ᴛɪᴛʟᴇ :</b> <a href="%s">%s</a>

⏳ <b>ᴅᴜʀᴀᴛɪᴏɴ :</b> %s ᴍɪɴᴜᴛᴇs
👀 <b>ᴠɪᴇᴡs :</b> <code>%s</code>
🎥 <b>ᴄʜᴀɴɴᴇʟ :</b> <a href="%s">%s</a>
⏰ <b>ᴘᴜʙʟɪsʜᴇᴅ ᴏɴ :</b> %s

<u><b>➻ ɪɴʟɪɴᴇ sᴇᴀʀᴄʜ ᴍᴏᴅᴇ ʙʏ %s</b></u>`, res.Link, res.Title, res.Duration, res.Views, res.ChannelLink, res.Channel, res.Published, b.Me.FirstName)

		// Inline button
		btn := tb.InlineButton{Text: "ʏᴏᴜᴛᴜʙᴇ 🎄", URL: res.Link}
		markup := [][]tb.InlineButton{{btn}}

		result := &tb.PhotoResult{
			ResultBase: tb.ResultBase{
				ID:          strconv.Itoa(i),
				Title:       res.Title,
				Description: desc,
				ReplyMarkup: markup,
			},
			URL:       res.Thumbnail,
			ThumbURL:  res.Thumbnail,
			Caption:   caption,
			ParseMode: tb.ModeHTML,
		}
		tbResults = append(tbResults, result)
	}

	b.Answer(q, &tb.QueryResponse{
		Results:   tbResults,
		CacheTime: 300,
	})
}
