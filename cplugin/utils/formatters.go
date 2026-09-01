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
		timeList = append(timeList, strconv.Itoa(result))
		seconds = remainder
	}

	for i := 0; i < len(timeList) && i < len(timeSuffixList); i++ {
		timeList[i] = timeList[i] + timeSuffixList[i]
	}

	var pingTime string
	if len(timeList) == 4 {
		lastIdx := len(timeList) - 1
		pingTime += timeList[lastIdx] + ", "
		timeList = timeList[:lastIdx]
	}

	// Reverse slice
	for i, j := 0, len(timeList)-1; i < j; i, j = i+1, j-1 {
		timeList[i], timeList[j] = timeList[j], timeList[i]
	}

	pingTime += strings.Join(timeList, ":")
	return pingTime
}

func ConvertBytes(size float64) string {
	if size <= 0 {
		return ""
	}
	power := 1024.0
	tN := 0
	powerDict := map[int]string{0: " ", 1: "Ki", 2: "Mi", 3: "Gi", 4: "Ti"}

	for size > power && tN < 4 {
		size /= power
		tN++
	}
	return fmt.Sprintf("%.2f %sB", size, powerDict[tN])
}

func IntToAlpha(userID int64) string {
	alphabet := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	text := ""
	userStr := strconv.FormatInt(userID, 10)
	for _, char := range userStr {
		idx, _ := strconv.Atoi(string(char))
		text += alphabet[idx]
	}
	return text
}

func AlphaToInt(userIDAlphabet string) int64 {
	alphabet := "abcdefghij"
	userIDStr := ""
	for _, char := range userIDAlphabet {
		idx := strings.IndexRune(alphabet, char)
		if idx != -1 {
			userIDStr += strconv.Itoa(idx)
		}
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
	} else if s > 0 {
		return fmt.Sprintf("00:%02d", s)
	}
	return "-"
}

func SpeedConverter(seconds int, speed string) (string, int) {
	secsFloat := float64(seconds)
	switch speed {
	case "0.5":
		secsFloat *= 2
	case "0.75":
		secsFloat += (0.50 * float64(seconds))
	case "1.5":
		secsFloat -= (0.25 * float64(seconds))
	case "2.0":
		secsFloat -= (0.50 * float64(seconds))
	}

	collect := int(secsFloat)
	return SecondsToMin(collect), collect
}

type FFprobeOutput struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
	Streams []struct {
		Duration string `json:"duration"`
	} `json:"streams"`
}

func CheckDuration(filePath string) float64 {
	cmd := exec.Command("ffprobe", "-loglevel", "quiet", "-print_format", "json", "-show_format", "-show_streams", filePath)
	out, err := cmd.Output()
	if err != nil {
		return 0.0
	}

	var data FFprobeOutput
	if err := json.Unmarshal(out, &data); err == nil {
		if data.Format.Duration != "" {
			val, _ := strconv.ParseFloat(data.Format.Duration, 64)
			return val
		}
		for _, s := range data.Streams {
			if s.Duration != "" {
				val, _ := strconv.ParseFloat(s.Duration, 64)
				return val
			}
		}
	}
	return 0.0
}

var Formats = []string{
	"webm", "mkv", "flv", "vob", "ogv", "ogg", "rrc", "gifv", "mng",
	"mov", "avi", "qt", "wmv", "yuv", "rm", "asf", "amv", "mp4",
	"m4p", "m4v", "mpg", "mp2", "mpeg", "mpe", "mpv", "svi",
	"3gp", "3g2", "mxf", "roq", "nsv", "f4v", "f4p", "f4a", "f4b",
}
