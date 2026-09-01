package platforms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var Themes = []string{
	"3024-night", "a11y-dark", "blackboard", "base16-dark", "cobalt",
	"dracula-pro", "hopscotch", "material", "monokai", "nord",
	"one-dark", "panda-syntax", "seti", "synthwave-84", "twilight",
}

var Colours = []string{
	"#FF0000", "#FF5733", "#FFFF00", "#008000", "#0000FF", "#800080",
	"#00FFFF", "#30D5C8", "#00FF00", "#4B0082", "#FFC0CB", "#000000",
}

type CarbonAPI struct {
	Language         string
	DropShadow       bool
	DropShadowBlur   string
	DropShadowOffset string
	FontFamily       string
	WidthAdjustment  bool
	Watermark        bool
}

func NewCarbonAPI() *CarbonAPI {
	rand.Seed(time.Now().UnixNano())
	return &CarbonAPI{
		Language:         "auto",
		DropShadow:       true,
		DropShadowBlur:   "68px",
		DropShadowOffset: "20px",
		FontFamily:       "JetBrains Mono",
		WidthAdjustment:  true,
		Watermark:        false,
	}
}

func (c *CarbonAPI) Generate(text string, userID int64) (string, error) {
	if err := os.MkdirAll("cache", 0755); err != nil {
		return "", err
	}

	payload := map[string]interface{}{
		"code":                 text,
		"backgroundColor":      Colours[rand.Intn(len(Colours))],
		"theme":                Themes[rand.Intn(len(Themes))],
		"dropShadow":           c.DropShadow,
		"dropShadowOffsetY":    c.DropShadowOffset,
		"dropShadowBlurRadius": c.DropShadowBlur,
		"fontFamily":           c.FontFamily,
		"language":             c.Language,
		"watermark":            c.Watermark,
		"widthAdjustment":      c.WidthAdjustment,
	}

	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://carbonara.solopov.dev/api/cook", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot reach the host")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("cache/carbon%d.png", userID)
	if err := os.WriteFile(filePath, body, 0644); err != nil {
		return "", err
	}

	return filepath.Abs(filePath)
}
