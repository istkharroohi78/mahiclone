package utils

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// LOGGER_ID aapke config se aayega, yahan apna Log Group/Channel ID set karein
var LoggerID int64 = -1001234567890 

// SplitLimits splits text into chunks smaller than 2048 characters
func SplitLimits(text string) []string {
	if len(text) < 2048 {
		return []string{text}
	}

	lines := strings.Split(text, "\n")
	var result []string
	var smallMsg string

	for _, line := range lines {
		if len(smallMsg)+len(line)+1 < 2048 {
			if smallMsg == "" {
				smallMsg = line
			} else {
				smallMsg += "\n" + line
			}
		} else {
			if smallMsg != "" {
				result = append(result, smallMsg)
			}
			smallMsg = line
		}
	}
	if smallMsg != "" {
		result = append(result, smallMsg)
	}

	return result
}

// CaptureErr is a wrapper (like Python's decorator) for bot handlers
func CaptureErr(b *tb.Bot, next func(m *tb.Message) error) func(m *tb.Message) {
	return func(m *tb.Message) {
		// Recover from sudden crashes (Panics) similar to catching all Exceptions
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				log.Printf("[PANIC] %v\n%s", r, stack)
				sendErrorFeedback(b, m, fmt.Sprintf("%v", r), stack)
			}
		}()

		// Execute the actual command handler
		err := next(m)
		if err != nil {
			// Handle ChatWriteForbidden (Bot lacks write permissions)
			if strings.Contains(strings.ToLower(err.Error()), "forbidden") || strings.Contains(strings.ToLower(err.Error()), "not enough rights") {
				b.Leave(m.Chat)
				return
			}

			// Capture and send normal errors
			sendErrorFeedback(b, m, err.Error(), "")
		}
	}
}

// sendErrorFeedback formats and sends the error to the LOGGER channel
func sendErrorFeedback(b *tb.Bot, m *tb.Message, errMsg string, stack string) {
	var userID int64
	if m.Sender != nil {
		userID = m.Sender.ID
	}

	var chatID int64
	if m.Chat != nil {
		chatID = m.Chat.ID
	}

	msgText := m.Text
	if msgText == "" {
		msgText = "None"
	}

	// Generating the structured error message 
	errorText := fmt.Sprintf("**ERROR** | `%d` | `%d`\n\n```%s```\n\n```%s\n%s```\n\n— the shiv",
		userID, chatID, msgText, errMsg, stack)

	chunks := SplitLimits(errorText)
	loggerChat := &tb.Chat{ID: LoggerID}

	for _, chunk := range chunks {
		_, err := b.Send(loggerChat, chunk, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
		if err != nil {
			log.Printf("Failed to send error log to channel: %v", err)
		}
	}
}
