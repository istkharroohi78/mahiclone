package utils

import (
	"image"
	"image/png"
	"math"
	"os"

	"github.com/nfnt/resize"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	StickerDimX = 512
	StickerDimY = 512
)

// ResizeFileToStickerSize resizes an image to fit 512x512 for Telegram stickers
func ResizeFileToStickerSize(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())

	var newWidth, newHeight uint
	if width < StickerDimX && height < StickerDimY {
		// If already smaller, Telegram will accept it, but we can scale up or leave as is.
		// Replicating Python logic:
		if width > height {
			scale := float64(StickerDimX) / float64(width)
			newWidth = StickerDimX
			newHeight = uint(math.Floor(float64(height) * scale))
		} else {
			scale := float64(StickerDimY) / float64(height)
			newWidth = uint(math.Floor(float64(width) * scale))
			newHeight = StickerDimY
		}
		img = resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	} else {
		// Thumbnail equivalent
		img = resize.Thumbnail(StickerDimX, StickerDimY, img, resize.Lanczos3)
	}

	outPath := filePath + ".png"
	out, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		return "", err
	}

	os.Remove(filePath)
	return outPath, nil
}

// GetDocument returns a telebot Document structure
func GetDocument(filePath string) *tb.Document {
	return &tb.Document{File: tb.FromDisk(filePath)}
}

// GetDocumentFromFileID returns a telebot Document by file ID
func GetDocumentFromFileID(fileID string) *tb.Document {
	return &tb.Document{File: tb.File{FileID: fileID}}
}
