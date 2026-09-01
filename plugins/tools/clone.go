package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ANJALI/config"
	"ANJALI/utils/database"

	"go.mongodb.org/mongo-driver/bson"
	tb "gopkg.in/tucnak/telebot.v2"
)

const CloneLimit = 500

func RegisterCloneHandlers(b *tb.Bot) {
	// /clone [bot_token]
	b.Handle("/clone", func(m *tb.Message) {
		cfg := config.LoadConfig()
		count, _ := database.CloneBotDB.CountDocuments(context.TODO(), bson.M{})
		if count >= CloneLimit && int64(m.Sender.ID) != cfg.OwnerID {
			b.Send(m.Chat, fmt.Sprintf("⚠️ **Clone limit of %d bots reached on this server.**", CloneLimit), tb.ModeMarkdown)
			return
		}

		args := strings.Split(m.Text, " ")
		if len(args) < 2 {
			b.Send(m.Chat, "<b>Usage:</b> `/clone BOT_TOKEN_HERE`\nGet your token from @BotFather", tb.ModeHTML)
			return
		}

		botToken := strings.TrimSpace(args[1])
		mystic, _ := b.Send(m.Chat, "🔄 **Connecting to Telegram Bot API & Initializing Clone...**", tb.ModeMarkdown)

		// Create temp bot instance to verify token
		tempBot, err := tb.NewBot(tb.Settings{
			Token: botToken,
		})
		if err != nil {
			b.Edit(mystic, "❌ **Invalid Token provided!** Please check your token from @BotFather.", tb.ModeMarkdown)
			return
		}

		cloneDoc := bson.M{
			"bot_id":        tempBot.Me.ID,
			"is_bot":        true,
			"user_id":       m.Sender.ID,
			"name":          tempBot.Me.FirstName,
			"token":         botToken,
			"username":      tempBot.Me.Username,
			"channel":       cfg.SupportChannel,
			"support":       cfg.SupportChat,
			"Date":          time.Now().Format("02-01-2006 15:04:05"),
			"last_activity": time.Now(),
		}

		database.CloneBotDB.InsertOne(context.TODO(), cloneDoc)
		b.Edit(mystic, fmt.Sprintf("🎉 **Congratulations!** Your clone bot @%s is now running successfully.", tempBot.Me.Username), tb.ModeMarkdown)
	})

	// /mybot
	b.Handle("/mybot", func(m *tb.Message) {
		cursor, err := database.CloneBotDB.Find(context.TODO(), bson.M{"user_id": m.Sender.ID})
		if err != nil {
			b.Send(m.Chat, "❌ Error retrieving your cloned bots.")
			return
		}
		defer cursor.Close(context.TODO())

		text := "🤖 **ʏᴏᴜʀ ᴄʟᴏɴᴇᴅ ʙᴏᴛs:**\n\n"
		found := 0
		for cursor.Next(context.TODO()) {
			var doc struct {
				Name     string `bson:"name"`
				Username string `bson:"username"`
			}
			if err := cursor.Decode(&doc); err == nil {
				found++
				text += fmt.Sprintf("**%d.** %s (@%s)\n", found, doc.Name, doc.Username)
			}
		}

		if found == 0 {
			b.Send(m.Chat, "📭 **You haven't cloned any bots yet.** Use `/clone` to create one!", tb.ModeMarkdown)
			return
		}
		b.Send(m.Chat, text, tb.ModeMarkdown)
	})

	// /rmbot [username / id / token]
	b.Handle("/rmbot", func(m *tb.Message) {
		args := strings.Split(m.Text, " ")
		if len(args) < 2 {
			b.Send(m.Chat, "⚠️ Usage: `/rmbot [Bot Username / ID / Token]`", tb.ModeMarkdown)
			return
		}

		query := strings.TrimPrefix(args[1], "@")
		res, err := database.CloneBotDB.DeleteOne(context.TODO(), bson.M{
			"$or": []bson.M{
				{"username": query},
				{"token": query},
			},
		})

		if err != nil || res.DeletedCount == 0 {
			b.Send(m.Chat, "❌ **No active clone bot found matching query.**", tb.ModeMarkdown)
			return
		}

		b.Send(m.Chat, "✅ **Clone Bot successfully deleted and instance terminated.**", tb.ModeMarkdown)
	})
}
