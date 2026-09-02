package tools

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"ANJALI/config"

	tb "gopkg.in/tucnak/telebot.v2"
)

var startTime = time.Now()

func RegisterPingHandlers(b *tb.Bot) {
	// 1. Logic ko ek function variable mein daal diya
	pingFunc := func(m *tb.Message) {
		start := time.Now()
		cfg := config.LoadConfig()

		pingImg := "https://files.catbox.moe/6r97s4.jpg"
		if len(cfg.StartImgURL) > 0 {
			pingImg = cfg.StartImgURL[rand.Intn(len(cfg.StartImgURL))]
		}

		responseMsg, err := b.Send(m.Chat, &tb.Photo{
			File:    tb.FromURL(pingImg),
			Caption: "⚡ **Pinging system statistics...**",
		}, tb.ModeMarkdown)

		if err != nil {
			return
		}

		latency := time.Since(start).Milliseconds()
		uptime := time.Since(startTime).Round(time.Second).String()

		var mStats runtime.MemStats
		runtime.ReadMemStats(&mStats)
		ramUsage := fmt.Sprintf("%.2f MB", float64(mStats.Alloc)/1024/1024)

		caption := fmt.Sprintf(`<blockquote><b>🏓 ᴘᴏɴɢ !</b>

<b>• ʟᴀᴛᴇɴᴄʏ :</b> <code>%d ms</code>
<b>• ᴜᴘᴛɪᴍᴇ :</b> <code>%s</code>
<b>• ʀᴀᴍ :</b> <code>%s</code>
<b>• ʙᴏᴛ :</b> %s</blockquote>`, latency, uptime, ramUsage, b.Me.FirstName)

		markup := &tb.ReplyMarkup{}
		markup.Inline(markup.Row(markup.URL("✨ sᴜᴘᴘᴏʀᴛ", cfg.SupportChat)))

		b.EditCaption(responseMsg, caption, markup, tb.ModeHTML)
	}

	// 2. Us same function ko dono commands ke liye register kar diya!
	b.Handle("/ping", pingFunc)
	b.Handle("/alive", pingFunc)
}
