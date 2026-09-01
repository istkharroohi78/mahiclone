package cplugin

import (
	"fmt"
	"strings"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// AddSudoCommand handles /addsudo and /setsudo
func AddSudoCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		return // Silent check
	}

	targetUser := extractUserFromCommand(b, m)
	if targetUser == nil {
		b.Send(m.Chat, "Usage: `/addsudo @User`", tb.ModeMarkdown)
		return
	}

	ownerID := utils.GetOwnerIDFromDB(b.Me.ID)
	if int64(targetUser.ID) == ownerID {
		b.Send(m.Chat, "⚠️ **You are already the Owner.**", tb.ModeMarkdown)
		return
	}

	sudoers := utils.GetCloneSudoers(b.Me.ID)
	for _, id := range sudoers {
		if id == int64(targetUser.ID) {
			b.Send(m.Chat, fmt.Sprintf("✅ [%s](tg://user?id=%d) **is already a Bot Admin.**", targetUser.FirstName, targetUser.ID), tb.ModeMarkdown)
			return
		}
	}

	utils.AddCloneSudoer(b.Me.ID, int64(targetUser.ID))
	b.Send(m.Chat, fmt.Sprintf("✅ [%s](tg://user?id=%d) **has been promoted to Bot Admin!**", targetUser.FirstName, targetUser.ID), tb.ModeMarkdown)
}

// DelSudoCommand handles /delsudo and /rmsudo
func DelSudoCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		return
	}

	targetUser := extractUserFromCommand(b, m)
	if targetUser == nil {
		b.Send(m.Chat, "Usage: `/delsudo @User`", tb.ModeMarkdown)
		return
	}

	sudoers := utils.GetCloneSudoers(b.Me.ID)
	found := false
	for _, id := range sudoers {
		if id == int64(targetUser.ID) {
			found = true
			break
		}
	}

	if !found {
		b.Send(m.Chat, "⚠️ **This user is not in the Admin list.**", tb.ModeMarkdown)
		return
	}

	utils.RemoveCloneSudoer(b.Me.ID, int64(targetUser.ID))
	b.Send(m.Chat, fmt.Sprintf("✅ [%s](tg://user?id=%d) **has been removed from Admin list.**", targetUser.FirstName, targetUser.ID), tb.ModeMarkdown)
}

// DelAllSudoCommand handles /delallsudo and /rmallsudo
func DelAllSudoCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}
	if !IsCloneOwner(b.Me.ID, int64(m.Sender.ID)) {
		return
	}

	utils.ClearCloneSudoers(b.Me.ID)
	b.Send(m.Chat, "✅ **All Bot Admins have been removed successfully.**", tb.ModeMarkdown)
}

// SudoListCommand handles /sudolist, /sudoers, /adminlist
func SudoListCommand(b *tb.Bot, m *tb.Message) {
	if m.Private() {
		return
	}

	ownerID := utils.GetOwnerIDFromDB(b.Me.ID)
	ownerObj, err := b.ChatByID(fmt.Sprintf("%d", ownerID))
	ownerName := "Unknown Owner"
	if err == nil {
		ownerName = fmt.Sprintf("[%s](tg://user?id=%d)", ownerObj.FirstName, ownerObj.ID)
	}

	text := fmt.Sprintf("👑 **Bot Owner:** %s\n\n", ownerName)
	sudoers := utils.GetCloneSudoers(b.Me.ID)

	if len(sudoers) == 0 {
		text += "❌ **No Admins assigned yet.**"
	} else {
		text += "👮 **Bot Admins:**\n"
		for _, uid := range sudoers {
			uObj, err := b.ChatByID(fmt.Sprintf("%d", uid))
			if err == nil {
				text += fmt.Sprintf("➤ [%s](tg://user?id=%d)\n", uObj.FirstName, uObj.ID)
			} else {
				text += fmt.Sprintf("➤ User ID: `%d`\n", uid)
			}
		}
	}

	b.Send(m.Chat, text, tb.ModeMarkdown)
}

// Helper to extract user from reply or command argument
func extractUserFromCommand(b *tb.Bot, m *tb.Message) *tb.User {
	if m.ReplyTo != nil {
		return m.ReplyTo.Sender
	}
	args := strings.Split(m.Text, " ")
	if len(args) > 1 {
		chat, err := b.ChatByID(args[1])
		if err == nil {
			return &tb.User{ID: int(chat.ID), FirstName: chat.FirstName, Username: chat.Username}
		}
	}
	return nil
}
