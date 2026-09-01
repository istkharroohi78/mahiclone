package cplugin

import (
	"fmt"
	"runtime"

	"ANJALI/utils"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	tb "gopkg.in/tucnak/telebot.v2"
)

func GetRandomStatsImg() string {
	return "https://files.catbox.moe/6r97s4.jpg"
}

// StatsCommand handles /stats and /gstats
func StatsCommand(b *tb.Bot, m *tb.Message) {
	markup := utils.StatsButtons(isAdminCheck(b, m.Chat, int64(m.Sender.ID)))

	caption := fmt.Sprintf("📊 **Stats for %s**\nClick below to fetch global stats.", b.Me.FirstName)
	photo := &tb.Photo{File: tb.FromURL(GetRandomStatsImg()), Caption: caption}

	b.Send(m.Chat, photo, markup, tb.ModeMarkdown)
}

// StatsCallback handles 'bot_stats_sudo' and 'TopOverall'
func StatsCallback(b *tb.Bot, c *tb.Callback) {
	if !isAdminCheck(b, c.Message.Chat, int64(c.Sender.ID)) {
		b.Respond(c, &tb.CallbackResponse{Text: "❌ Sudo required!", ShowAlert: true})
		return
	}

	b.Respond(c, &tb.CallbackResponse{Text: "Fetching Stats..."})

	vm, _ := mem.VirtualMemory()
	ram := fmt.Sprintf("%.2f GB", float64(vm.Total)/(1024*1024*1024))

	pCore, _ := cpu.Counts(false)
	tCore, _ := cpu.Counts(true)

	diskStat, _ := disk.Usage("/")
	totalDisk := float64(diskStat.Total) / (1024 * 1024 * 1024)
	usedDisk := float64(diskStat.Used) / (1024 * 1024 * 1024)
	freeDisk := float64(diskStat.Free) / (1024 * 1024 * 1024)

	text := fmt.Sprintf(`**🖥 Server Stats for %s:**

**OS:** %s
**RAM:** %s
**Cores (P/T):** %d / %d
**Go Version:** %s

**Disk (Total/Used/Free):**
%.2fGB / %.2fGB / %.2fGB`,
		b.Me.FirstName, runtime.GOOS, ram, pCore, tCore, runtime.Version(),
		totalDisk, usedDisk, freeDisk)

	markup := utils.BackStatsButtons()
	photo := &tb.Photo{File: tb.FromURL(GetRandomStatsImg()), Caption: text}

	b.Edit(c.Message, photo, markup, tb.ModeMarkdown)
}
