package play

import (
	"fmt"
	"strings"

	"ANJALI/config"
	"ANJALI/utils/decorators"
	"ANJALI/utils/stream"

	tb "gopkg.in/tucnak/telebot.v2"
)

var BannedWords = []string{
	"porn", "pornhub", "xvideos", "xnxx", "brazzers",
	"onlyfans", "xhamster", "hot bhabhi", "redtube",
	"child porn", "loli", "shota", "incest", "bestiality",
}

const (
	MsgDownloading = "<tg-emoji emoji-id=\"6156513311585211842\">🌀</tg-emoji> 𝐃𝐨𝐰𝐧𝐥𝐨𝐚𝐝𝐢𝐧𝐠 𝐁𝐚𝐛𝐲 𝐩𝐥𝐞𝐚𝐬𝐞 𝐰𝐚𝐢𝐭...."
	MsgStarting    = "<tg-emoji emoji-id=\"6129476453802188018\">▶️</tg-emoji> 𝐒𝐭𝐚𝐫𝐭𝐢𝐧𝐠 𝐒𝐭𝐫𝐞𝐚𝐦 𝐄𝐧𝐣𝐨𝐲...."
	FallbackImg    = "https://telegra.ph/file/2e3d368e77c449c287430.jpg"
)

func isNsfwContent(text string) bool {
	lowerText := strings.ToLower(text)
	for _, word := range BannedWords {
		if strings.Contains(lowerText, word) {
			return true
		}
	}
	return false
}

func isMaliciousLink(text string) bool {
	lowerText := strings.ToLower(text)
	badExtensions := []string{".sh", ".exe", ".bat", "ngrok", "webhook"}
	for _, ext := range badExtensions {
		if strings.Contains(lowerText, ext) {
			return true
		}
	}
	return false
}

func RegisterPlayHandlers(b *tb.Bot) {
	// God-mode filter checking
	b.Handle(tb.OnText, func(m *tb.Message) {
		if isMaliciousLink(m.Text) || isNsfwContent(m.Text) {
			b.Delete(m)
			b.Send(m.Chat, "🚫 **sᴇᴄᴜʀɪᴛʏ ᴀʟᴇʀᴛ: ᴍᴀʟɪᴄɪᴏᴜs / ᴀᴅᴜʟᴛ ᴄᴏɴᴛᴇɴᴛ ɪs sᴛʀɪᴄᴛʟʏ ᴘʀᴏʜɪʙɪᴛᴇᴅ!**", tb.ModeMarkdown)

			// Send to Secure Logger
			cfg := config.LoadConfig()
			if cfg.LoggerID != 0 {
				logChat, _ := b.ChatByID(fmt.Sprintf("%d", cfg.LoggerID))
				b.Send(logChat, fmt.Sprintf("🚨 **SECURITY BREACH**\nUser: [%s](tg://user?id=%d)\nContent: `%s`", m.Sender.FirstName, m.Sender.ID, m.Text), tb.ModeMarkdown)
			}
			return
		}
	})

	playCmds := []string{"/play", "/vplay", "/cplay", "/cvplay", "/playforce"}
	for _, cmd := range playCmds {
		// FIXED: Using PlayWrapper as a middleware correctly with PlayContext
		b.Handle(cmd, decorators.PlayWrapper(b, func(m *tb.Message, ctx *decorators.PlayContext) {

			if isNsfwContent(ctx.URL) || isNsfwContent(m.Text) {
				b.Send(m.Chat, "🚫 **sᴇᴄᴜʀɪᴛʏ ᴀʟᴇʀᴛ: ᴀᴅᴜʟᴛ ᴄᴏɴᴛᴇɴᴛ ɪs sᴛʀɪᴄᴛʟʏ ᴘʀᴏʜɪʙɪᴛᴇᴅ!**")
				return
			}

			mystic, _ := b.Send(m.Chat, MsgDownloading, tb.ModeHTML)

			// Detect Telegram Media
			if m.ReplyTo != nil && (m.ReplyTo.Audio != nil || m.ReplyTo.Voice != nil || m.ReplyTo.Video != nil || m.ReplyTo.Document != nil) {
				b.Edit(mystic, MsgStarting, tb.ModeHTML)
				details := map[string]string{"title": "Telegram File", "duration_min": "Unknown", "vidid": "tg_file"}
				
				// FIXED: stream.Stream does not return an error, just execute it
				stream.Stream(b, m, ctx.ChatID, m.Chat.ID, details, ctx.Video, "telegram", ctx.FPlay)
				b.Delete(mystic)
				return
			}

			if ctx.URL != "" {
				b.Edit(mystic, MsgStarting, tb.ModeHTML)
				// External URL & Search Logic
				details := map[string]string{
					"title":        "YouTube Audio",
					"duration_min": "03:45",
					"vidid":        "dummy_id",
					"thumb":        FallbackImg,
				}
				
				// FIXED: stream.Stream does not return an error, removed err assignment
				stream.Stream(b, m, ctx.ChatID, m.Chat.ID, details, ctx.Video, "youtube", ctx.FPlay)
				b.Delete(mystic)
			}
		}))
	}
}
