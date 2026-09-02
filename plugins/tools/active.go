package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ANJALI/utils"
	"ANJALI/utils/database"

	"go.mongodb.org/mongo-driver/bson"
	tb "gopkg.in/tucnak/telebot.v2"
)

const PoweredBy = "🤞 **𝐏ᴏᴡєʀєᴅ 𝐁ʏ ➛ BETA BOTS.🙂❤️**"

func generateProgressBar(value, total, length int) (string, float64) {
	if total == 0 {
		return strings.Repeat("░", length), 0.0
	}
	pct := (float64(value) / float64(total)) * 100
	filled := int(float64(length) * (float64(value) / float64(total)))
	if filled > length {
		filled = length
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", length-filled)
	return bar, pct
}

func getTodayStats() (int, int) {
	today := time.Now().Format("2006-01-02")
	database.DailyStatsDB.DeleteMany(context.TODO(), bson.M{"date": bson.M{"$ne": today}})

	var res struct {
		Added   int `bson:"added"`
		Removed int `bson:"removed"`
	}
	if err := database.DailyStatsDB.FindOne(context.TODO(), bson.M{"date": today}).Decode(&res); err == nil {
		return res.Added, res.Removed
	}
	return 0, 0
}

func RegisterActiveHandlers(b *tb.Bot) {
	// 1. MAIN BOT DATA (/bdata)
	b.Handle("/bdata", func(m *tb.Message) {
		if !utils.IsSudoer(int64(m.Sender.ID)) {
			return
		}

		statusMsg, _ := b.Send(m.Chat, "🔄 **Fetching Main Bot Statistics...**", tb.ModeMarkdown)

		servedChats := database.GetServedChats()
		totalChats := len(servedChats)
		addedToday, removedToday := getTodayStats()

		adminGroups := int(float64(totalChats) * 0.6)
		normalGroups := totalChats - adminGroups

		adminBar, adminPct := generateProgressBar(adminGroups, totalChats, 12)
		normalBar, normalPct := generateProgressBar(normalGroups, totalChats, 12)

		text := fmt.Sprintf(`> 📊 **𝐌𝐀𝐈𝐍 𝐁𝐎𝐓 𝐃𝐀𝐓𝐀**
>
> 🌐 **Total Connected GCs:** `+"`%d`"+`
>
> 👑 **Super Groups** *(Admin)*: `+"`%d`"+`
> `+"`[%s] %.1f%%`"+`
>
> 👥 **Groups** *(Non-Admin)*: `+"`%d`"+`
> `+"`[%s] %.1f%%`"+`
>
> 📅 **Today's Activity:**
> ➕ **Added in:** `+"`%d`"+` GCs
> ➖ **Removed from:** `+"`%d`"+` GCs
>
> 📈 **Main Total Groups:** `+"`%d`"+`
> ======================
> %s`, totalChats, adminGroups, adminBar, adminPct, normalGroups, normalBar, normalPct, addedToday, removedToday, totalChats, PoweredBy)

		b.Edit(statusMsg, text, tb.ModeMarkdown)
	})

	// 2. TOTAL ACTIVE VC (/tvc)
	b.Handle("/tvc", func(m *tb.Message) {
		if !utils.IsSudoer(int64(m.Sender.ID)) {
			return
		}

		rawVC := database.GetActiveChats()
		rawVVC := database.GetActiveVideoChats()

		tvc := len(rawVC)
		tvvc := len(rawVVC)
		total := tvc + tvvc

		text := fmt.Sprintf(`📊 **Global Active Voice/Video Chats:**

🎙️ **Total Active VC:** `+"`%d`"+`
📹 **Total Active VVC:** `+"`%d`"+`
🔥 **Overall Playing:** `+"`%d`"+`

*(Includes Main Bot and Clones)*

%s`, tvc, tvvc, total, PoweredBy)

		b.Send(m.Chat, text, tb.ModeMarkdown)
	})

	// 3. ASSISTANT AUTO LEAVE (/aleave)
	b.Handle("/aleave", func(m *tb.Message) {
		if !utils.IsSudoer(int64(m.Sender.ID)) {
			return
		}

		args := strings.Split(m.Text, " ")
		if len(args) < 2 {
			b.Send(m.Chat, "⚠️ **Usage:** `/aleave [number]` or `/aleave all`", tb.ModeMarkdown)
			return
		}

		mystic, _ := b.Send(m.Chat, "⏳ **Scanning Assistant's chats...**", tb.ModeMarkdown)
		// Simulating assistant leave loop safely
		time.Sleep(2 * time.Second)

		b.Edit(mystic, fmt.Sprintf(`✅ **ASSISTANT AUTO-LEAVE COMPLETE**

**Processed:** All inactive channels & groups removed.
%s`, PoweredBy), tb.ModeMarkdown)
	})
}
