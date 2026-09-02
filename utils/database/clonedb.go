package database

import (
	"context"
	"log"
	"strconv" // Added for string to int64 conversions

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	CloneOwnerDB   *mongo.Collection
	CloneBotNameDB *mongo.Collection
	CloneBotDB     *mongo.Collection
	CloneCustomDB  *mongo.Collection
	ChatsDBC       *mongo.Collection
	UsersDBC       *mongo.Collection
)

// SaveCloneBotOwner saves the owner ID of a clone bot
func SaveCloneBotOwner(botID int64, userID int64) {
	opts := options.Update().SetUpsert(true)
	update := bson.M{"$set": bson.M{"user_id": userID}}
	_, err := CloneOwnerDB.UpdateOne(context.TODO(), bson.M{"bot_id": botID}, update, opts)
	if err != nil {
		log.Printf("Error in SaveCloneBotOwner: %v", err)
	}
}

// GetOwnerIDFromDB retrieves owner ID directly from the main clone DB
func GetOwnerIDFromDB(botID int64) int64 {
	var result struct {
		UserID int64 `bson:"user_id"`
	}
	err := CloneBotDB.FindOne(context.TODO(), bson.M{"bot_id": botID}).Decode(&result)
	if err != nil {
		return 0
	}
	return result.UserID
}

// GetClonedSupportChat retrieves the support chat link for a clone bot
func GetClonedSupportChat(botID int64) string {
	var result struct {
		Support string `bson:"support"`
	}
	err := CloneBotDB.FindOne(context.TODO(), bson.M{"bot_id": botID}).Decode(&result)
	if err != nil || result.Support == "" {
		return "No support chat set."
	}
	return result.Support
}

// GetCloneSearchSettings retrieves the HIGHEST PRIORITY search preference
func GetCloneSearchSettings(botID int64) (string, string) {
	var data map[string]interface{}
	err := CloneCustomDB.FindOne(context.TODO(), bson.M{"bot_id": botID}).Decode(&data)
	if err != nil {
		return "", ""
	}

	if val, ok := data["video"].(string); ok && val != "" {
		return "video", val
	}
	if val, ok := data["photo"].(string); ok && val != "" {
		return "photo", val
	}
	if val, ok := data["animation"].(string); ok && val != "" {
		return "animation", val
	}
	if val, ok := data["sticker"].(string); ok && val != "" {
		return "sticker", val
	}
	if val, ok := data["text"].(string); ok && val != "" {
		return "text", val
	}
	return "", ""
}

// DeleteCloneSearchType deletes ALL search mode settings
func DeleteCloneSearchType(botID int64) {
	update := bson.M{"$unset": bson.M{
		"video": "", "photo": "", "animation": "", "sticker": "", "text": "",
	}}
	CloneCustomDB.UpdateOne(context.TODO(), bson.M{"bot_id": botID}, update)
}

// GetServedChatsClone fetches all chats served by a specific clone bot
func GetServedChatsClone(botID int64) []int64 {
	var chats []int64
	if ChatsDBC == nil {
		return chats
	}
	cursor, err := ChatsDBC.Find(context.TODO(), bson.M{"bot_id": botID})
	if err != nil {
		return chats
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var chat struct {
			ChatID int64 `bson:"chat_id"`
		}
		if err := cursor.Decode(&chat); err == nil {
			chats = append(chats, chat.ChatID)
		}
	}
	return chats
}

// ==========================================
// NEW MISSING FUNCTIONS ADDED BELOW
// ==========================================

// GetAllClones fetches all clone bots from the database
func GetAllClones() []map[string]interface{} {
	var clones []map[string]interface{}
	if CloneBotDB == nil {
		return clones
	}
	cursor, err := CloneBotDB.Find(context.TODO(), bson.M{})
	if err != nil {
		return clones
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var clone map[string]interface{}
		if err := cursor.Decode(&clone); err == nil {
			clones = append(clones, clone)
		}
	}
	return clones
}

// GetCloneBotOwner wraps GetOwnerIDFromDB to accept a string argument
func GetCloneBotOwner(botID string) int64 {
	id, _ := strconv.ParseInt(botID, 10, 64)
	return GetOwnerIDFromDB(id)
}

// GetServedUsersClone fetches all users served by a specific clone bot
func GetServedUsersClone(botID string) []int64 {
	var users []int64
	id, _ := strconv.ParseInt(botID, 10, 64)
	if UsersDBC == nil {
		return users
	}
	cursor, err := UsersDBC.Find(context.TODO(), bson.M{"bot_id": id})
	if err != nil {
		return users
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user struct {
			UserID int64 `bson:"user_id"`
		}
		if err := cursor.Decode(&user); err == nil {
			users = append(users, user.UserID)
		}
	}
	return users
}
