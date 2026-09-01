package database

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	activeChats      = make(map[int64]bool)
	activeVideoChats = make(map[int64]bool)
	assistantDict    = make(map[int64]int)
	autoEndMode      = make(map[int64]bool)
	countMode        = make(map[int64]int)
	channelConnect   = make(map[int64]int64)
	langMap          = make(map[int64]string)
	loopMode         = make(map[int64]int)
	maintenanceMode  []int
	nonAdminChat     = make(map[int64]bool)
	musicPause       = make(map[int64]bool)
	playMode         = make(map[int64]string)
	playType         = make(map[int64]string)
	skipMode         = make(map[int64]bool)
	muteMode         = make(map[int64]bool)
	autoplayMode     = make(map[int64]bool)

	dbMutex sync.RWMutex
)

// Database Collections
var (
	AuthDB    *mongo.Collection
	SudoersDB *mongo.Collection
	BlockedDB *mongo.Collection
	GBanDB    *mongo.Collection
	// Add other collections here...
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// --- Assistant Logic ---
func GetAssistantNumber(chatID int64) int {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return assistantDict[chatID]
}

func SetAssistant(chatID int64, num int) {
	dbMutex.Lock()
	assistantDict[chatID] = num
	dbMutex.Unlock()
}

// --- Play Mode & Type ---
func GetPlaytype(chatID int64) string {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	if val, ok := playType[chatID]; ok {
		return val
	}
	return "Everyone"
}

func SetPlaytype(chatID int64, mode string) {
	dbMutex.Lock()
	playType[chatID] = mode
	dbMutex.Unlock()
}

// --- Status Modifiers ---
func IsSkipMode(chatID int64) bool {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	if val, ok := skipMode[chatID]; ok {
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

func IsAutoplayOn(chatID int64) bool {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return autoplayMode[chatID]
}

func AutoplayOn(chatID int64) {
	dbMutex.Lock()
	autoplayMode[chatID] = true
	dbMutex.Unlock()
}

func AutoplayOff(chatID int64) {
	dbMutex.Lock()
	autoplayMode[chatID] = false
	dbMutex.Unlock()
}

// --- Sudo & Bans ---
func GetSudoers() []int64 {
	var result struct {
		Sudoers []int64 `bson:"sudoers"`
	}
	err := SudoersDB.FindOne(context.TODO(), bson.M{"sudo": "sudo"}).Decode(&result)
	if err != nil {
		return []int64{}
	}
	return result.Sudoers
}

func AddSudo(userID int64) {
	sudoers := GetSudoers()
	sudoers = append(sudoers, userID)
	SudoersDB.UpdateOne(context.TODO(), bson.M{"sudo": "sudo"}, bson.M{"$set": bson.M{"sudoers": sudoers}}, options.Update().SetUpsert(true))
}

func RemoveSudo(userID int64) {
	sudoers := GetSudoers()
	var newSudoers []int64
	for _, id := range sudoers {
		if id != userID {
			newSudoers = append(newSudoers, id)
		}
	}
	SudoersDB.UpdateOne(context.TODO(), bson.M{"sudo": "sudo"}, bson.M{"$set": bson.M{"sudoers": newSudoers}}, options.Update().SetUpsert(true))
}

// --- Maintenance ---
func IsMaintenance() bool {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	if len(maintenanceMode) == 0 {
		return false
	}
	return maintenanceMode[0] == 1
}

func MaintenanceOn() {
	dbMutex.Lock()
	maintenanceMode = []int{1}
	dbMutex.Unlock()
}

func MaintenanceOff() {
	dbMutex.Lock()
	maintenanceMode = []int{2}
	dbMutex.Unlock()
}
