package main

import (
	"io"
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Write to both stdout and log.txt
	multiWriter := io.MultiWriter(os.Stdout, file)
	Logger = log.New(multiWriter, "[SHIVMUSIC] ", log.Ldate|log.Ltime|log.Lshortfile)
}
