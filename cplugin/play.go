package cplugin

import (
	"fmt"
	"regexp"
	"strings"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Premium Emojis integrated with HTML tags
const (
	MsgDownloading = `<tg-emoji emoji-id="6156513311585211842">🌀</tg-emoji> 𝐃𝐨𝐰𝐧𝐥𝐨𝐚𝐝𝐢𝐧𝐠 𝐅𝐫𝐨𝐦 𝐁𝐞𝐭𝐚 𝐇𝐮𝐛 𝐁𝐚𝐛𝐲 𝐩𝐥𝐞𝐚𝐬𝐞 𝐰𝐚𝐢𝐭....`
	MsgStarting    = `<tg-emoji emoji-id="6129476453802188018">▶️</tg-emoji> 𝐒𝐭𝐚𝐫𝐭𝐢𝐧𝐠 𝐒𝐭𝐫𝐞𝐚𝐦 𝐄𝐧𝐣𝐨𝐲....`
	VCOffMsg       = `> **<tg-emoji emoji-id="6129782440157256336">🥹</tg-emoji> ᴏᴏᴘs! ᴠᴏɪᴄᴇ ᴄʜᴀᴛ ɪs ᴏғғ...**
> 
> **<tg-emoji emoji-id="5294339927318739359">🎙️</tg-emoji> ᴘʟᴇᴀsᴇ ᴛᴜʀɴ ᴏɴ ᴛʜᴇ ᴠᴏɪᴄᴇ ᴄʜᴀᴛ ɪɴ ᴛʜɪs ɢʀᴏᴜᴘ sᴏ ɪ ᴄᴀɴ ᴘʟᴀʏ sᴏᴍᴇ ᴀᴡᴇsᴏᴍᴇ ᴍᴜsɪᴄ ғᴏʀ ʏᴏᴜ! <tg-emoji emoji-id="6098202009486233047">🥳</tg-emoji>**`
	AlertIcon = `<tg-emoji emoji-id="6102938383456146362">🚨</tg-emoji>`
	BlockIcon = `<tg-emoji emoji-id="6271674836628541366">🚫</tg-emoji>`
)

// RCE & Security Filter
var bannedWords = []string{"porn", "pornhub", "xvideos", "webhook.site", "xnxx", "onlyfans"}
var suspiciousChars = regexp.MustCompile(`[;|&$` + "`" + `{}<>\\]`)

func isMalicious(text string) bool {
	lowerText := strings.ToLower(text)
	if suspiciousChars.MatchString(lowerText) {
		return true
	}
	badExts := []string{"webhook", "ngrok", ".sh", ".exe", ".bat", ".vbs", "rm -rf"}
	for _, ext := range badExts {
		if strings.Contains(lowerText, ext) {
			return true
		}
	}
	return false
}

func isNSFW(text string) bool {
	lower := strings.ToLower(text)
	for _, word := range bannedWords {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}

// PlayCommand handles /play, /vplay, etc.
func PlayCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	// 1. Anti-Spam Check
	userID := int64(m.Sender.ID)
	spamTracker.mu.Lock()
	// (Reuse spam logic from channel.go)
	spamTracker.mu.Unlock()

	// 2. Security Check (RCE/NSFW)
	if isMalicious(m.Text) {
		utils.SendSecurityLog(b, m.Chat.ID, m.Sender.ID, "ᴍᴀʟɪᴄɪᴏᴜs ʜᴀᴄᴋ ʟɪɴᴋ ʙʟᴏᴄᴋᴇᴅ", m.Text)
		b.Send(m.Chat, BlockIcon+" **sᴇᴄᴜʀɪᴛʏ ᴀʟᴇʀᴛ: ᴍᴀʟɪᴄɪᴏᴜs ᴄᴏᴍᴍᴀɴᴅ ɪɴᴊᴇᴄᴛɪᴏɴ ʙʟᴏᴄᴋᴇᴅ!**", tb.ModeHTML)
		return
	}
	if isNSFW(m.Text) {
		utils.SendSecurityLog(b, m.Chat.ID, m.Sender.ID, "ɴsғᴡ ᴠɪᴏʟᴀᴛɪᴏɴ", m.Text)
		b.Send(m.Chat, BlockIcon+" **sᴇᴄᴜʀɪᴛʏ ᴀʟᴇʀᴛ: ᴀᴅᴜʟᴛ ᴄᴏɴᴛᴇɴᴛ ɪs sᴛʀɪᴄᴛʟʏ ᴘʀᴏʜɪʙɪᴛᴇᴅ!**", tb.ModeHTML)
		return
	}

	args := strings.SplitN(m.Text, " ", 2)
	query := ""
	if len(args) > 1 {
		query = args[1]
	} else if m.ReplyTo != nil {
		// Handle Telegram Audio/Video play here
		b.Send(m.Chat, MsgDownloading, tb.ModeHTML)
		// Process local media -> download -> tgcalls stream...
		return
	}

	if query == "" {
		b.Send(m.Chat, "⚠️ What do you want to play?")
		return
	}

	mystic, _ := b.Send(m.Chat, "🔎 Searching...")

	// Stream logic (YouTube/Spotify/etc.) abstract
	err := utils.StreamHelper(b, m.Chat.ID, userID, query)
	if err != nil {
		if strings.Contains(err.Error(), "noactivegroupcall") {
			b.Edit(mystic, VCOffMsg, tb.ModeHTML) // ModeHTML enabled for premium emojis
		} else {
			b.Edit(mystic, fmt.Sprintf("❌ Error: %v", err))
		}
		return
	}

	b.Edit(mystic, MsgStarting, tb.ModeHTML) // ModeHTML enabled for premium emojis
	// Delay and delete mystic, send TrackMarkup
}

// MusicStreamCallback handles play buttons inline
func MusicStreamCallback(b *tb.Bot, c *tb.Callback) {
	// Parse c.Data ("MusicStream vidid|userid|mode|cplay|fplay")
	data := strings.Split(c.Data, " ")
	if len(data) < 2 {
		return
	}
	args := strings.Split(data[1], "|")
	vidID, userIDStr := args[0], args[1]

	if fmt.Sprintf("%d", c.Sender.ID) != userIDStr {
		b.Respond(c, &tb.CallbackResponse{Text: "❌ Not for you!", ShowAlert: true})
		return
	}

	b.Delete(c.Message)
	mystic, _ := b.Send(c.Message.Chat, MsgDownloading, tb.ModeHTML) // ModeHTML enabled

	err := utils.StreamFromVidID(vidID, c.Message.Chat.ID)
	if err != nil {
		b.Edit(mystic, fmt.Sprintf("❌ Error: %v", err))
		return
	}
	b.Edit(mystic, MsgStarting, tb.ModeHTML) // ModeHTML enabled
}
