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
	adminCache       = make(map[int64][]int64) 

	dbMutex sync.RWMutex
)

// Database Collections
var (
	AuthDB             *mongo.Collection
	SudoersDB          *mongo.Collection
	BlockedDB          *mongo.Collection
	GBanDB             *mongo.Collection
	BlacklistedChatsDB *mongo.Collection
	ServedChatsDB      *mongo.Collection 
	DailyStatsDB       *mongo.Collection 
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// --- General Functions ---

func GetLang(chatID int64) string {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	if val, ok := langMap[chatID]; ok {
		return val
	}
	return "en" 
}

func SetLang(chatID int64, lang string) {
	dbMutex.Lock()
	langMap[chatID] = lang
	dbMutex.Unlock()
}

func IsNonAdminChat(chatID int64) bool {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return nonAdminChat[chatID]
}

func GetAdminCache(chatID int64) []int64 {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return adminCache[chatID]
}

func GetActiveChats() []int64 {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	var chats []int64
	for chatID, active := range activeChats {
		if active {
			chats = append(chats, chatID)
		}
	}
	return chats
}

func GetActiveVideoChats() []int64 {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	var chats []int64
	for chatID, active := range activeVideoChats {
		if active {
			chats = append(chats, chatID)
		}
	}
	return chats
}

// --- MongoDB Chat Blacklist Functions (Fixed returns) ---

func GetBlacklistedChats() []int64 {
	var result []int64
	if BlacklistedChatsDB == nil {
		return result
	}
	cursor, err := BlacklistedChatsDB.Find(context.TODO(), bson.M{})
	if err != nil {
		return result
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var doc struct {
			ChatID int64 `bson:"chat_id"`
		}
		if err := cursor.Decode(&doc); err == nil {
			result = append(result, doc.ChatID)
		}
	}
	return result
}

func BlacklistChat(chatID int64) bool {
	if BlacklistedChatsDB == nil {
		return false
	}
	_, err := BlacklistedChatsDB.UpdateOne(context.TODO(), bson.M{"chat_id": chatID}, bson.M{"$set": bson.M{"chat_id": chatID}}, options.Update().SetUpsert(true))
	return err == nil
}

func WhitelistChat(chatID int64) bool {
	if BlacklistedChatsDB == nil {
		return false
	}
	_, err := BlacklistedChatsDB.DeleteOne(context.TODO(), bson.M{"chat_id": chatID})
	return err == nil
}

// --- MongoDB Banned Users Functions (NEW) ---

func AddBannedUser(userID int64) bool {
	if BlockedDB == nil { return false }
	_, err := BlockedDB.UpdateOne(context.TODO(), bson.M{"user_id": userID}, bson.M{"$set": bson.M{"user_id": userID}}, options.Update().SetUpsert(true))
	return err == nil
}

func RemoveBannedUser(userID int64) bool {
	if BlockedDB == nil { return false }
	_, err := BlockedDB.DeleteOne(context.TODO(), bson.M{"user_id": userID})
	return err == nil
}

// --- MongoDB Served Chats & Users Functions ---

func GetServedChats() []int64 {
	var result []int64
	if ServedChatsDB == nil {
		return result
	}
	cursor, err := ServedChatsDB.Find(context.TODO(), bson.M{})
	if err != nil {
		return result
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var doc struct {
			ChatID int64 `bson:"chat_id"`
		}
		if err := cursor.Decode(&doc); err == nil {
			result = append(result, doc.ChatID)
		}
	}
	return result
}

func GetServedUsers() []int64 {
	var result []int64
	if AuthDB == nil { return result }
	cursor, err := AuthDB.Find(context.TODO(), bson.M{})
	if err != nil { return result }
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var doc struct {
			UserID int64 `bson:"user_id"`
		}
		if err := cursor.Decode(&doc); err == nil {
			result = append(result, doc.UserID)
		}
	}
	return result
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

// Wrapper to fix spelling error in decorators/play.go
func GetPlayType(chatID int64) string {
	return GetPlaytype(chatID)
}

func SetPlaytype(chatID int64, mode string) {
	dbMutex.Lock()
	playType[chatID] = mode
	dbMutex.Unlock()
}

func GetPlayMode(chatID int64) string {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	if val, ok := playMode[chatID]; ok {
		return val
	}
	return "Direct"
}

func GetCMode(chatID int64) int64 {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	if val, ok := channelConnect[chatID]; ok {
		return val
	}
	return 0
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
	if SudoersDB == nil {
		return []int64{}
	}
	err := SudoersDB.FindOne(context.TODO(), bson.M{"sudo": "sudo"}).Decode(&result)
	if err != nil {
		return []int64{}
	}
	return result.Sudoers
}

func AddSudo(userID int64) {
	if SudoersDB == nil {
		return
	}
	sudoers := GetSudoers()
	sudoers = append(sudoers, userID)
	SudoersDB.UpdateOne(context.TODO(), bson.M{"sudo": "sudo"}, bson.M{"$set": bson.M{"sudoers": sudoers}}, options.Update().SetUpsert(true))
}

func RemoveSudo(userID int64) {
	if SudoersDB == nil {
		return
	}
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
