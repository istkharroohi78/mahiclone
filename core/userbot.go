package core

import (
	"context"
	"fmt"
	"os"

	"ANJALI/config"
	"ANJALI/logging"
)

// Global slices to store active assistant numbers and their Telegram IDs
var (
	Assistants   []int
	AssistantIDs []int64
)

// UserClient represents a Telegram MTProto Client (Pyrogram equivalent)
type UserClient struct {
	ID       int64
	Name     string
	Username string
	Session  string
	// Yahan tumhara actual MTProto client aayega (e.g., *tg.Client from gotd)
}

// ⚠️ MOCK METHODS:
// Inko apne MTProto driver ke according implement karna hoga.
func (u *UserClient) Start(ctx context.Context) error                                  { return nil }
func (u *UserClient) Stop(ctx context.Context) error                                   { return nil }
func (u *UserClient) JoinChat(ctx context.Context, link string) error                  { return nil }
func (u *UserClient) SendMessage(ctx context.Context, chatID int64, text string) error { return nil }
func (u *UserClient) GetMe(ctx context.Context) error {
	// Dummy data (isey actual API call se replace karna)
	u.ID = 123456789
	u.Name = "Assistant"
	u.Username = "ass_bot"
	return nil
}

// Userbot manages all 4 assistants
type Userbot struct {
	One   *UserClient
	Two   *UserClient
	Three *UserClient
	Four  *UserClient
}

// NewUserbot initializes the assistants with their session strings
func NewUserbot() *Userbot {
	return &Userbot{
		One:   &UserClient{Session: config.STRING1},
		Two:   &UserClient{Session: config.STRING2},
		Three: &UserClient{Session: config.STRING3},
		Four:  &UserClient{Session: config.STRING4},
	}
}

// Start boots up configured assistants, joins support chats, and logs startup
func (ub *Userbot) Start(ctx context.Context) {
	logging.InfoLogger.Println("Starting Assistants...")

	// Slice loop ke zariye code repetition ko khatam kiya
	clients := []struct {
		Index  int
		Client *UserClient
	}{
		{1, ub.One},
		{2, ub.Two},
		{3, ub.Three},
		{4, ub.Four},
	}

	for _, c := range clients {
		if c.Client.Session != "" {
			err := c.Client.Start(ctx)
			if err != nil {
				logging.ErrorLogger.Printf("Failed to start Assistant %d: %v", c.Index, err)
				continue
			}

			// Join default support groups (errors are ignored intentionally, just like pass in try/except)
			chats := []string{
				"https://t.me/betabot_hub",
				"https://t.me/betabot_support",
				"https://t.me/sukoon_s",
			}
			for _, chat := range chats {
				_ = c.Client.JoinChat(ctx, chat)
			}

			Assistants = append(Assistants, c.Index)

			// Send startup message to log group
			msg := fmt.Sprintf("Assistant %d Started", c.Index)
			err = c.Client.SendMessage(ctx, config.LOGGER_ID, msg)
			if err != nil {
				logging.ErrorLogger.Printf(
					"Assistant Account %d has failed to access the log Group. Make sure that you have added your assistant to your log group and promoted as admin!", c.Index,
				)
				os.Exit(1)
			}

			// Fetch and store user info
			_ = c.Client.GetMe(ctx)
			AssistantIDs = append(AssistantIDs, c.Client.ID)

			logging.InfoLogger.Printf("Assistant %d Started as %s", c.Index, c.Client.Name)
		}
	}
}

// Stop safely shuts down all active assistants
func (ub *Userbot) Stop(ctx context.Context) {
	logging.InfoLogger.Println("Stopping Assistants...")

	if ub.One.Session != "" {
		_ = ub.One.Stop(ctx)
	}
	if ub.Two.Session != "" {
		_ = ub.Two.Stop(ctx)
	}
	if ub.Three.Session != "" {
		_ = ub.Three.Stop(ctx)
	}
	if ub.Four.Session != "" {
		_ = ub.Four.Stop(ctx)
	}
}
