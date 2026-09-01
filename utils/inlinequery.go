package utils

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

// GetStaticInlineAnswers returns the default static inline query results
func GetStaticInlineAnswers() tb.Results {
	return tb.Results{
		&tb.ArticleResult{
			ResultBase:  tb.ResultBase{ID: "pause"},
			Title:       "Pᴀᴜsᴇ",
			Description: "ᴩᴀᴜsᴇ ᴛʜᴇ ᴄᴜʀʀᴇɴᴛ ᴩʟᴀʏɪɴɢ sᴛʀᴇᴀᴍ ᴏɴ ᴠɪᴅᴇᴏᴄʜᴀᴛ.",
			ThumbURL:    "https://telegra.ph/file/c5952790fa8235f499749.jpg",
			Text:        "/pause",
		},
		&tb.ArticleResult{
			ResultBase:  tb.ResultBase{ID: "resume"},
			Title:       "Rᴇsᴜᴍᴇ",
			Description: "ʀᴇsᴜᴍᴇ ᴛʜᴇ ᴩᴀᴜsᴇᴅ sᴛʀᴇᴀᴍ ᴏɴ ᴠɪᴅᴇᴏᴄʜᴀᴛ.",
			ThumbURL:    "https://telegra.ph/file/c5952790fa8235f499749.jpg",
			Text:        "/resume",
		},
		&tb.ArticleResult{
			ResultBase:  tb.ResultBase{ID: "skip"},
			Title:       "Sᴋɪᴩ",
			Description: "sᴋɪᴩ ᴛʜᴇ ᴄᴜʀʀᴇɴᴛ ᴩʟᴀʏɪɴɢ sᴛʀᴇᴀᴍ ᴏɴ ᴠɪᴅᴇᴏᴄʜᴀᴛ.",
			ThumbURL:    "https://telegra.ph/file/c5952790fa8235f499749.jpg",
			Text:        "/skip",
		},
		&tb.ArticleResult{
			ResultBase:  tb.ResultBase{ID: "end"},
			Title:       "Eɴᴅ",
			Description: "ᴇɴᴅ ᴛʜᴇ ᴄᴜʀʀᴇɴᴛ ᴩʟᴀʏɪɴɢ sᴛʀᴇᴀᴍ ᴏɴ ᴠɪᴅᴇᴏᴄʜᴀᴛ.",
			ThumbURL:    "https://telegra.ph/file/c5952790fa8235f499749.jpg",
			Text:        "/end",
		},
		&tb.ArticleResult{
			ResultBase:  tb.ResultBase{ID: "shuffle"},
			Title:       "Sʜᴜғғʟᴇ",
			Description: "sʜᴜғғʟᴇ ᴛʜᴇ ǫᴜᴇᴜᴇᴅ sᴏɴɢs ɪɴ ᴩʟᴀʏʟɪsᴛ.",
			ThumbURL:    "https://telegra.ph/file/c5952790fa8235f499749.jpg",
			Text:        "/shuffle",
		},
		&tb.ArticleResult{
			ResultBase:  tb.ResultBase{ID: "loop"},
			Title:       "Lᴏᴏᴩ",
			Description: "ʟᴏᴏᴩ ᴛʜᴇ ᴄᴜʀʀᴇɴᴛ ᴩʟᴀʏɪɴɢ ᴛʀᴀᴄᴋ ᴏɴ ᴠɪᴅᴇᴏᴄʜᴀᴛ.",
			ThumbURL:    "https://telegra.ph/file/c5952790fa8235f499749.jpg",
			Text:        "/loop 3",
		},
	}
}
