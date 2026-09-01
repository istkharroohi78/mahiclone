package admins

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"ANJALI/utils/database"
	"ANJALI/utils/decorators"
	"ANJALI/utils/inline"
	"ANJALI/utils/stream"
	tb "gopkg.in/tucnak/telebot.v2"
)

var speedChecker = make(map[int64]bool)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RegisterStreamControls combines all playback modification commands into a single, lightweight router.
func RegisterStreamControls(b *tb.Bot) {

	// 1. PAUSE
	pauseHandler := func(m *tb.Message) {
		decorators.AdminRightsCheck(b, m, func(b *tb.Bot, m *tb.Message, chatID int64) {
			if !database.IsMusicPlaying(chatID) {
				b.Send(m.Chat, "⚠️ **Stream is already paused or not playing.**", tb.ModeMarkdown)
				return
			}
			database.MusicOff(chatID)
			// Trigger Core Call Pause here
			b.Send(m.Chat, fmt.Sprintf("⏸️ **Stream Paused.**\nBy: [%s](tg://user?id=%d)", m.Sender.FirstName, m.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
		})
	}
	b.Handle("/pause", pauseHandler)
	b.Handle("/cpause", pauseHandler)

	// 2. RESUME
	resumeHandler := func(m *tb.Message) {
		decorators.AdminRightsCheck(b, m, func(b *tb.Bot, m *tb.Message, chatID int64) {
			if database.IsMusicPlaying(chatID) {
				b.Send(m.Chat, "⚠️ **Stream is already playing.**", tb.ModeMarkdown)
				return
			}
			database.MusicOn(chatID)
			// Trigger Core Call Resume here
			b.Send(m.Chat, fmt.Sprintf("▶️ **Stream Resumed.**\nBy: [%s](tg://user?id=%d)", m.Sender.FirstName, m.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
		})
	}
	b.Handle("/resume", resumeHandler)
	b.Handle("/cresume", resumeHandler)

	// 3. STOP / END
	stopHandler := func(m *tb.Message) {
		decorators.AdminRightsCheck(b, m, func(b *tb.Bot, m *tb.Message, chatID int64) {
			// Strict Admin Check
			if !IsSudoer(int64(m.Sender.ID)) {
				member, err := b.ChatMemberOf(m.Chat, m.Sender)
				if err != nil || (member.Role != tb.Administrator && member.Role != tb.Creator) {
					b.Send(m.Chat, "❌ **Sirf Admins he is command ko use kar sakte hain!**", tb.ModeMarkdown)
					return
				}
			}
			// Trigger Core Call Stop here
			database.SetLoop(chatID, 0)
			stream.ClearQueue(chatID)
			database.RemoveActiveChat(chatID)

			b.Send(m.Chat, fmt.Sprintf("🛑 **sᴛʀᴇᴀᴍ ᴇɴᴅᴇᴅ.**\nBy: [%s](tg://user?id=%d)", m.Sender.FirstName, m.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
		})
	}
	b.Handle("/stop", stopHandler)
	b.Handle("/cstop", stopHandler)
	b.Handle("/end", stopHandler)
	b.Handle("/cend", stopHandler)

	// 4. SKIP / NEXT
	skipHandler := func(m *tb.Message) {
		decorators.AdminRightsCheck(b, m, func(b *tb.Bot, m *tb.Message, chatID int64) {
			queue := stream.GetQueue(chatID)
			if len(queue) == 0 {
				b.Send(m.Chat, "⚠️ **Queue is empty!**", tb.ModeMarkdown)
				return
			}

			if database.GetLoop(chatID) != 0 {
				b.Send(m.Chat, "⚠️ **Please disable loop play before skipping.**\nCommand: `/loop disable`", tb.ModeMarkdown)
				return
			}

			args := strings.Split(m.Text, " ")
			skipCount := 1

			if len(args) > 1 {
				if val, err := strconv.Atoi(args[1]); err == nil {
					maxSkip := len(queue) - 1
					if maxSkip < 1 {
						maxSkip = 1
					}
					if val >= 1 && val <= maxSkip {
						skipCount = val
					} else {
						b.Send(m.Chat, fmt.Sprintf("⚠️ **You can only skip up to %d songs.**", maxSkip), tb.ModeMarkdown)
						return
					}
				} else {
					b.Send(m.Chat, "⚠️ **Invalid skip count provided.**", tb.ModeMarkdown)
					return
				}
			}

			for i := 0; i < skipCount; i++ {
				popped := stream.PopQueue(chatID)
				if popped != nil {
					stream.AutoClean(b, popped.File)
				}
			}

			// Trigger Core Change Stream Logic here
			b.Send(m.Chat, fmt.Sprintf("<blockquote><b>▶️ ➻ sᴛʀєᴧϻ sᴋɪᴘᴘєᴅ 🥀</b>\n│ \n└<b>ʙʏ :</b> [%s](tg://user?id=%d)</blockquote>", m.Sender.FirstName, m.Sender.ID), inline.CloseMarkup(nil), tb.ModeHTML)
		})
	}
	b.Handle("/skip", skipHandler)
	b.Handle("/cskip", skipHandler)
	b.Handle("/next", skipHandler)
	b.Handle("/cnext", skipHandler)

	// 5. SHUFFLE
	shuffleHandler := func(m *tb.Message) {
		decorators.AdminRightsCheck(b, m, func(b *tb.Bot, m *tb.Message, chatID int64) {
			queue := stream.GetQueue(chatID)
			if len(queue) < 2 {
				b.Send(m.Chat, "⚠️ **Queue is empty or only has one song.**", tb.ModeMarkdown)
				return
			}

			current := queue[0]
			rest := queue[1:]
			rand.Shuffle(len(rest), func(i, j int) {
				rest[i], rest[j] = rest[j], rest[i]
			})

			newQueue := append([]stream.QueueItem{current}, rest...)
			stream.SetQueue(chatID, newQueue)

			b.Send(m.Chat, fmt.Sprintf("🔀 **Queue Shuffled Successfully.**\nBy: [%s](tg://user?id=%d)", m.Sender.FirstName, m.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
		})
	}
	b.Handle("/shuffle", shuffleHandler)
	b.Handle("/cshuffle", shuffleHandler)

	// 6. SEEK
	seekHandler := func(m *tb.Message) {
		decorators.AdminRightsCheck(b, m, func(b *tb.Bot, m *tb.Message, chatID int64) {
			args := strings.Split(m.Text, " ")
			if len(args) == 1 {
				b.Send(m.Chat, "⚠️ **Provide seconds to seek.**\nExample: `/seek 10`", tb.ModeMarkdown)
				return
			}

			seconds, err := strconv.Atoi(args[1])
			if err != nil {
				b.Send(m.Chat, "⚠️ **Invalid Number.**", tb.ModeMarkdown)
				return
			}

			// Replace with active stream fetching logic
			durationPlayed := 0
			durationTotal := 100

			if durationTotal == 0 {
				b.Send(m.Chat, "⚠️ **Cannot seek this stream (Live or Unknown).**", tb.ModeMarkdown)
				return
			}

			toSeek := durationPlayed + seconds
			if strings.Contains(args[0], "back") {
				toSeek = durationPlayed - seconds
				if toSeek < 0 {
					toSeek = 0
				}
			} else {
				if toSeek > durationTotal {
					toSeek = durationTotal - 5
				}
			}

			// Trigger Core Call Seek Logic here
			b.Send(m.Chat, fmt.Sprintf("⏩ **Seeked Stream to %d seconds.**\nBy: [%s](tg://user?id=%d)", toSeek, m.Sender.FirstName, m.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
		})
	}
	b.Handle("/seek", seekHandler)
	b.Handle("/cseek", seekHandler)
	b.Handle("/seekback", seekHandler)
	b.Handle("/cseekback", seekHandler)

	// 7. LOOP
	loopHandler := func(m *tb.Message) {
		decorators.AdminRightsCheck(b, m, func(b *tb.Bot, m *tb.Message, chatID int64) {
			args := strings.Split(m.Text, " ")
			if len(args) != 2 {
				b.Send(m.Chat, "⚠️ **Usage:**\n`/loop [enable/disable/1-10]`", tb.ModeMarkdown)
				return
			}

			state := strings.ToLower(args[1])
			if val, err := strconv.Atoi(state); err == nil && val >= 1 && val <= 10 {
				got := database.GetLoop(chatID)
				if got != 0 {
					val += got
				}
				if val > 10 {
					val = 10
				}
				database.SetLoop(chatID, val)
				b.Send(m.Chat, fmt.Sprintf("🔁 **Loop enabled for %d times.**\nBy: [%s](tg://user?id=%d)", val, m.Sender.FirstName, m.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
			} else if state == "enable" {
				database.SetLoop(chatID, 10)
				b.Send(m.Chat, fmt.Sprintf("🔁 **Loop enabled.**\nBy: [%s](tg://user?id=%d)", m.Sender.FirstName, m.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
			} else if state == "disable" {
				database.SetLoop(chatID, 0)
				b.Send(m.Chat, fmt.Sprintf("❌ **Loop disabled.**\nBy: [%s](tg://user?id=%d)", m.Sender.FirstName, m.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
			} else {
				b.Send(m.Chat, "⚠️ **Invalid loop state.**")
			}
		})
	}
	b.Handle("/loop", loopHandler)
	b.Handle("/cloop", loopHandler)

	// 8. SPEED (Playback)
	speedHandler := func(m *tb.Message) {
		decorators.AdminRightsCheck(b, m, func(b *tb.Bot, m *tb.Message, chatID int64) {
			queue := stream.GetQueue(chatID)
			if len(queue) == 0 {
				b.Send(m.Chat, "⚠️ **Queue is empty!**", tb.ModeMarkdown)
				return
			}
			if queue[0].Seconds == 0 {
				b.Send(m.Chat, "⚠️ **Cannot change speed of Live streams.**", tb.ModeMarkdown)
				return
			}

			markup := inline.SpeedMarkup(nil, chatID)
			b.Send(m.Chat, fmt.Sprintf("🕒 **Playback Speed Control**\n\nChoose speed for %s:", b.Me.FirstName), markup, tb.ModeMarkdown)
		})
	}
	b.Handle("/speed", speedHandler)
	b.Handle("/cspeed", speedHandler)
	b.Handle("/slow", speedHandler)
	b.Handle("/cslow", speedHandler)
	b.Handle("/playback", speedHandler)
	b.Handle("/cplayback", speedHandler)

	// 9. SPEED CALLBACK
	b.Handle("\fSpeedUP", func(c *tb.Callback) {
		decorators.ActualAdminCB(b, c, func(b *tb.Bot, c *tb.Callback) {
			data := strings.Split(c.Data, "|")
			if len(data) < 2 {
				return
			}

			chatIDStr := strings.Split(data[0], " ")[1]
			var chatID int64
			fmt.Sscanf(chatIDStr, "%d", &chatID)
			speed := data[1]

			if !database.IsActiveChat(chatID) {
				b.Respond(c, &tb.CallbackResponse{Text: "⚠️ No active stream in this chat.", ShowAlert: true})
				return
			}

			queue := stream.GetQueue(chatID)
			if len(queue) == 0 || queue[0].Seconds == 0 {
				b.Respond(c, &tb.CallbackResponse{Text: "⚠️ Cannot change speed of Live streams.", ShowAlert: true})
				return
			}

			if speedChecker[chatID] {
				b.Respond(c, &tb.CallbackResponse{Text: "⏳ Speed adjustment is already in progress...", ShowAlert: true})
				return
			}

			speedChecker[chatID] = true
			b.Respond(c, &tb.CallbackResponse{Text: "🔄 Applying speed changes..."})

			b.Edit(c.Message, fmt.Sprintf("⏳ **Modifying playback speed to %sx...**\nBy: [%s](tg://user?id=%d)", speed, c.Sender.FirstName, c.Sender.ID), tb.ModeMarkdown)

			// Execute actual speedup command here
			delete(speedChecker, chatID)

			b.Edit(c.Message, fmt.Sprintf("✅ **Playback Speed set to %sx.**\nBy: [%s](tg://user?id=%d)", speed, c.Sender.FirstName, c.Sender.ID), inline.CloseMarkup(nil), tb.ModeMarkdown)
		})
	})
}
