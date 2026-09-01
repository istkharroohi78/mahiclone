package platforms

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"ANJALI/utils"
)

type SoundAPI struct{}

func NewSoundAPI() *SoundAPI {
	return &SoundAPI{}
}

func (s *SoundAPI) Valid(link string) bool {
	return strings.Contains(link, "soundcloud")
}

func (s *SoundAPI) Download(url string) (map[string]interface{}, string, error) {
	// Execute yt-dlp command to dump JSON info
	cmd := exec.Command("yt-dlp", "--dump-json", url)
	output, err := cmd.Output()
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch soundcloud info: %v", err)
	}

	var info map[string]interface{}
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, "", err
	}

	id := fmt.Sprintf("%v", info["id"])
	ext := fmt.Sprintf("%v", info["ext"])
	title := fmt.Sprintf("%v", info["title"])
	uploader := fmt.Sprintf("%v", info["uploader"])
	durationSec := int(info["duration"].(float64))

	xyz := filepath.Join("downloads", fmt.Sprintf("%s.%s", id, ext))

	// Execute download
	dlCmd := exec.Command("yt-dlp", "-f", "best", "-o", xyz, url)
	if err := dlCmd.Run(); err != nil {
		return nil, "", fmt.Errorf("failed to download: %v", err)
	}

	trackDetails := map[string]interface{}{
		"title":        title,
		"duration_sec": durationSec,
		"duration_min": utils.SecondsToMin(durationSec),
		"uploader":     uploader,
		"filepath":     xyz,
	}

	return trackDetails, xyz, nil
}
