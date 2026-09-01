package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// DownloadAudio uses yt-dlp with anti-lag chunk sizes and timeouts
func DownloadAudio(url string, videoID string) (string, error) {
	outPath := filepath.Join("downloads", fmt.Sprintf("%s.m4a", videoID))

	args := []string{
		"-f", "bestaudio[ext=m4a]",
		"-o", outPath,
		"--geo-bypass",
		"--no-check-certificate",
		"--quiet",
		"--no-warnings",
		"--http-chunk-size", "10485760", // 10MB Buffer
		"--socket-timeout", "10",
		"--retries", "10",
		"--no-progress",
		url,
	}

	cmd := exec.Command("yt-dlp", args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("download failed: %v", err)
	}

	return outPath, nil
}
