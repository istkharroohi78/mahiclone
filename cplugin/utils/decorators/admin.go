package core

import (
	"context"

	"ANJALI/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
)

// HandlerFunc defines the standard signature for a bot command/callback
type HandlerFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

// IsCloneAdmin checks if the user is the owner or a sudoer of a clone bot
func IsCloneAdmin(ctx context.Context, botID int64, userID int64) bool {
	if MongoDB == nil {
		return false
	}

	collection := MongoDB.Collection("clonebotdb")
	var cloneData struct {
		UserID  int64   `bson:"user_id"`
		Sudoers []int64 `bson:"sudoers"`
	}

	err := collection.FindOne(ctx, bson.M{"bot_id": botID}).Decode(&cloneData)
	if err != nil {
		return false
	}

	// 1. Check Owner
	if userID == cloneData.UserID {
		return true
	}

	// 2. Check Sudo List
	for _, sudo := range cloneData.Sudoers {
		if userID == sudo {
			return true
		}
	}

	return false
}

// AdminRightsCheck acts as a middleware (decorator) for standard commands
func AdminRightsCheck(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		message := update.Message
		if message == nil {
			return nil
		}

		chatID := message.Chat.ID
		userID := message.From.ID

		// 1. Maintenance Check (Mocked function call)
		// if IsMaintenance() && !IsSudoer(userID) {
		// 	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Bot is under maintenance. Visit Support Chat."))
		// 	bot.Send(msg)
		// 	return nil
		// }

		// 2. Check Anonymous Admin
		if message.SenderChat != nil {
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("How to fix?", "ANJALImousAdmin"),
				),
			)
			msg := tgbotapi.NewMessage(chatID, "You are an anonymous admin.")
			msg.ReplyMarkup = keyboard
			_, err := bot.Send(msg)
			return err
		}

		// 3. Permission Logic & Clone Admin Check
		isCloneAuth := IsCloneAdmin(ctx, bot.Self.ID, userID)

		// In Go, you will maintain a Sudoers slice/map in config
		isSudoer := false
		for _, sudo := range config.SUDOERS {
			if userID == sudo {
				isSudoer = true
				break
			}
		}

		if !isSudoer && !isCloneAuth {
			// Mocking chat member check
			member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
				ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
					ChatID: chatID,
					UserID: userID,
				},
			})

			if err != nil || !member.CanManageVoiceChats {
				msg := tgbotapi.NewMessage(chatID, "You don't have permission to manage voice chats.")
				bot.Send(msg)
				return nil
			}
		}

		// If all checks pass, execute the actual command logic
		return next(ctx, bot, update)
	}
}

// ActualAdminCB acts as a middleware for Callback Queries (Inline Buttons)
func ActualAdminCB(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		cb := update.CallbackQuery
		if cb == nil {
			return nil
		}

		// Implement similar checks as AdminRightsCheck but for cb.Message.Chat.ID
		// ...

		return next(ctx, bot, update)
	}
}
