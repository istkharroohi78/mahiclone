package tools

import (
	"context"
	"fmt"
	"runtime"

	"ANJALI/config"
	"ANJALI/utils/database"
	"ANJALI/utils/inline"

	"go.mongodb.org/mongo-driver/bson"
	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterStatsHandlers(b *tb.Bot) {
	b.Handle("/stats", func(m *tb.Message) {
		cfg := config.LoadConfig()
		isOwnerOrSudo := int64(m.Sender.ID) == cfg.OwnerID || IsSudoer(int64(m.Sender.ID))

		// Restrict stats view strictly to Owner & Sudoers
		if !isOwnerOrSudo {
			b.Send(m.Chat, "❌ **Only the Bot Owner & Sudoers can view hardware statistics.**", tb.ModeMarkdown)
			return
		}

		var mStats runtime.MemStats
		runtime.ReadMemStats(&mStats)

		totalUsers := len(database.GetServedUsers())
		totalChats := len(database.GetServedChats())
		activeVCs := len(database.GetActiveChats())

		// MongoDB stats check
		mongoStatus := "Online 🟢"
		if err := database.CloneBotDB.Database().RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
			mongoStatus = "Offline 🔴"
		}

		caption := fmt.Sprintf(`<blockquote><b>📊 𝐒𝐘𝐒𝐓𝐄𝐌 & 𝐁𝐎𝐓 𝐒𝐓𝐀𝐓𝐈𝐒𝐓𝐈𝐂𝐒</b>

<b>• ʀᴀᴍ ᴜsᴀɢᴇ :</b> <code>%.2f MB</code>
<b>• ɢᴏʀᴏᴜᴛɪɴᴇs :</b> <code>%d</code>
<b>• ᴀᴄᴛɪᴠᴇ sᴛʀᴇᴀᴍs :</b> <code>%d</code>
<b>• sᴇʀᴠᴇᴅ ᴄʜᴀᴛs :</b> <code>%d</code>
<b>• sᴇʀᴠᴇᴅ ᴜsᴇʀs :</b> <code>%d</code>
<b>• ᴅᴀᴛᴀʙᴀsᴇ :</b> <code>%s</code>
<b>• ᴇɴɢɪɴᴇ :</b> <code>Golang v%s</code></blockquote>`,
			float64(mStats.Alloc)/1024/1024,
			runtime.NumGoroutine(),
			activeVCs,
			totalChats,
			totalUsers,
			mongoStatus,
			runtime.Version(),
		)

		markup := &tb.ReplyMarkup{}
		markup.Inline(markup.Row(inline.CreateBtn(markup, "🗑️ ᴄʟᴏsᴇ", "close", "", inline.Danger, 0, false)))

		b.Send(m.Chat, &tb.Photo{
			File:    tb.FromURL("https://files.catbox.moe/6r97s4.jpg"),
			Caption: caption,
		}, markup, tb.ModeHTML)
	})
}
