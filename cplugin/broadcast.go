package cplugin

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"ANJALI/config"
	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// PremiumEmojis list (Updated with the new ones provided)
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

const PoweredByHtml = "\n\n🤞 𝐏ᴏᴡєʀєᴅ 𝐁ʏ ➛ <a href=\"https://t.me/betabot_hub\">[˹the shiv.🙂❤️˼]</a>"

// Global flag to prevent overlapping broadcasts
var isBroadcasting bool

// BroadcastCommand handles the /broadcast command
func BroadcastCommand(b *tb.Bot, m *tb.Message) {
	botID := b.Me.ID
	userID := m.Sender.ID

	// Clone Owner Check
	cloneData := utils.GetCloneBotData(botID)
	cfg := config.LoadConfig()

	if cloneData == nil || (cloneData.UserID != int64(userID) && int64(userID) != cfg.OwnerID) {
		// Normal unauthorized error for users who aren't the owner
		b.Send(m.Chat, fmt.Sprintf("❌ You are not authorized to use this command. Ask the bot owner or visit [Support](%s).", cfg.SupportChat), tb.ModeMarkdown)
		return
	}

	if isBroadcasting {
		b.Send(m.Chat, "⏳ **Broadcast is already running. Please wait.**", tb.ModeMarkdown)
		return
	}

	query := ""
	if m.ReplyTo != nil {
		args := strings.SplitN(m.Text, " ", 2)
		if len(args) > 1 {
			query = args[1]
		}
	} else {
		args := strings.SplitN(m.Text, " ", 2)
		if len(args) < 2 {
			b.Send(m.Chat, "⚠️ Please provide a message or reply to a message to broadcast.", tb.ModeMarkdown)
			return
		}
		query = args[1]
	}

	flags := []string{"-pin", "-nobot", "-pinloud", "-user"}
	queryToSend := query
	for _, flag := range flags {
		queryToSend = strings.ReplaceAll(queryToSend, flag, "")
	}
	queryToSend = strings.TrimSpace(queryToSend)

	if m.ReplyTo == nil && queryToSend == "" {
		b.Send(m.Chat, "❌ Cannot send an empty broadcast.", tb.ModeMarkdown)
		return
	}

	isBroadcasting = true
	defer func() { isBroadcasting = false }()

	b.Send(m.Chat, "📢 **Broadcast started...**", tb.ModeMarkdown)

	// Add a random premium emoji to the text just for styling flair (optional, Telebot doesn't natively render custom emojis but standard clients will if sent as text)
	randEmoji := PremiumEmojis[rand.Intn(len(PremiumEmojis))]
	_ = randEmoji // Stored to prevent unused variable error, can be injected into text if desired.

	// PART A: BROADCAST TO GROUPS
	if !strings.Contains(m.Text, "-nobot") {
		sent := 0
		pinCount := 0
		servedChats := utils.GetServedChatsClone(botID)

		for _, chat := range servedChats {
			targetChat, err := b.ChatByID(fmt.Sprintf("%d", chat.ChatID))
			if err != nil {
				continue
			}

			var sentMsg *tb.Message
			if m.ReplyTo != nil {
				sentMsg, err = b.Forward(targetChat, m.ReplyTo)
			} else {
				sentMsg, err = b.Send(targetChat, queryToSend)
			}

			if err == nil {
				sent++
				if strings.Contains(m.Text, "-pin") || strings.Contains(m.Text, "-pinloud") {
					notify := !strings.Contains(m.Text, "-pinloud") // pinloud means disable_notification=false
					err = b.Pin(sentMsg, notify)
					if err == nil {
						pinCount++
					}
				}
			}
			time.Sleep(200 * time.Millisecond) // Flood protection
		}

		resultText := fmt.Sprintf("✅ **Broadcast Completed!**\n\n🏢 **Sent to Groups:** `%d`\n📌 **Pinned in:** `%d`%s", sent, pinCount, PoweredByHtml)
		b.Send(m.Chat, resultText, tb.ModeHTML)
	}

	// PART B: BROADCAST TO USERS
	if strings.Contains(m.Text, "-user") {
		susr := 0
		servedUsers := utils.GetServedUsersClone(botID)

		for _, user := range servedUsers {
			targetUser := &tb.User{ID: int(user.UserID)}

			var err error
			if m.ReplyTo != nil {
				_, err = b.Forward(targetUser, m.ReplyTo)
			} else {
				_, err = b.Send(targetUser, queryToSend)
			}

			if err == nil {
				susr++
			}
			time.Sleep(200 * time.Millisecond) // Flood protection
		}

		resultText := fmt.Sprintf("✅ **User Broadcast Completed!**\n\n👤 **Sent to Users:** `%d`%s", susr, PoweredByHtml)
		b.Send(m.Chat, resultText, tb.ModeHTML)
	}
}
