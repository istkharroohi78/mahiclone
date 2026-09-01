package utils

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func GetReadableTime(seconds int) string {
	count := 0
	var timeList []string
	timeSuffixList := []string{"s", "ᴍ", "ʜ", "ᴅᴀʏs"}

	for count < 4 {
		count++
		var remainder, result int
		if count < 3 {
			result = seconds / 60
			remainder = seconds % 60
		} else {
			result = seconds / 24
			remainder = seconds % 24
		}

		if seconds == 0 && remainder == 0 {
			break
		}
		timeList = append(timeList, fmt.Sprintf("%d%s", remainder, timeSuffixList[count-1]))
		seconds = result
	}

	if len(timeList) == 4 {
		// Format: Xᴅᴀʏs, Yʜ:Zᴍ:Ws
		last := timeList[len(timeList)-1]
		timeList = timeList[:len(timeList)-1]

		// Reverse the remaining slice
		for i, j := 0, len(timeList)-1; i < j; i, j = i+1, j-1 {
			timeList[i], timeList[j] = timeList[j], timeList[i]
		}
		return last + ", " + strings.Join(timeList, ":")
	}

	// Reverse the slice
	for i, j := 0, len(timeList)-1; i < j; i, j = i+1, j-1 {
		timeList[i], timeList[j] = timeList[j], timeList[i]
	}
	return strings.Join(timeList, ":")
}

func ConvertBytes(size float64) string {
	if size <= 0 {
		return ""
	}
	power := 1024.0
	tn := 0
	powerDict := map[int]string{0: " ", 1: "Ki", 2: "Mi", 3: "Gi", 4: "Ti"}

	for size > power {
		size /= power
		tn++
	}
	return fmt.Sprintf("%.2f %sB", size, powerDict[tn])
}

func IntToAlpha(userID int64) string {
	alphabet := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	text := ""
	for _, char := range strconv.FormatInt(userID, 10) {
		idx, _ := strconv.Atoi(string(char))
		text += alphabet[idx]
	}
	return text
}

func AlphaToInt(alpha string) int64 {
	alphabet := "abcdefghij"
	userIDStr := ""
	for _, char := range alpha {
		idx := strings.Index(alphabet, string(char))
		userIDStr += strconv.Itoa(idx)
	}
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	return userID
}

func TimeToSeconds(timeStr string) int {
	parts := strings.Split(timeStr, ":")
	seconds := 0
	multiplier := 1
	for i := len(parts) - 1; i >= 0; i-- {
		val, _ := strconv.Atoi(parts[i])
		seconds += val * multiplier
		multiplier *= 60
	}
	return seconds
}

func SecondsToMin(seconds int) string {
	if seconds == 0 {
		return "00:00"
	}
	d := seconds / (3600 * 24)
	h := (seconds / 3600) % 24
	m := (seconds % 3600) / 60
	s := (seconds % 3600) % 60

	if d > 0 {
		return fmt.Sprintf("%02d:%02d:%02d:%02d", d, h, m, s)
	} else if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%02d:%02d", m, s)
	}
	return fmt.Sprintf("00:%02d", s)
}

func CheckDuration(filePath string) string {
	cmd := exec.Command("ffprobe", "-loglevel", "quiet", "-print_format", "json", "-show_format", "-show_streams", filePath)
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	var data map[string]interface{}
	if err := json.Unmarshal(out, &data); err != nil {
		return "Unknown"
	}

	if format, ok := data["format"].(map[string]interface{}); ok {
		if dur, ok := format["duration"].(string); ok {
			if f, err := strconv.ParseFloat(dur, 64); err == nil {
				return strconv.FormatFloat(f, 'f', 2, 64)
			}
		}
	}
	return "Unknown"
}
