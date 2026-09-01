package cplugin

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"ANJALI/config"
	"ANJALI/utils"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Boot time tracking
var bootTime = time.Now()

func GetRandomPingImg() string {
	cfg := config.LoadConfig()
	if len(cfg.PingImgURLs) > 0 {
		return cfg.PingImgURLs[rand.Intn(len(cfg.PingImgURLs))]
	}
	return "https://telegra.ph/file/2e3d368e77c449c287430.jpg"
}

// PingCommand handles /ping
func PingCommand(b *tb.Bot, m *tb.Message) {
	start := time.Now()
	botID := b.Me.ID

	cfg := config.LoadConfig()
	cloneData := utils.GetCloneBotData(botID)
	supportChat := cfg.SupportChat
	if cloneData != nil && cloneData.SupportChat != "" {
		supportChat = cloneData.SupportChat
	}

	captionText := fmt.Sprintf("%s is pinging...", b.Me.FirstName)
	msg, _ := b.Send(m.Chat, captionText)

	// Fetch custom ping media
	pingImg := utils.GetClonePingImage(botID)

	// Stats calculation
	upt := int(time.Since(bootTime).Seconds())
	cpuPercent, _ := cpu.Percent(0, false)
	cpuVal := 0.0
	if len(cpuPercent) > 0 {
		cpuVal = cpuPercent[0]
	}
	vm, _ := mem.VirtualMemory()
	diskStat, _ := disk.Usage("/")
	resp := time.Since(start).Milliseconds()

	statsText := fmt.Sprintf("➻ Pong: %dms\n\n%s System Stats:\n\n๏ Uptime: %s\n๏ Ram: %.2f%%\n๏ Cpu: %.2f%%\n๏ Disk: %.2f%%",
		resp, b.Me.FirstName, utils.GetReadableTime(upt), vm.UsedPercent, cpuVal, diskStat.UsedPercent)

	markup := &tb.ReplyMarkup{}
	btnSupport := CreateBtn(markup, "Support", "", supportChat, Primary, false)
	markup.Inline(markup.Row(btnSupport))

	// If custom image exists, edit message to photo
	if pingImg != "" {
		photo := &tb.Photo{File: tb.FromURL(pingImg), Caption: statsText}
		b.Edit(msg, photo, markup, tb.ModeMarkdown)
	} else {
		photo := &tb.Photo{File: tb.FromURL(GetRandomPingImg()), Caption: statsText}
		b.Edit(msg, photo, markup, tb.ModeMarkdown)
	}
}

// SetPingImage handles /setpingimg and /addpingimg
func SetPingImage(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "🚫 Only the Bot Owner can use this command.", tb.ModeMarkdown)
		return
	}

	if m.ReplyTo != nil && m.ReplyTo.Photo != nil {
		utils.SetClonePingImage(b.Me.ID, m.ReplyTo.Photo.FileID)
		b.Send(m.Chat, "✅ Photo added to Ping List!")
		return
	}

	args := strings.SplitN(m.Text, " ", 2)
	if len(args) == 2 {
		utils.SetClonePingImage(b.Me.ID, args[1])
		b.Send(m.Chat, "✅ Image URL added to Ping List!")
		return
	}
	b.Send(m.Chat, "❗ Reply to a photo or use `/addpingimg [URL]`", tb.ModeMarkdown)
}

// DelPingImage handles /delpingimg
func DelPingImage(b *tb.Bot, m *tb.Message) {
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		b.Send(m.Chat, "🚫 Only the Bot Owner can use this command.", tb.ModeMarkdown)
		return
	}
	utils.DeleteClonePingImage(b.Me.ID)
	b.Send(m.Chat, "🗑 All custom ping images removed! (Reset to Default)")
}
