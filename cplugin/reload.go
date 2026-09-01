package cplugin

import (
	"fmt"
	"time"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

var reloadLimits = make(map[int64]int64)

// ReloadCommand handles /reload, /admincache
func ReloadCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	now := time.Now().Unix()
	if lastRun, exists := reloadLimits[m.Chat.ID]; exists && lastRun > now {
		b.Send(m.Chat, fmt.Sprintf("⏳ Please wait %s before reloading again.", utils.GetReadableTime(int(lastRun-now))), tb.ModeMarkdown)
		return
	}

	// Rebuild admin cache
	admins, err := b.AdminsOf(m.Chat)
	if err == nil {
		var adminIDs []int64
		for _, a := range admins {
			adminIDs = append(adminIDs, int64(a.User.ID))
		}
		utils.UpdateAdminCache(m.Chat.ID, adminIDs)
	}

	reloadLimits[m.Chat.ID] = now + 180 // 3 minutes cooldown
	b.Send(m.Chat, "✅ **Admin cache successfully reloaded!**", tb.ModeMarkdown)
}

// RebootCommand handles /reboot
func RebootCommand(b *tb.Bot, m *tb.Message) {
	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		return
	}

	mystic, _ := b.Send(m.Chat, fmt.Sprintf("🔄 **%s is rebooting...**", b.Me.FirstName), tb.ModeMarkdown)

	utils.ClearChatQueue(m.Chat.ID)
	utils.ForceStopStream(m.Chat.ID)

	time.Sleep(1 * time.Second)
	b.Edit(mystic, fmt.Sprintf("✅ **%s successfully rebooted for this chat!**", b.Me.FirstName), tb.ModeMarkdown)
}

// CloseMenuCallback handles 'close' buttons
func CloseMenuCallback(b *tb.Bot, c *tb.Callback) {
	b.Respond(c, &tb.CallbackResponse{})
	b.Delete(c.Message)

	tempMsg, _ := b.Send(c.Message.Chat, fmt.Sprintf("ᴄʟᴏꜱᴇ ʙʏ : [%s](tg://user?id=%d)", c.Sender.FirstName, c.Sender.ID), tb.ModeMarkdown)
	go func() {
		time.Sleep(2 * time.Second)
		b.Delete(tempMsg)
	}()
}
