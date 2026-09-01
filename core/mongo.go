package core

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	tb "gopkg.in/tucnak/telebot.v2"

	// Apne project module ka naam 'ANJALI' lagayein
	"ANJALI/config"
)

// MongoDB exposes the database instance globally from the core package
var MongoDB *mongo.Database

const tempMongoDBURI = "mongodb://localhost:27017"

// InitMongoDB initializes the database connection
func InitMongoDB() {
	uri := config.MongoDBURI
	dbName := "Anon"

	// Fallback logic if no Mongo URI is provided
	if uri == "" || uri == "None" {
		log.Println("[WARNING] No MONGO DB URL found. LOL — the shiv")
		uri = tempMongoDBURI

		// Initialize a temporary bot to fetch its username
		tempBot, err := tb.NewBot(tb.Settings{
			Token: config.BotToken,
		})
		
		if err == nil && tempBot.Me != nil && tempBot.Me.Username != "" {
			dbName = tempBot.Me.Username
		} else {
			log.Printf("[WARNING] Could not fetch bot username, defaulting DB name. Error: %v\n", err)
			dbName = "FallbackDB"
		}
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect to MongoDB: %v\n", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("[ERROR] MongoDB Ping failed: %v\n", err)
	}

	// Initialize the global database variable
	MongoDB = client.Database(dbName)
	log.Println("[INFO] MongoDB connected successfully. — the shiv")
}
