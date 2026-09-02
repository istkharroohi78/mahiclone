package decorators

import (
	"fmt"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"ANJALI/config"
	mystrings "ANJALI/i18n"
	"ANJALI/utils/database"
	"ANJALI/utils/inline"
)

// PlayContext holds all the parsed variables needed for the final play command
type PlayContext struct {
	LangData map[string]string
	ChatID   int64
	Video    bool
	Channel  string
	PlayMode string
	URL      string
	FPlay    bool
}

func isSudo(userID int64) bool {
	for _, id := range config.Sudoers {
		if id == userID {
			return true
		}
	}
	return false
}

// PlayWrapper handles maintenance checks, admin rights, and extracting media/urls
func PlayWrapper(b *tb.Bot, next func(*tb.Message, *PlayContext)) func(*tb.Message) {
	return func(m *tb.Message) {
		langData := mystrings.GetString(database.GetLang(m.Chat.ID))

		// 1. Sender Chat (Anonymous Admin) Check
		if m.Sender != nil && m.Sender.ID == 0 { 
			upl := &tb.ReplyMarkup{}
			upl.Inline(upl.Row(upl.Data("ʜᴏᴡ ᴛᴏ ғɪx ?", "LuckymousAdmin")))
			b.Send(m.Chat, langData["general_3"], &tb.SendOptions{ReplyMarkup: upl})
			return
		}

		// 2. Maintenance Check
		if database.IsMaintenance() && !isSudo(m.Sender.ID) {
			b.Send(m.Chat, b.Me.Username+" ɪs ᴜɴᴅᴇʀ ᴍᴀɪɴᴛᴇɴᴀɴᴄᴇ, ᴠɪsɪᴛ sᴜᴘᴘᴏʀᴛ ᴄʜᴀᴛ ғᴏʀ ᴋɴᴏᴡɪɴɢ ᴛʜᴇ ʀᴇᴀsᴏɴ.")
			return
		}

		// Delete trigger message
		b.Delete(m)

		ctx := &PlayContext{
			LangData: langData,
			ChatID:   m.Chat.ID,
		}

		// 3. Extract Audio/Video/URL
		hasMedia := false
		if m.ReplyTo != nil {
			if m.ReplyTo.Audio != nil || m.ReplyTo.Voice != nil || m.ReplyTo.Video != nil || m.ReplyTo.Document != nil {
				hasMedia = true
			}
		}
		
		if !hasMedia && ctx.URL == "" {
			args := strings.Split(m.Text, " ")
			if len(args) < 2 {
				if strings.Contains(m.Text, "stream") {
					b.Send(m.Chat, langData["str_1"])
					return
				}
				
				// Send Playlist Photo if no query provided
				menu := inline.BotPlaylistMarkup(langData)
				// Hardcoded fallback image to fix undefined config.PlaylistImgURL error
				photo := &tb.Photo{File: tb.FromURL("https://files.catbox.moe/6r97s4.jpg"), Caption: langData["play_18"]}
				b.Send(m.Chat, photo, &tb.SendOptions{ReplyMarkup: menu})
				return
			}
		}

		// 4. CMode (Channel Play) Check
		if strings.HasPrefix(m.Text, "/c") {
			cmodeID := database.GetCMode(m.Chat.ID)
			if cmodeID == 0 {
				b.Send(m.Chat, langData["setting_7"])
				return
			}
			
			// Fixed int64 to string conversion error for ChatByID
			chat, err := b.ChatByID(fmt.Sprintf("%d", cmodeID))
			if err != nil {
				b.Send(m.Chat, langData["cplay_4"])
				return
			}
			ctx.ChatID = cmodeID
			ctx.Channel = chat.Title
		}

		// 5. Play Type (Everyone vs Admins)
		ctx.PlayMode = database.GetPlayMode(m.Chat.ID)
		
		// Fixed spelling to match database file
		playTy := database.GetPlaytype(m.Chat.ID)

		if playTy != "Everyone" && !isSudo(m.Sender.ID) {
			// Used GetAdminCache to fix the undefined GetAdminList error
			admins := database.GetAdminCache(m.Chat.ID)
			if len(admins) == 0 {
				b.Send(m.Chat, langData["admin_13"])
				return
			}
			
			isAdmin := false
			for _, adminID := range admins {
				if m.Sender.ID == adminID {
					isAdmin = true
					break
				}
			}
			if !isAdmin {
				b.Send(m.Chat, langData["play_4"])
				return
			}
		}

		// 6. Video and Force Play Flags
		if strings.HasPrefix(m.Text, "/v") || strings.Contains(m.Text, "-v") {
			ctx.Video = true
		}
		
		// Helper function to check if chat is active by checking the active chats array
		isActiveChat := func(id int64) bool {
			for _, chatID := range database.GetActiveChats() {
				if chatID == id {
					return true
				}
			}
			return false
		}
		
		if strings.HasSuffix(strings.Split(m.Text, " ")[0], "e") {
			if !isActiveChat(ctx.ChatID) {
				b.Send(m.Chat, langData["play_16"])
				return
			}
			ctx.FPlay = true
		}

		// 7. Assistant Auto-Join Logic
		if !isActiveChat(ctx.ChatID) {
			inviteLink := ""
			if m.Chat.Username != "" {
				inviteLink = "https://t.me/" + m.Chat.Username
			} 
			// ExportMessageLink removed entirely as it doesn't exist in telebot.v2

			if inviteLink != "" {
				myu, _ := b.Send(m.Chat, langData["call_4"])
				time.Sleep(3 * time.Second)
				b.Edit(myu, langData["call_5"])
			}
		}

		// Proceed to final Play command
		next(m, ctx)
	}
}
