package platforms

import (
	"fmt"
	"path/filepath"

	"ANJALI/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

type TeleAPI struct{}

func NewTeleAPI() *TeleAPI {
	return &TeleAPI{}
}

func (t *TeleAPI) GetFilename(doc *tb.Document, audio bool) string {
	if doc != nil && doc.FileName != "" {
		return doc.FileName
	}
	if audio {
		return "ᴛᴇʟᴇɢʀᴀᴍ ᴀᴜᴅɪᴏ"
	}
	return "ᴛᴇʟᴇɢʀᴀᴍ ᴠɪᴅᴇᴏ"
}

func (t *TeleAPI) GetDuration(duration int) string {
	if duration == 0 {
		return "Unknown"
	}
	return utils.SecondsToMin(duration)
}

func (t *TeleAPI) GetFilepath(audio *tb.Audio, video *tb.Video) string {
	fileName := "unknown"
	if audio != nil {
		fileName = audio.FileID + ".ogg"
	} else if video != nil {
		fileName = video.FileID + ".mp4"
	}

	path, _ := filepath.Abs(filepath.Join("downloads", fileName))
	return path
}

func (t *TeleAPI) Download(b *tb.Bot, m *tb.Message, dest string) error {
	// Telebot provides an easy download method for media
	if m.Audio != nil {
		return b.Download(&m.Audio.File, dest)
	} else if m.Video != nil {
		return b.Download(&m.Video.File, dest)
	} else if m.Document != nil {
		return b.Download(&m.Document.File, dest)
	}
	return fmt.Errorf("no media found")
}
