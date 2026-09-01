package utils

import (
	"sync"
)

var (
	activeChats      = make(map[int64]bool)
	activeVideoChats = make(map[int64]bool)
	skipMode         = make(map[int64]bool)
	autoEndMode      = make(map[int64]bool)
	musicPause       = make(map[int64]bool)
	playMode         = make(map[int64]string)
	playType         = make(map[int64]string)
	channelConnect   = make(map[int64]int64)
	upvoteCount      = make(map[int64]int)
	loopMode         = make(map[int64]int)

	dbMutex sync.RWMutex
)

// Skip Mode
func IsSkipMode(chatID int64) bool {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	val, exists := skipMode[chatID]
	if exists {
		return val
	}
	return true
}

func SkipOn(chatID int64) {
	dbMutex.Lock()
	skipMode[chatID] = true
	dbMutex.Unlock()
}

func SkipOff(chatID int64) {
	dbMutex.Lock()
	skipMode[chatID] = false
	dbMutex.Unlock()
}

// Playmode & Playtype
func GetPlaymode(chatID int64) string {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	mode, ok := playMode[chatID]
	if ok {
		return mode
	}
	return "Direct"
}

func SetPlaymode(chatID int64, mode string) {
	dbMutex.Lock()
	playMode[chatID] = mode
	dbMutex.Unlock()
}

func GetPlaytype(chatID int64) string {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	mode, ok := playType[chatID]
	if ok {
		return mode
	}
	return "Everyone"
}

func SetPlaytype(chatID int64, mode string) {
	dbMutex.Lock()
	playType[chatID] = mode
	dbMutex.Unlock()
}

// Upvotes & Channel Mode
func GetUpvoteCount(chatID int64) int {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	val, ok := upvoteCount[chatID]
	if ok {
		return val
	}
	return 5
}

func SetUpvotes(chatID int64, mode int) {
	dbMutex.Lock()
	upvoteCount[chatID] = mode
	dbMutex.Unlock()
}

func GetCMode(chatID int64) int64 {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return channelConnect[chatID]
}

func SetCMode(chatID, targetChatID int64) {
	dbMutex.Lock()
	channelConnect[chatID] = targetChatID
	dbMutex.Unlock()
}

// Stream Status
func IsMusicPlaying(chatID int64) bool {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return musicPause[chatID]
}

func MusicOn(chatID int64) {
	dbMutex.Lock()
	musicPause[chatID] = true
	dbMutex.Unlock()
}

func MusicOff(chatID int64) {
	dbMutex.Lock()
	musicPause[chatID] = false
	dbMutex.Unlock()
}

// Loop Settings
func GetLoop(chatID int64) int {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return loopMode[chatID]
}

func SetLoop(chatID int64, mode int) {
	dbMutex.Lock()
	loopMode[chatID] = mode
	dbMutex.Unlock()
}
