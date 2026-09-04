package admins

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	mystrings "ANJALI/i18n"
	"ANJALI/utils/database"
	"ANJALI/utils/inline"
)

// Helper: Check if sender is an actual admin in the group
func isAdmin(b *tb.Bot, m *tb.Message) bool {
	member, err := b.ChatMemberOf(m.Chat, m.Sender)
	if err != nil {
		return false
	}
	return member.Role == tb.Administrator || member.Role == tb.Creator
}

// Helper: Extract user from reply or command arguments
func extractUser(b *tb.Bot, m *tb.Message) (*tb.User, error) {
	if m.ReplyTo != nil && m.ReplyTo.Sender != nil {
		return m.ReplyTo.Sender, nil
	}
	
	args := strings.Split(m.Text, " ")
	if len(args) < 2 {
		return nil, fmt.Errorf("no user provided")
	}

	// Try parsing as ID
	userID, err := strconv.ParseInt(args[1], 10, 64)
	if err == nil {
		chat, err := b.ChatByID(fmt.Sprintf("%d", userID))
		if err == nil {
			return &tb.User{ID: chat.ID, FirstName: chat.FirstName, Username: chat.Username}, nil
		}
	}
	
	return nil, fmt.Errorf("user not found")
}

func RegisterAuthHandlers(b *tb.Bot) {

	// 1. Auth Command
	authHandler := func(m *tb.Message) {
		if !m.FromGroup() {
			return
		}
		if !isAdmin(b, m) {
			return 
		}

		langData := mystrings.GetString(database.GetLang(m.Chat.ID))
		targetUser, err := extractUser(b, m)
		if err != nil {
			b.Send(m.Chat, langData["general_1"])
			return
		}

		token := fmt.Sprintf("%d", targetUser.ID)
		authNames := database.GetAuthUserNames(m.Chat.ID)
		
		if len(authNames) >= 25 {
			b.Send(m.Chat, langData["auth_1"])
			return
		}

		isAlreadyAuth := false
		for _, name := range authNames {
			if name == token {
				isAlreadyAuth = true
				break
			}
		}

		if !isAlreadyAuth {
			authData := map[string]interface{}{
				"auth_user_id": targetUser.ID,
				"auth_name":    targetUser.FirstName,
				"admin_id":     m.Sender.ID,
				"admin_name":   m.Sender.FirstName,
			}
			
			// Update local cache
			database.AddAdminList(m.Chat.ID, targetUser.ID)
			
			// Save to DB: marshal map to string, and pass targetUser.ID as int64
			jsonData, _ := json.Marshal(authData)
			database.SaveAuthUser(m.Chat.ID, int64(targetUser.ID), string(jsonData))
			
			mention := fmt.Sprintf("<a href='tg://user?id=%d'>%s</a>", targetUser.ID, targetUser.FirstName)
			msg := strings.ReplaceAll(langData["auth_2"], "{0}", mention)
			b.Send(m.Chat, msg+"\n\n— the shiv", &tb.SendOptions{ParseMode: tb.ModeHTML})
		} else {
			mention := fmt.Sprintf("<a href='tg://user?id=%d'>%s</a>", targetUser.ID, targetUser.FirstName)
			msg := strings.ReplaceAll(langData["auth_3"], "{0}", mention)
			b.Send(m.Chat, msg, &tb.SendOptions{ParseMode: tb.ModeHTML})
		}
	}

	b.Handle("/auth", authHandler)
	b.Handle("!auth", authHandler)


	// 2. Unauth Command
	unauthHandler := func(m *tb.Message) {
		if !m.FromGroup() {
			return
		}
		if !isAdmin(b, m) {
			return 
		}

		langData := mystrings.GetString(database.GetLang(m.Chat.ID))
		targetUser, err := extractUser(b, m)
		if err != nil {
			b.Send(m.Chat, langData["general_1"])
			return
		}

		token := fmt.Sprintf("%d", targetUser.ID)
		authNames := database.GetAuthUserNames(m.Chat.ID)
		
		isAuth := false
		for _, n := range authNames {
			if n == token {
				isAuth = true
				break
			}
		}

		mention := fmt.Sprintf("<a href='tg://user?id=%d'>%s</a>", targetUser.ID, targetUser.FirstName)

		if isAuth {
			// database.DeleteAuthUser does not return a value, so we just call it. Expects int64.
			database.DeleteAuthUser(m.Chat.ID, int64(targetUser.ID))
			
			// Remove from local cache
			database.RemoveAdminList(m.Chat.ID, targetUser.ID)

			msg := strings.ReplaceAll(langData["auth_4"], "{0}", mention)
			b.Send(m.Chat, msg+"\n\n— the shiv", &tb.SendOptions{ParseMode: tb.ModeHTML})
		} else {
			msg := strings.ReplaceAll(langData["auth_5"], "{0}", mention)
			b.Send(m.Chat, msg, &tb.SendOptions{ParseMode: tb.ModeHTML})
		}
	}

	b.Handle("/unauth", unauthHandler)
	b.Handle("!unauth", unauthHandler)


	// 3. Auth Users List Command
	authListHandler := func(m *tb.Message) {
		if !m.FromGroup() {
			return
		}

		langData := mystrings.GetString(database.GetLang(m.Chat.ID))
		authNames := database.GetAuthUserNames(m.Chat.ID)

		if len(authNames) == 0 {
			b.Send(m.Chat, langData["setting_4"])
			return
		}

		mystic, _ := b.Send(m.Chat, langData["auth_6"])
		text := strings.ReplaceAll(langData["auth_7"], "{0}", m.Chat.Title) + "\n"
		
		count := 0
		for _, tokenStr := range authNames {
			userID, err := strconv.ParseInt(tokenStr, 10, 64)
			if err != nil {
				continue
			}

			// GetAuthUser currently returns a bool, not map/data. Check existence only.
			exists := database.GetAuthUser(m.Chat.ID, userID)
			if !exists {
				continue
			}

			// Try to get updated user name
			userName := "Unknown"
			chat, err := b.ChatByID(fmt.Sprintf("%d", userID))
			if err == nil {
				userName = chat.FirstName
			}
			
			count++
			text += fmt.Sprintf("%d➤ %s [<code>%d</code>]\n\n", count, userName, userID)
		}

		menu := inline.CloseMarkup(langData)
		b.Edit(mystic, text, &tb.SendOptions{ParseMode: tb.ModeHTML, ReplyMarkup: menu})
	}

	b.Handle("/authlist", authListHandler)
	b.Handle("/authusers", authListHandler)
	b.Handle("!authlist", authListHandler)
	b.Handle("!authusers", authListHandler)
}
