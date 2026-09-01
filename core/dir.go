package core

import (
	"os"
	"strings"

	"ANJALI/logging" // Tumhare root module 'ANJALI' se logging package import ho raha hai
)

// Dirr clears image files and ensures required directories exist
func Dirr() {
	// Current directory ki saari files aur folders ko read karna
	entries, err := os.ReadDir(".")
	if err != nil {
		logging.ErrorLogger.Printf("Error reading directory: %v\n", err)
		return
	}

	// Loop chalakar .jpg, .jpeg, .png files delete karna
	for _, entry := range entries {
		// Agar entry folder hai, toh skip karo
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		// Extension check karne ke liye lowercase me convert kar lena safe rehta hai
		lowerName := strings.ToLower(fileName)

		if strings.HasSuffix(lowerName, ".jpg") ||
			strings.HasSuffix(lowerName, ".jpeg") ||
			strings.HasSuffix(lowerName, ".png") {

			err := os.Remove(fileName)
			if err != nil {
				logging.ErrorLogger.Printf("Failed to remove file %s: %v\n", fileName, err)
			}
		}
	}

	// 'downloads' aur 'cache' directories create karna
	if err := os.MkdirAll("downloads", 0755); err != nil {
		logging.ErrorLogger.Printf("Failed to create 'downloads' directory: %v\n", err)
	}

	if err := os.MkdirAll("cache", 0755); err != nil {
		logging.ErrorLogger.Printf("Failed to create 'cache' directory: %v\n", err)
	}

	logging.InfoLogger.Println("Directories Updated.")
}
