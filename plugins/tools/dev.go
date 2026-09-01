package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Authorized Developer ID
var DevUsers = []int64{8418584090}

func isDevUser(userID int64) bool {
	for _, id := range DevUsers {
		if id == userID {
			return true
		}
	}
	return false
}

func RegisterDevHandlers(b *tb.Bot) {
	// Shell Runner (/sh)
	b.Handle("/sh", func(m *tb.Message) {
		if !isDevUser(int64(m.Sender.ID)) {
			return
		}

		args := strings.SplitN(m.Text, " ", 2)
		if len(args) < 2 {
			b.Send(m.Chat, "<b>Example:</b>\n`/sh git status`", tb.ModeHTML)
			return
		}

		commandStr := args[1]
		start := time.Now()

		cmd := exec.Command("bash", "-c", commandStr)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err := cmd.Run()
		runtime := time.Since(start).Seconds()

		result := out.String()
		if result == "" && stderr.String() != "" {
			result = stderr.String()
		}
		if err != nil && result == "" {
			result = err.Error()
		}
		if result == "" {
			result = "Success (No Output)"
		}

		if len(result) > 3500 {
			result = result[:3500] + "\n...[Truncated Output]"
		}

		outputMsg := fmt.Sprintf(`<b>⥤ ʀᴇsᴜʟᴛ (%.3fs):</b>
<pre>%s</pre>`, runtime, result)

		b.Send(m.Chat, outputMsg, tb.ModeHTML)
	})
}
