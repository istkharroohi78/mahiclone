package utils

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ClearText truncates text with an ellipsis if it exceeds max length
func ClearText(text string, maxLen int) string {
	text = strings.TrimSpace(text)
	runes := []rune(text)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}
	return text
}

// DownloadUserPhoto fetches the user's profile picture from Telegram
func DownloadUserPhoto(bot *tgbotapi.BotAPI, userID int64) string {
	os.MkdirAll("cache", 0755)
	filePath := fmt.Sprintf("cache/%d.jpg", userID)

	photos, err := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
		UserID: userID,
		Limit:  1,
	})

	if err != nil || photos.TotalCount == 0 {
		return ""
	}

	fileID := photos.Photos[0][0].FileID
	fileURL, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		return ""
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return ""
	}
	defer out.Close()

	io.Copy(out, resp.Body)
	return filePath
}

// GetGlowingCircle creates a circular cropped image with a fake "glow" and border using 'gg'
func GetGlowingCircle(img image.Image, dc *gg.Context, x, y, size float64) {
	// Circular Crop and Border
	radius := size / 2

	// Glow effect (simulated with fading translucent circles)
	for i := 1.0; i <= 5.0; i++ {
		alpha := int(50 / i)
		dc.DrawCircle(x, y, radius+(i*5))
		dc.SetRGBA255(255, 105, 180, alpha)
		dc.Fill()
	}

	// White Border
	dc.DrawCircle(x, y, radius+4)
	dc.SetColor(color.White)
	dc.SetLineWidth(8)
	dc.Stroke()

	// Circular Clip for Image
	dc.DrawCircle(x, y, radius)
	dc.Clip()

	// Draw Image centered
	bounds := img.Bounds()
	imgW, imgH := float64(bounds.Dx()), float64(bounds.Dy())
	dc.DrawImageAnchored(img, int(x), int(y), 0.5, 0.5)

	dc.ResetClip()
}

// GetThumb generates the glowing thumbnail (Bot and Owner text removed)
func GetThumb(bot *tgbotapi.BotAPI, videoID string, userID int64, title, channel, views, duration, ytThumbURL string) string {
	os.MkdirAll("cache", 0755)
	botID := bot.Self.ID
	filename := fmt.Sprintf("cache/%s_%d.png", videoID, botID)

	if _, err := os.Stat(filename); err == nil {
		return filename
	}

	// Download YT Thumbnail
	tempBg := fmt.Sprintf("cache/temp_%s.jpg", videoID)
	resp, err := http.Get(ytThumbURL)
	if err == nil {
		out, _ := os.Create(tempBg)
		io.Copy(out, resp.Body)
		out.Close()
		resp.Body.Close()
	}
	defer os.Remove(tempBg)

	// Load and process background
	bgImg, err := imaging.Open(tempBg)
	if err != nil {
		bgImg = image.NewRGBA(image.Rect(0, 0, 1920, 1080)) // Fallback solid color
	}
	bgImg = imaging.Fill(bgImg, 1920, 1080, imaging.Center, imaging.Lanczos)
	bgImg = imaging.Blur(bgImg, 25.0)
	bgImg = imaging.AdjustBrightness(bgImg, -40) // Darken the background

	// Setup Drawing Context
	dc := gg.NewContext(1920, 1080)
	dc.DrawImage(bgImg, 0, 0)

	// Draw Black Card (Rounded Rectangle)
	dc.DrawRoundedRectangle(40, 40, 1840, 1000, 60) // Adjusted height since top text is removed
	dc.SetRGBA255(0, 0, 0, 200)                     // Semi-transparent black
	dc.FillPreserve()
	dc.SetRGBA255(132, 224, 240, 200) // Cyan outline
	dc.SetLineWidth(6)
	dc.Stroke()

	// Load Fonts
	fontPath1 := "SHIVMUSIC/assets/font.ttf"
	fontPath2 := "SHIVMUSIC/assets/font2.ttf"

	// YT Image Circular Placement
	ytImg, err := imaging.Open(tempBg)
	if err == nil {
		ytImg = imaging.Fill(ytImg, 500, 500, imaging.Center, imaging.Lanczos)
		GetGlowingCircle(ytImg, dc, 350, 350, 450)
	}

	// User Image Circular Placement
	uPhoto := DownloadUserPhoto(bot, userID)
	if uPhoto != "" {
		uImg, err := imaging.Open(uPhoto)
		if err == nil {
			uImg = imaging.Fill(uImg, 450, 450, imaging.Center, imaging.Lanczos)
			uImg = imaging.Blur(uImg, 6.0) // Specific Gaussian Blur 6 for user photo
			GetGlowingCircle(uImg, dc, 1550, 350, 400)
		}
		defer os.Remove(uPhoto)
	}

	// Text Placement
	if err := dc.LoadFontFace(fontPath1, 65); err == nil {
		dc.SetColor(color.White)
		dc.DrawString(ClearText(title, 25), 650, 300)
	}

	if err := dc.LoadFontFace(fontPath2, 45); err == nil {
		dc.SetRGBA255(200, 200, 200, 255)
		dc.DrawString(fmt.Sprintf("Artist: %s", channel), 650, 400)

		dc.SetRGBA255(150, 150, 150, 255)
		dc.DrawString(fmt.Sprintf("Views: %s", views), 650, 470)
		dc.DrawString(fmt.Sprintf("Duration: %s", duration), 650, 530)
	}

	// Uniform Waveform
	barCount := 64
	barWidth := 4.0
	barGap := 10.0
	totalWidth := float64(barCount) * barGap
	startX := (1920.0 - totalWidth) / 2
	baseY := 780.0

	for i := 0; i < barCount; i++ {
		distFromCenter := math.Abs(float64(i - (barCount / 2)))
		h := 20.0
		if distFromCenter < 5 {
			h = 35.0
		}

		x0 := startX + (float64(i) * barGap)
		y0 := baseY - h

		if i < (barCount / 2) {
			dc.SetRGBA255(255, 255, 255, 255)
		} else {
			dc.SetRGBA255(150, 150, 150, 200)
		}
		dc.DrawRoundedRectangle(x0, y0, barWidth, h*2, 2)
		dc.Fill()
	}

	// Processing Line & Icons
	lineY := baseY + 60
	dc.SetRGBA255(100, 100, 100, 255)
	dc.DrawLine(startX, lineY, startX+totalWidth, lineY)
	dc.SetLineWidth(1)
	dc.Stroke()

	dc.SetRGBA255(255, 255, 255, 255)
	dc.DrawLine(startX, lineY, startX+(totalWidth/2), lineY)
	dc.SetLineWidth(2)
	dc.Stroke()

	thumbX := startX + (totalWidth / 2)
	dc.DrawCircle(thumbX, lineY, 8)
	dc.SetColor(color.White)
	dc.Fill()
	dc.DrawCircle(thumbX, lineY, 3)
	dc.SetColor(color.Black)
	dc.Fill()

	if err := dc.LoadFontFace(fontPath2, 30); err == nil {
		dc.SetColor(color.White)
		dc.DrawString("00:00", startX, lineY+35)
		dc.DrawString(duration, startX+totalWidth-80, lineY+35)
	}

	// Control Buttons
	ctrlY := lineY + 50
	midX := 960.0

	// Play/Pause Circle
	dc.DrawCircle(midX, ctrlY, 25)
	dc.SetColor(color.White)
	dc.SetLineWidth(2)
	dc.Stroke()

	// Play Triangle
	dc.DrawRegularPolygon(3, midX, ctrlY, 10, math.Pi)
	dc.SetColor(color.White)
	dc.Fill()

	// Previous & Next Ellipses
	dc.DrawEllipse(midX-47, ctrlY, 12, 15)
	dc.SetLineWidth(1)
	dc.Stroke()

	dc.DrawEllipse(midX+47, ctrlY, 12, 15)
	dc.SetLineWidth(1)
	dc.Stroke()

	// Save Output
	if err := dc.SavePNG(filename); err != nil {
		log.Printf("Failed to save thumbnail: %v", err)
	}

	return filename
}
