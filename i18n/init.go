package i18n

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Global maps for languages
var (
	Languages        = make(map[string]map[string]string)
	LanguagesPresent = make(map[string]string)
)

// GetString returns the translation map for a specific language
func GetString(lang string) map[string]string {
	if data, exists := Languages[lang]; exists {
		return data
	}
	// Default to English if the requested language is not found
	return Languages["en"]
}

// LoadLanguages reads and processes all YAML language files
func LoadLanguages() {
	langDir := "./strings/langs/"

	// 1. Load English first as the base/fallback language
	enFilePath := filepath.Join(langDir, "en.yml")
	enData, err := os.ReadFile(enFilePath)
	if err != nil {
		log.Fatalf("Failed to read en.yml: %v", err)
	}

	var enParsed map[string]string
	if err := yaml.Unmarshal(enData, &enParsed); err != nil {
		log.Fatalf("Failed to parse en.yml: %v", err)
	}

	Languages["en"] = enParsed
	if name, ok := enParsed["name"]; ok {
		LanguagesPresent["en"] = name
	} else {
		fmt.Println("There is some issue with the language file inside bot.")
		os.Exit(1)
	}

	// 2. Loop through all files in the directory
	files, err := os.ReadDir(langDir)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v", langDir, err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
			langCode := strings.TrimSuffix(file.Name(), ".yml")

			// Skip "en" as it's already loaded
			if langCode == "en" {
				continue
			}

			filePath := filepath.Join(langDir, file.Name())
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Failed to read %s: %v", file.Name(), err)
				continue
			}

			var parsed map[string]string
			if err := yaml.Unmarshal(fileData, &parsed); err != nil {
				log.Printf("Failed to parse %s: %v", file.Name(), err)
				continue
			}

			// Initialize the map for this language
			Languages[langCode] = make(map[string]string)

			// 3. Fallback logic: Copy everything from English first
			for k, v := range Languages["en"] {
				Languages[langCode][k] = v
			}

			// 4. Overwrite with actual translated strings
			for k, v := range parsed {
				Languages[langCode][k] = v
			}

			// 5. Store the language name
			if name, ok := Languages[langCode]["name"]; ok {
				LanguagesPresent[langCode] = name
			} else {
				fmt.Println("There is some issue with the language file inside bot.")
				os.Exit(1)
			}
		}
	}

	log.Println("✅ Language files successfully loaded.")
}
