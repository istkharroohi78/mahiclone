package stream

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"ANJALI/config"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	WeekInSeconds     = 7 * 24 * 60 * 60
	OneMinInSeconds   = 60
	MaxCacheSize      = 5 * 1024 * 1024 * 1024 // 5 GB
	MaxFileSize       = 500 * 1024 * 1024      // 500 MB
	AvgSongSize       = 10 * 1024 * 1024
	BufferSpace       = 15 * AvgSongSize
	AdvancedThreshold = MaxCacheSize - BufferSpace
)

type FileInfo struct {
	Path string
	Name string
	Size int64
	Age  float64
}

// AutoClean handles cache cleanup asynchronously
func AutoClean(b *tb.Bot, filePath string) {
	go func() {
		currentTime := time.Now()

		// 1. Image Cleanup
		imageDir := "./downloads"
		files, _ := ioutil.ReadDir(imageDir)
		for _, file := range files {
			if !file.IsDir() && (filepath.Ext(file.Name()) == ".png" || filepath.Ext(file.Name()) == ".jpg") {
				imgPath := filepath.Join(imageDir, file.Name())
				if stat, err := os.Stat(imgPath); err == nil {
					age := currentTime.Sub(stat.ModTime()).Seconds()
					if age > OneMinInSeconds {
						os.Remove(imgPath)
						log.Printf("🗑️ Auto-deleted old image: %s", file.Name())
					}
				}
			}
		}

		// 2. Advanced Song Cache Cleanup
		if filePath == "" {
			return
		}

		dir := filepath.Dir(filePath)
		var allFiles []FileInfo
		var currentCacheSize int64

		mediaFiles, _ := ioutil.ReadDir(dir)
		for _, file := range mediaFiles {
			if !file.IsDir() {
				fPath := filepath.Join(dir, file.Name())
				stat, err := os.Stat(fPath)
				if err == nil {
					// Delete huge files immediately
					if stat.Size() > MaxFileSize {
						os.Remove(fPath)
						continue
					}

					age := currentTime.Sub(stat.ModTime()).Seconds()
					allFiles = append(allFiles, FileInfo{Path: fPath, Name: file.Name(), Size: stat.Size(), Age: age})
					currentCacheSize += stat.Size()
				}
			}
		}

		// Sort by Age (Oldest First)
		sort.Slice(allFiles, func(i, j int) bool {
			return allFiles[i].Age > allFiles[j].Age
		})

		deletedInAdv := 0
		var deletedFiles []string

		for _, f := range allFiles {
			if currentCacheSize > AdvancedThreshold && deletedInAdv < 10 {
				os.Remove(f.Path)
				currentCacheSize -= f.Size
				deletedInAdv++
				deletedFiles = append(deletedFiles, f.Name)
				continue
			}

			if f.Age > WeekInSeconds || currentCacheSize > MaxCacheSize {
				os.Remove(f.Path)
				currentCacheSize -= f.Size
				deletedFiles = append(deletedFiles, f.Name)
			}
		}

		cfg := config.LoadConfig()
		if cfg.LoggerID != 0 && len(deletedFiles) > 0 {
			logChat, err := b.ChatByID(fmt.Sprintf("%d", cfg.LoggerID))
			if err == nil {
				text := fmt.Sprintf("🗑️ **Storage Cache Cleaned**\nCleared %d old files to free space.", len(deletedFiles))
				b.Send(logChat, text, tb.ModeMarkdown)
			}
		}
	}()
}
