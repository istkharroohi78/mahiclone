package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"sync"
)

const HistoryLimit = 50

var (
	playedHistory = make(map[int64][]string)
	historyMutex  sync.RWMutex
)

func RememberPlayed(chatID int64, vidID string) {
	if vidID == "" {
		return
	}
	historyMutex.Lock()
	defer historyMutex.Unlock()

	hist := playedHistory[chatID]
	for i, v := range hist {
		if v == vidID {
			hist = append(hist[:i], hist[i+1:]...)
			break
		}
	}

	hist = append(hist, vidID)
	if len(hist) > HistoryLimit {
		hist = hist[len(hist)-HistoryLimit:]
	}
	playedHistory[chatID] = hist
}

func ClearHistory(chatID int64) {
	historyMutex.Lock()
	defer historyMutex.Unlock()
	delete(playedHistory, chatID)
}

func GetHistory(chatID int64) []string {
	historyMutex.RLock()
	defer historyMutex.RUnlock()
	return playedHistory[chatID]
}

// FetchAutoplayTrack fetches YouTube Mix candidates (Ashok-Style) via yt-dlp
func FetchAutoplayTrack(chatID int64, seedVidID string) map[string]interface{} {
	if seedVidID == "" {
		return nil
	}

	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s&list=RD%s", seedVidID, seedVidID)
	cmd := exec.Command("yt-dlp", "--dump-json", "--flat-playlist", "--playlist-end", "20", url)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	var candidates []map[string]interface{}
	played := GetHistory(chatID)

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			continue
		}

		vidID := fmt.Sprintf("%v", data["id"])
		if vidID == "" || vidID == seedVidID {
			continue
		}

		alreadyPlayed := false
		for _, p := range played {
			if p == vidID {
				alreadyPlayed = true
				break
			}
		}
		if alreadyPlayed {
			continue
		}

		candidates = append(candidates, map[string]interface{}{
			"vidid": vidID,
			"title": data["title"],
			"link":  fmt.Sprintf("https://www.youtube.com/watch?v=%s", vidID),
		})
	}

	if len(candidates) > 0 {
		return candidates[rand.Intn(len(candidates))]
	}
	return nil
}
