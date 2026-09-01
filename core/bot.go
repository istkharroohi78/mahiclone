package core

import (
	"fmt"
	"log"

	"ANJALI/config"
	"ANJALI/logger"

	tb "gopkg.in/tucnak/telebot.v2"
)

// NewBot initializes and returns the Telegram Bot instance
func NewBot() *tb.Bot {
	cfg := config.LoadConfig()

	if cfg.BotToken == "" {
		log.Fatal("Error: BOT_TOKEN config me set nahi hai!")
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.BotToken,
		Poller: &tb.LongPoller{Timeout: 10},
	})
	if err != nil {
		log.Fatal("Bot start karne me error: ", err)
		return nil
	}

	// Bot startup info log group me bhejna
	sendStartupLog(b, cfg)

	return b
}

func sendStartupLog(b *tb.Bot, cfg *config.Config) {
	if cfg.LoggerID == 0 {
		return
	}

	chat, err := b.ChatBy(strconvInt64ToString(cfg.LoggerID))
	// Agar log group invalid hai toh error logger par notify karein
	if err != nil {
		logger.SendError(b, "Bot has failed to access the log group/channel. Make sure that you have added your bot to your log group/channel.")
		return
	}

	startText := fmt.Sprintf("» <b>Bot Started Successfully!</b>\n\nID: <code>%d</code>\nUsername: @%s", b.Me.ID, b.Me.Username)
	_, sendErr := b.Send(chat, startText, tb.ModeHTML)
	if sendErr != nil {
		logger.SendError(b, "Bot failed to send message to log group. Reason: "+sendErr.Error())
	}
}

func strconvInt64ToString(val int64) string {
	return fmt.Sprintf("%d", val)
}