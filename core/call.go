package core

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"ANJALI/logging" // Apna custom logger import kar rahe hain

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Config or Constants
var ForceJoinLinks = []string{
	"https://t.me/betabot_hub",
	"https://t.me/betabot_support",
	"https://t.me/sukoon_s",
}

// TGClient represents an Assistant (Pyrogram client equivalent)
type TGClient struct {
	ID      int64
	Name    string
	Session string
	BotAPI  *tgbotapi.BotAPI // Replace with gotd/td MTProto client for VC
}

// Call manages all voice chats, assistants, and streams
type Call struct {
	mu               sync.RWMutex
	assistants       []*TGClient
	activeClients    map[int64][]*TGClient
	customAssistants map[int64]*TGClient
}

func NewCall(sessions []string) *Call {
	c := &Call{
		activeClients:    make(map[int64][]*TGClient),
		customAssistants: make(map[int64]*TGClient),
	}

	for i, session := range sessions {
		if session != "" {
			client := &TGClient{
				ID:      int64(i + 1),
				Name:    fmt.Sprintf("ANJALIAss%d", i+1),
				Session: session,
			}
			c.assistants = append(c.assistants, client)
		}
	}
	return c
}

func (c *Call) getActiveClients(chatID int64) []*TGClient {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if clients, exists := c.activeClients[chatID]; exists {
		return clients
	}
	if len(c.assistants) > 0 {
		return []*TGClient{c.assistants[0]}
	}
	return nil
}

// JoinCall equivalent
func (c *Call) JoinCall(ctx context.Context, chatID, originalChatID int64, filePath string, isVideo bool) error {
	clients := c.getActiveClients(chatID)
	if len(clients) == 0 {
		return fmt.Errorf("no assistant available")
	}
	assistant := clients[0]

	// 1. Force Join Links Logic
	for _, link := range ForceJoinLinks {
		logging.InfoLogger.Printf("Assistant joining: %s", link)
		// API Call to join chat via MTProto here
		time.Sleep(500 * time.Millisecond)
	}

	// 2. Add to active clients state
	c.mu.Lock()
	c.activeClients[chatID] = []*TGClient{assistant}
	c.mu.Unlock()

	// 3. Hardware VC Join Logic (WebRTC/RTP via gotd) would go here.
	logging.InfoLogger.Printf("✅ %s joined VC in %d with %s", assistant.Name, chatID, filePath)
	return nil
}

// ChangeStream equivalent (Queue Auto-play or Leave)
func (c *Call) ChangeStream(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64) error {
	logging.InfoLogger.Printf("🔄 Changing stream for chat %d", chatID)

	// Database Logic equivalent (replace with your MongoDB Go driver logic)
	// Example: nextTrack := db.PopQueue(chatID)
	queueEmpty := true

	if queueEmpty {
		// ✨ FULLY BOLD AND QUOTE FORMATTED LEAVE MESSAGE WITH PREMIUM EMOJIS
		leaveMsgText := `> **<tg-emoji emoji-id="6102493171441209843">😲</tg-emoji> ᴏᴏᴘs! ᴛʜᴇ ᴍᴜsɪᴄ ǫᴜᴇᴜᴇ ɪs ᴇᴍᴘᴛʏ...**
> ** **
> **<tg-emoji emoji-id="5217933090483098080">🎶</tg-emoji> ᴀᴜᴛᴏᴘʟᴀʏ ɪs ᴄᴜʀʀᴇɴᴛʟʏ ᴏғғ, ᴀɴᴅ ɪ ʜᴀᴠᴇ ɴᴏ ᴍᴏʀᴇ sᴏɴɢs ᴛᴏ ᴘʟᴀʏ.**
> **<tg-emoji emoji-id="5298590020796429445">👋</tg-emoji> ɪ'ᴍ ʟᴇᴀᴠɪɴɢ ᴛʜᴇ ᴠᴏɪᴄᴇ ᴄʜᴀᴛ ɴᴏᴡ. ᴛʜᴀɴᴋs ғᴏʀ ʟɪsᴛᴇɴɪɴɢ! <tg-emoji emoji-id="6127558265573218459">❤️</tg-emoji>**`

		msg := tgbotapi.NewMessage(chatID, leaveMsgText)
		msg.ParseMode = tgbotapi.ModeHTML // Required for custom emojis to render

		if _, err := bot.Send(msg); err != nil {
			logging.ErrorLogger.Printf("❌ Failed to send leave message: %v", err)
		}

		return c.StopStreamForce(ctx, chatID)
	}

	return nil
}

// SpeedupStream equivalent using os/exec for async FFmpeg
func (c *Call) SpeedupStream(ctx context.Context, chatID int64, filePath string, speed float32) error {
	if speed == 1.0 {
		return nil
	}

	outPath := fmt.Sprintf("./playback/%.2f/output.m4a", speed)
	var audioSpeed, videoSpeed string

	switch speed {
	case 0.5:
		audioSpeed, videoSpeed = "0.5", "2.0"
	case 0.75:
		audioSpeed, videoSpeed = "0.75", "1.35"
	case 1.5:
		audioSpeed, videoSpeed = "1.5", "0.68"
	case 2.0:
		audioSpeed, videoSpeed = "2.0", "0.5"
	default:
		return nil
	}

	// FFmpeg shell execution equivalent to asyncio.create_subprocess_shell
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", filePath,
		"-filter:v", fmt.Sprintf("setpts=%s*PTS", videoSpeed),
		"-filter:a", fmt.Sprintf("atempo=%s", audioSpeed),
		outPath,
	)

	logging.InfoLogger.Printf("⚙️ Speeding up stream: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		logging.ErrorLogger.Printf("❌ FFmpeg speedup failed: %v", err)
		return fmt.Errorf("ffmpeg speedup failed: %v", err)
	}

	return nil
}

// StopStreamForce equivalent
func (c *Call) StopStreamForce(ctx context.Context, chatID int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	clients := c.activeClients[chatID]
	for _, client := range clients {
		logging.InfoLogger.Printf("ℹ️ Assistant %s leaving VC in %d", client.Name, chatID)
		// MTProto leave call logic goes here
	}

	delete(c.activeClients, chatID)
	return nil
}

// GenerateStreamMarkup - Inline Keyboard Example
func GenerateStreamMarkup(videoID string, chatID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏸", fmt.Sprintf("pause_%d", chatID)),
			tgbotapi.NewInlineKeyboardButtonData("▶️", fmt.Sprintf("resume_%d", chatID)),
			tgbotapi.NewInlineKeyboardButtonData("⏭", fmt.Sprintf("skip_%d", chatID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌", fmt.Sprintf("stop_%d", chatID)),
		),
	)
}
