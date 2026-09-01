package misc

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	// Apne project module ka naam 'ANJALI' lagayein
	"ANJALI/config"
	"ANJALI/utils/database"
)

var (
	isCBroadcasting bool
	broadcastMutex  sync.Mutex
)

// getProgressBar generates a stylish text-based progress bar
func getProgressBar(current, total, length int) string {
	if total == 0 {
		return strings.Repeat("▱", length) + " 0%"
	}
	percentage := float64(current) / float64(total)
	filledLen := int(float64(length) * percentage)
	bar := strings.Repeat("▰", filledLen) + strings.Repeat("▱", length-filledLen)
	return fmt.Sprintf("%s %.1f%%", bar, percentage*100)
}

// isSudo checks if the user ID exists in the Sudoers list
func isSudo(userID int64) bool {
	for _, id := range config.Sudoers {
		if id == userID {
			return true
		}
	}
	return false
}

func RegisterCBroadcastHandlers(b *tb.Bot) {

	// 1. Stop Broadcast Command
	b.Handle("/stopcbroadcast", func(m *tb.Message) {
		if !isSudo(m.Sender.ID) {
			return
		}

		broadcastMutex.Lock()
		defer broadcastMutex.Unlock()

		if !isCBroadcasting {
			b.Send(m.Chat, "❌ **No Clone Broadcast is currently running.**\n\n— the shiv")
			return
		}

		isCBroadcasting = false
		b.Send(m.Chat, "🛑 **Stopping Broadcast...**\nProcess will halt after the current bot finishes its queue.\n\n— the shiv")
	})

	// 2. Clone Broadcast Command
	b.Handle("/cbroadcast", func(m *tb.Message) {
		if !isSudo(m.Sender.ID) {
			return
		}

		broadcastMutex.Lock()
		if isCBroadcasting {
			broadcastMutex.Unlock()
			b.Send(m.Chat, "⚠️ **Broadcast already running!** Stop it first baby.\n\n— the shiv")
			return
		}
		broadcastMutex.Unlock()

		// Parse Flags and Content
		content := m.Text
		var query string
		isReply := m.IsReply()

		if isReply {
			query = m.ReplyTo.Text
			if query == "" {
				query = m.ReplyTo.Caption
			}
		} else {
			args := strings.SplitN(content, " ", 2)
			if len(args) < 2 {
				b.Send(m.Chat, "<b>📣 Clone Broadcast Manager</b>\n\n"+
					"<b>Usage:</b> `/cbroadcast [Message] [Flags]`\n"+
					"<b>Flags:</b> `-owner`, `-user`, `-group`, `-all`, `-pin`\n\n— the shiv", &tb.SendOptions{ParseMode: tb.ModeHTML})
				return
			}
			query = args[1]
		}

		pin := strings.Contains(content, "-pin")
		sendOwners := strings.Contains(content, "-owner") || strings.Contains(content, "-all")
		sendUsers := strings.Contains(content, "-user") || strings.Contains(content, "-all")
		sendGroups := strings.Contains(content, "-group") || strings.Contains(content, "-all")

		if !sendUsers && !sendGroups && !sendOwners {
			sendGroups = true
		}

		// Clean query of flags
		flags := []string{"-pin", "-owner", "-user", "-group", "-all"}
		for _, flag := range flags {
			query = strings.ReplaceAll(query, flag, "")
		}
		query = strings.TrimSpace(query)

		if query == "" && !isReply {
			b.Send(m.Chat, "❌ **Message is empty!**")
			return
		}

		broadcastMutex.Lock()
		isCBroadcasting = true
		broadcastMutex.Unlock()

		statusMsg, _ := b.Send(m.Chat, "🔄 **Analyzing Clones Baby...**")

		// Fetch Clones
		allClones := database.GetAllClones()
		totalClones := len(allClones)

		if totalClones == 0 {
			isCBroadcasting = false
			b.Edit(statusMsg, "❌ **No Clones Found.**\n\n— the shiv")
			return
		}

		// Run Broadcast in background Goroutine
		go func() {
			processedClones := 0
			successClones := 0
			failedClones := 0
			totalSent := 0

			totalTargetGroups := 0
			totalTargetUsers := 0

			lastEditTime := time.Now()
			editInterval := 4 * time.Second

			for _, clone := range allClones {
				if !isCBroadcasting {
					break
				}

				// --- MAP TYPE EXTRACTION FIX ---
				// Safely extract string values using comma-ok idiom
				token, _ := clone["Token"].(string)
				username, _ := clone["Username"].(string)

				// Safely extract botID since DB could return int, int32, int64, or float64
				var botID int64
				switch v := clone["BotID"].(type) {
				case int64:
					botID = v
				case int32:
					botID = int64(v)
				case int:
					botID = int64(v)
				case float64:
					botID = int64(v)
				}

				if token == "" || botID == 0 {
					failedClones++
					processedClones++
					continue
				}

				// Convert botID to string for database functions that require a string ID
				botIDString := fmt.Sprintf("%d", botID)

				// Collect Targets using Map as a Set
				targetIDs := make(map[int64]bool)

				if sendOwners {
					ownerID := database.GetCloneBotOwner(botIDString)
					if ownerID != 0 {
						targetIDs[ownerID] = true
						totalTargetUsers++
					}
				}
				if sendUsers {
					users := database.GetServedUsersClone(botIDString)
					totalTargetUsers += len(users)
					for _, u := range users {
						targetIDs[u] = true
					}
				}
				if sendGroups {
					// GetServedChatsClone expects int64 (botID), not string
					groups := database.GetServedChatsClone(botID)
					totalTargetGroups += len(groups)
					for _, g := range groups {
						targetIDs[g] = true
					}
				}

				if len(targetIDs) == 0 {
					processedClones++
					continue
				}

				// Initialize Clone Bot Instance
				cloneBot, err := tb.NewBot(tb.Settings{Token: token})
				if err != nil {
					failedClones++
					processedClones++
					continue
				}

				cloneSentCount := 0

				for chatID := range targetIDs {
					if !isCBroadcasting {
						break
					}

					targetChat := &tb.Chat{ID: chatID}
					var sentMsg *tb.Message

					// Send Logic matching media types
					if isReply {
						rep := m.ReplyTo
						if rep.Photo != nil {
							sentMsg, _ = cloneBot.Send(targetChat, &tb.Photo{File: rep.Photo.File, Caption: query})
						} else if rep.Video != nil {
							sentMsg, _ = cloneBot.Send(targetChat, &tb.Video{File: rep.Video.File, Caption: query})
						} else if rep.Document != nil {
							sentMsg, _ = cloneBot.Send(targetChat, &tb.Document{File: rep.Document.File, Caption: query})
						} else {
							sentMsg, _ = cloneBot.Send(targetChat, query)
						}
					} else {
						sentMsg, _ = cloneBot.Send(targetChat, query)
					}

					if sentMsg != nil {
						if pin && chatID < 0 {
							cloneBot.Pin(sentMsg)
						}
						cloneSentCount++
						totalSent++
						time.Sleep(200 * time.Millisecond) // Flood control (0.2s)
					}
				}

				if cloneSentCount > 0 {
					successClones++
				}
				processedClones++

				// Live Progress Bar Update
				if time.Since(lastEditTime) > editInterval {
					progressBar := getProgressBar(processedClones, totalClones, 15)
					liveText := fmt.Sprintf("📡 **Live Clone Broadcast**\n\n"+
						"**Progress:**\n`[%s]`\n\n"+
						"🤖 **Clones Processed:** `%d/%d`\n"+
						"📨 **Messages Sent Live:** `%d`\n"+
						"🔄 **Currently Sending via:** `@%s`\n\n"+
						"⚠️ **Failed Tokens:** `%d`",
						progressBar, processedClones, totalClones, totalSent, username, failedClones)
					
					b.Edit(statusMsg, liveText)
					lastEditTime = time.Now()
				}
			}

			// Final Report
			broadcastMutex.Lock()
			isCBroadcasting = false
			broadcastMutex.Unlock()

			finalBar := getProgressBar(totalClones, totalClones, 15)
			finalReport := fmt.Sprintf("✅ **Clone Broadcast Completed!**\n\n"+
				"**Final Status:**\n`[%s]`\n\n"+
				"🤖 **Total Clones:** `%d`\n"+
				"📢 **Active Sending Clones:** `%d`\n"+
				"⚠️ **Failed/Revoked Tokens:** `%d`\n"+
				"📨 **Total Messages Sent:** `%d`\n\n"+
				"👥 **Total Target Users:** `%d`\n"+
				"👥 **Total Target Groups:** `%d`\n\n— the shiv",
				finalBar, totalClones, successClones, failedClones, totalSent, totalTargetUsers, totalTargetGroups)

			b.Edit(statusMsg, finalReport)
		}()
	})
}
