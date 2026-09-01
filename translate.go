package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var targets = map[string]string{
	"hi": "hi",    // Hindi
	"mm": "my",    // Myanmar (Burmese)
	"ar": "ar",    // Arabic
	"bn": "bn",    // Bengali
	"zh": "zh-CN", // Chinese
	"ta": "ta",    // Tamil
	"te": "te",    // Telugu
	"ru": "ru",    // Russian
	"pa": "pa",    // Punjabi
	"kn": "kn",    // Kannada
	"ja": "ja",    // Japanese
}

func main() {
	enFile := "strings/langs/en.yml"
	data, err := os.ReadFile(enFile)
	if err != nil {
		log.Fatalf("❌ Error reading %s: %v", enFile, err)
	}

	var enData map[string]string
	if err := yaml.Unmarshal(data, &enData); err != nil {
		log.Fatalf("❌ Error parsing YAML: %v", err)
	}

	log.Println("✅ Successfully loaded en.yml")

	for langCode, gCode := range targets {
		log.Printf("⏳ Translating to %s...", langCode)
		translatedData := make(map[string]string)

		for key, val := range enData {
			if key == "name" {
				translatedData[key] = langCode
				continue
			}

			translated, err := translateText(val, gCode)
			if err != nil {
				log.Printf("⚠️ Skipped %s: %v", key, err)
				translatedData[key] = val // Fallback to English if translation fails
			} else {
				translatedData[key] = translated
			}

			// ⏱️ DELAY INCREASED: 2 seconds to avoid Google Rate Limit (429 Too Many Requests)
			time.Sleep(2 * time.Second)
		}

		outData, err := yaml.Marshal(&translatedData)
		if err != nil {
			log.Printf("❌ Error marshalling %s.yml: %v", langCode, err)
			continue
		}

		outFile := filepath.Join("strings/langs", langCode+".yml")
		if err := os.WriteFile(outFile, outData, 0644); err != nil {
			log.Printf("❌ Error writing %s.yml: %v", langCode, err)
		} else {
			log.Printf("✅ Successfully saved %s", outFile)
		}
	}
	log.Println("🎉 All translations completed!")
}

func translateText(text, targetLang string) (string, error) {
	if text == "" {
		return "", nil
	}

	apiURL := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=%s&dt=t&q=%s",
		targetLang, url.QueryEscape(text))

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 🛡️ STATUS CODE CHECK: Return error if Google sends HTML error page instead of JSON
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Google API blocked request (Status: %d). Increase delay", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data []interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("JSON parse error (might be HTML): %v", err)
	}

	if len(data) > 0 && data[0] != nil {
		parsed := data[0].([]interface{})
		var finalStr strings.Builder
		for _, chunk := range parsed {
			chunkSlice := chunk.([]interface{})
			finalStr.WriteString(chunkSlice[0].(string))
		}

		res := finalStr.String()
		res = strings.ReplaceAll(res, "{ ", "{")
		res = strings.ReplaceAll(res, " }", "}")
		res = strings.ReplaceAll(res, "< tg - emoji", "<tg-emoji")
		res = strings.ReplaceAll(res, "</ tg - emoji >", "</tg-emoji>")
		res = strings.ReplaceAll(res, "emoji - id", "emoji-id")
		res = strings.ReplaceAll(res, "=\" ", "=\"")
		return res, nil
	}
	return text, nil
}
