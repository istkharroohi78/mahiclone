package cplugin

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

const DeleteDelay = 7 * time.Second

var (
	enabledChats  = make(map[int64]bool)
	userJoinCount = make(map[string]int)   // Key format: "chatID_userID"
	userCache     = make(map[int64]string) // Map userID to Name
	vcMutex       sync.Mutex
)

// DeleteMessageAfterDelay acts like the async sleep-then-delete function
func DeleteMessageAfterDelay(b *tb.Bot, m *tb.Message) {
	go func() {
		time.Sleep(DeleteDelay)
		b.Delete(m)
	}()
}

// SendJoinNotification formats and sends the join alert
func SendJoinNotification(b *tb.Bot, chatID int64, user *tb.User) {
	vcMutex.Lock()
	key := fmt.Sprintf("%d_%d", chatID, user.ID)
	userJoinCount[key]++
	count := userJoinCount[key]
	vcMutex.Unlock()

	name := user.FirstName
	if user.LastName != "" {
		name += " " + user.LastName
	}
	username := "Ignored"
	if user.Username != "" {
		username = "@" + user.Username
	}

	mention := fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, user.ID, name)

	text := fmt.Sprintf("<b>#JoinVideoChat</b>\n\n<b>● ɴᴀᴍᴇ ➛</b> %s\n<b>● ɪᴅ ➛</b><code>%d</code>\n<b>● ᴜsᴇʀɴᴀᴍᴇ ➛</b> %s", mention, user.ID, username)

	if count > 1 {
		text += fmt.Sprintf("\n\n<b>🔁 ᴊᴏɪɴ ᴄᴏᴜɴᴛ ➛</b> <code>%d</code>", count)
	}

	chat, err := b.ChatByID(fmt.Sprintf("%d", chatID))
	if err == nil {
		msg, _ := b.Send(chat, text, tb.ModeHTML)
		if msg != nil {
			DeleteMessageAfterDelay(b, msg)
		}
	}
}

// SendLeaveNotification formats and sends the leave alert
func SendLeaveNotification(b *tb.Bot, chatID int64, user *tb.User) {
	name := user.FirstName
	if user.LastName != "" {
		name += " " + user.LastName
	}
	username := "Ignored"
	if user.Username != "" {
		username = "@" + user.Username
	}

	mention := fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, user.ID, name)
	text := fmt.Sprintf("<b>#LeaveVideoChat</b>\n\n<b>● ɴᴀᴍᴇ ➛</b> %s\n<b>● ɪᴅ ➛</b><code>%d</code>\n<b>● ᴜsᴇʀɴᴀᴍᴇ ➛</b> %s", mention, user.ID, username)

	chat, err := b.ChatByID(fmt.Sprintf("%d", chatID))
	if err == nil {
		msg, _ := b.Send(chat, text, tb.ModeHTML)
		if msg != nil {
			DeleteMessageAfterDelay(b, msg)
		}
	}
}

// VCLoggerCmd handles /vclogger and /vclog
func VCLoggerCmd(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "**❌ ᴀᴅᴍɪɴ ᴏɴʟʏ!**", tb.ModeMarkdown)
		return
	}

	args := strings.Split(m.Text, " ")
	if len(args) < 2 {
		vcMutex.Lock()
		status := enabledChats[m.Chat.ID]
		vcMutex.Unlock()

		statusText := "❌ OFF"
		if status {
			statusText = "✅ ON"
		}

		b.Send(m.Chat, fmt.Sprintf("**📊 ᴠᴄ ʟᴏɢɢᴇʀ :** %s\n\n**ᴄᴏᴍᴍᴀɴᴅs :**\n\n**• /vclogger on**\n•** /vclogger off**", statusText), tb.ModeMarkdown)
		return
	}

	action := strings.ToLower(args[1])
	vcMutex.Lock()
	defer vcMutex.Unlock()

	if action == "on" {
		enabledChats[m.Chat.ID] = true
		b.Send(m.Chat, "**✅ ᴠᴄ ʟᴏɢɢᴇʀ ᴇɴᴀʙʟᴇᴅ!**", tb.ModeMarkdown)
	} else if action == "off" {
		delete(enabledChats, m.Chat.ID)
		// Clear counts for this chat
		for k := range userJoinCount {
			if strings.HasPrefix(k, fmt.Sprintf("%d_", m.Chat.ID)) {
				delete(userJoinCount, k)
			}
		}
		b.Send(m.Chat, "**❌ ᴠᴄ ʟᴏɢɢᴇʀ ᴅɪsᴀʙʟᴇᴅ!**", tb.ModeMarkdown)
	} else {
		b.Send(m.Chat, "**ᴜsᴇ:** `/vclogger on | off`", tb.ModeMarkdown)
	}
}
