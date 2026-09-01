package utils

import (
	"context"
	"os"
	"strings"

	"ANJALI/config"
	"ANJALI/utils/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Check if running on Heroku
func IsHeroku() bool {
	hostname, err := os.Hostname()
	if err != nil {
		return false
	}
	return strings.Contains(hostname, "heroku")
}

// Loads sudoers from DB and ensures Owner is inside
func InitSudoers() {
	cfg := config.LoadConfig()

	var res struct {
		Sudoers []int64 `bson:"sudoers"`
	}

	// Fetch sudoers from MongoDB
	database.SudoersDB.FindOne(context.TODO(), bson.M{"sudo": "sudo"}).Decode(&res)

	// Ensure OwnerID is inside Sudoers list
	ownerExists := false
	for _, id := range res.Sudoers {
		if id == cfg.OwnerID {
			ownerExists = true
			break
		}
	}

	if !ownerExists {
		res.Sudoers = append(res.Sudoers, cfg.OwnerID)
		opts := options.Update().SetUpsert(true)
		database.SudoersDB.UpdateOne(context.TODO(), bson.M{"sudo": "sudo"}, bson.M{"$set": bson.M{"sudoers": res.Sudoers}}, opts)
	}

	// Inject back to Config
	cfg.Sudoers = res.Sudoers
	Logger.Println("𝗦𝗨𝗗𝗢 𝗨𝗦𝗘𝗥 𝗗𝗢𝗡𝗘✨🎋.")
}
