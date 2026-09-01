package utils

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

func DownloadImage(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// DrawGlowingCircle creates a circular crop and simulates a glow border
func DrawGlowingCircle(dc *gg.Context, img image.Image, x, y, size float64) {
	// Draw Outer Glow / Border
	dc.DrawCircle(x+size/2, y+size/2, size/2+8)
	dc.SetColor(color.RGBA{255, 255, 255, 200})
	dc.Fill()

	// Draw Image Circularly
	dc.DrawCircle(x+size/2, y+size/2, size/2)
	dc.Clip()
	dc.DrawImage(img, int(x), int(y))
	dc.ResetClip()
}

// GetThumb generates the aesthetic thumbnail using github.com/fogleman/gg
func GetThumb(videoID, userIDStr, userName, title, duration, views, channel, thumbURL string) (string, error) {
	os.MkdirAll("cache", 0755)
	finalPath := fmt.Sprintf("cache/%s_%s.png", videoID, userIDStr)
	if _, err := os.Stat(finalPath); err == nil {
		return finalPath, nil
	}

	tempPath := fmt.Sprintf("cache/temp_%s.jpg", videoID)
	if err := DownloadImage(thumbURL, tempPath); err != nil {
		// Fallback handling
	}

	// 1. Load Background & Blur
	bgImg, err := gg.LoadImage(tempPath)
	if err != nil {
		// Create a fallback black image
		dcFallback := gg.NewContext(1920, 1080)
		dcFallback.SetColor(color.RGBA{30, 30, 30, 255})
		dcFallback.Clear()
		bgImg = dcFallback.Image()
	}

	bgImg = imaging.Resize(bgImg, 1920, 1080, imaging.Lanczos)
	bgImg = imaging.Blur(bgImg, 25.0)            // Gaussian blur
	bgImg = imaging.AdjustBrightness(bgImg, -40) // Darken

	dc := gg.NewContextForImage(bgImg)

	// 2. Draw Rounded Black Rectangle
	dc.DrawRoundedRectangle(40, 40, 1840, 1000, 60)
	dc.SetColor(color.RGBA{0, 0, 0, 180})
	dc.FillPreserve()
	dc.SetColor(color.RGBA{132, 224, 240, 200})
	dc.SetLineWidth(6)
	dc.Stroke()

	// 3. Load Fonts (Assuming assets/font.ttf exists, fallback to basic otherwise)
	err = dc.LoadFontFace("SHIVMUSIC/assets/font.ttf", 65)
	if err != nil {
		// Just relies on internal fallback if font not found
	}

	// 4. Texts
	if len(title) > 22 {
		title = title[:22] + "..."
	}
	dc.SetColor(color.White)
	dc.DrawString(title, 650, 300)

	dc.LoadFontFace("SHIVMUSIC/assets/font2.ttf", 45)
	dc.SetColor(color.RGBA{200, 200, 200, 255})
	dc.DrawString("Artist: "+channel, 650, 400)
	dc.SetColor(color.RGBA{150, 150, 150, 255})
	dc.DrawString("Views: "+views, 650, 470)
	dc.DrawString("Duration: "+duration, 650, 530)

	// 5. Draw Uniform Waveform
	barCount := 64
	barGap := 12
	totalWidth := float64(barCount * barGap)
	startX := (1920.0 - totalWidth) / 2
	baseY := 760.0

	for i := 0; i < barCount; i++ {
		h := float64(rand.Intn(30) + 15)
		x0 := startX + float64(i*barGap)
		if i < barCount/2 {
			dc.SetColor(color.White)
		} else {
			dc.SetColor(color.RGBA{150, 150, 150, 200})
		}
		dc.DrawRoundedRectangle(x0, baseY-h, 5, h*2, 2.5)
		dc.Fill()
	}

	// 6. Draw Line & Icons
	lineY := baseY + 55
	dc.SetColor(color.RGBA{80, 80, 80, 255})
	dc.DrawLine(startX, lineY, startX+totalWidth, lineY)
	dc.Stroke()

	dc.SetColor(color.White)
	dc.DrawLine(startX, lineY, startX+(totalWidth/2), lineY)
	dc.SetLineWidth(2)
	dc.Stroke()

	dc.DrawCircle(startX+(totalWidth/2), lineY, 8)
	dc.Fill()

	// 7. Branding - "THE SHIV"
	dc.LoadFontFace("SHIVMUSIC/assets/font.ttf", 55)
	dc.SetColor(color.RGBA{255, 60, 160, 255})
	dc.DrawString("BETA BOT HUB", 80, 975)
	dc.DrawString("THE SHIV", 1480, 975) // Specific brand requirement mapped

	// Save
	dc.SavePNG(finalPath)
	return finalPath, nil
}
