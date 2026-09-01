package cplugin

import (
	"fmt"
	"log"
	"strings"

	"ANJALI/config"
	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// isAdminCheck validates if the user is an admin or owner (Bulletproof Check)
func isAdminCheck(b *tb.Bot, chat *tb.Chat, userID int64) bool {
	cfg := config.LoadConfig()
	if userID == cfg.OwnerID {
		return true
	}

	member, err := b.ChatMemberOf(chat, &tb.User{ID: int(userID)})
	if err != nil {
		return false
	}

	if member.Role == tb.Administrator || member.Role == tb.Creator {
		return true
	}
	return false
}

// AuthCommand handles /auth and /cauth
func AuthCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	// 🟢 THE FIX: BULLETPROOF ADMIN CHECK
	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Sirf Admins he is command ko use kar sakte hain!**", tb.ModeMarkdown)
		return
	}

	args := strings.Split(m.Text, " ")
	var targetUser *tb.User

	if m.ReplyTo != nil {
		targetUser = m.ReplyTo.Sender
	} else if len(args) == 2 {
		// Yahan par username se ID nikalne ka logic lagega (extract_user)
		// Go telebot me direct ID pass karni padti hai ya API se fetch karna hota hai
		targetUser = &tb.User{ID: 0} // Placeholder for extract_user logic
	} else {
		b.Send(m.Chat, "User ko reply karein ya uski ID/Username dein.", tb.ModeMarkdown)
		return
	}

	token := utils.IntToAlpha(int64(targetUser.ID))
	authNames := utils.GetAuthUserNames(m.Chat.ID)

	if len(authNames) >= 25 {
		b.Send(m.Chat, "❌ You can only have 25 authorized users.", tb.ModeMarkdown)
		return
	}

	alreadyAuth := false
	for _, n := range authNames {
		if n == token {
			alreadyAuth = true
			break
		}
	}

	if !alreadyAuth {
		assis := map[string]interface{}{
			"auth_user_id": targetUser.ID,
			"auth_name":    targetUser.FirstName,
			"admin_id":     m.Sender.ID,
			"admin_name":   m.Sender.FirstName,
		}

		utils.SaveAuthUser(m.Chat.ID, token, assis)

		// Update AdminList Cache
		utils.AddAdminCache(m.Chat.ID, int64(targetUser.ID))

		b.Send(m.Chat, fmt.Sprintf("✅ Added [%s](tg://user?id=%d) to Authorized Users List.", targetUser.FirstName, targetUser.ID), tb.ModeMarkdown)
	} else {
		b.Send(m.Chat, fmt.Sprintf("⚠️ [%s](tg://user?id=%d) is already in the Authorized Users List.", targetUser.FirstName, targetUser.ID), tb.ModeMarkdown)
	}
}

// UnauthCommand handles /unauth and /cunauth
func UnauthCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	// 🟢 THE FIX: BULLETPROOF ADMIN CHECK
	if !isAdminCheck(b, m.Chat, int64(m.Sender.ID)) {
		b.Send(m.Chat, "❌ **Sirf Admins he is command ko use kar sakte hain!**", tb.ModeMarkdown)
		return
	}

	args := strings.Split(m.Text, " ")
	var targetUser *tb.User

	if m.ReplyTo != nil {
		targetUser = m.ReplyTo.Sender
	} else if len(args) == 2 {
		targetUser = &tb.User{ID: 0} // Placeholder for extract_user logic
	} else {
		b.Send(m.Chat, "User ko reply karein ya uski ID/Username dein.", tb.ModeMarkdown)
		return
	}

	token := utils.IntToAlpha(int64(targetUser.ID))
	deleted := utils.DeleteAuthUser(m.Chat.ID, token)

	if deleted {
		utils.RemoveAdminCache(m.Chat.ID, int64(targetUser.ID))
		b.Send(m.Chat, fmt.Sprintf("✅ Removed [%s](tg://user?id=%d) from Authorized Users List.", targetUser.FirstName, targetUser.ID), tb.ModeMarkdown)
	} else {
		b.Send(m.Chat, fmt.Sprintf("⚠️ [%s](tg://user?id=%d) is not in the Authorized Users List.", targetUser.FirstName, targetUser.ID), tb.ModeMarkdown)
	}
}

// AuthListCommand handles /authlist, /authusers, /cauthlist
func AuthListCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	authNames := utils.GetAuthUserNames(m.Chat.ID)
	if len(authNames) == 0 {
		b.Send(m.Chat, "❌ No authorized users found in this group.", tb.ModeMarkdown)
		return
	}

	waitingMsg, _ := b.Send(m.Chat, "🔄 **Fetching Authorized Users...**", tb.ModeMarkdown)
	text := fmt.Sprintf("📜 **Authorized Users List for %s:**\n\n", m.Chat.Title)

	count := 0
	for _, token := range authNames {
		authData := utils.GetAuthUser(m.Chat.ID, token)
		if authData == nil {
			continue
		}

		userID := int(authData["auth_user_id"].(float64))
		adminID := int(authData["admin_id"].(float64))
		adminName := authData["admin_name"].(string)

		userObj, err := b.ChatByID(fmt.Sprintf("%d", userID))
		userName := "Unknown"
		if err == nil {
			userName = userObj.FirstName
		}

		count++
		text += fmt.Sprintf("%d➤ %s [`%d`]\n   Added By: %s [`%d`]\n\n", count, userName, userID, adminName, adminID)
	}

	markup := utils.CloseMarkup()
	_, err := b.Edit(waitingMsg, text, markup, tb.ModeHTML)
	if err != nil {
		log.Println("Failed to edit authlist message:", err)
	}
}
