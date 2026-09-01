package config

import (
	"os"
	"strconv"
	"strings"
)

// Struct for files that use: cfg := config.LoadConfig()
type Config struct {
	APIID                int
	APIHash              string
	BotToken             string
	BotID                string
	OwnerUsername        string
	BotUsername          string
	BotName              string
	AssUsername          string
	BotLink              string
	MongoDBURI           string
	JioSaavnAPI          string
	APIURL               string
	APIKey               string
	YTProxyURL           string
	YTAPIKey             string
	WorkerFallbackAPIURL string
	WorkerFallbackAPIKey string
	InflexAPIURL         string
	InflexAPIKey         string
	DurationLimit        int
	LoggerID             int64
	Logger2ID            int64
	CloneLogger          int64
	CloneLogger2         int64
	ErrorLoggerID        int64
	OwnerID              int64
	Sudoers              []int64
	SupportChannel       string
	SupportChat          string
	Github               string
	SpotifyClientID      string
	SpotifyClientSecret  string
	StringSessions       []string
	StartImgURL          []string
	CMBot                []string
	EffectID             []int64
}

// Global variables for files that use: config.SupportChat
var (
	APIID                int
	APIHash              string
	BotToken             string
	BotID                string
	OwnerUsername        string
	BotUsername          string
	BotName              string
	AssUsername          string
	BotLink              string
	MongoDBURI           string
	JioSaavnAPI          string
	APIURL               string
	APIKey               string
	YTProxyURL           string
	YTAPIKey             string
	WorkerFallbackAPIURL string
	WorkerFallbackAPIKey string
	InflexAPIURL         string
	InflexAPIKey         string
	DurationLimit        int
	LoggerID             int64
	Logger2ID            int64
	CloneLogger          int64
	CloneLogger2         int64
	ErrorLoggerID        int64
	OwnerID              int64
	Sudoers              []int64
	SupportChannel       string
	SupportChat          string
	Github               string
	SpotifyClientID      string
	SpotifyClientSecret  string
	StringSessions       []string
	StartImgURL          []string
	CMBot                []string
	EffectID             []int64
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}

func getEnvAsInt64(key string, fallback int64) int64 {
	strValue := getEnv(key, "")
	if value, err := strconv.ParseInt(strValue, 10, 64); err == nil {
		return value
	}
	return fallback
}

func getEnvAsSlice(key, fallback string) []string {
	strValue := getEnv(key, fallback)
	return strings.Fields(strValue)
}

func timeToSeconds(timeStr string) int {
	parts := strings.Split(timeStr, ":")
	seconds := 0
	multiplier := 1
	for i := len(parts) - 1; i >= 0; i-- {
		val, _ := strconv.Atoi(parts[i])
		seconds += val * multiplier
		multiplier *= 60
	}
	return seconds
}

func init() {
	LoadConfig()
}

// LoadConfig globals update karega aur end mein struct return karega
func LoadConfig() *Config {
	durationLimitMin := getEnv("DURATION_LIMIT", "17000")

	APIID = getEnvAsInt("API_ID", 0)
	APIHash = getEnv("API_HASH", "")
	BotToken = getEnv("BOT_TOKEN", "")
	BotID = getEnv("BOT_ID", "")
	OwnerUsername = getEnv("OWNER_USERNAME", "")
	BotUsername = getEnv("BOT_USERNAME", "")
	BotName = getEnv("BOT_NAME", "")
	AssUsername = getEnv("ASSUSERNAME", "")
	BotLink = getEnv("BOT_LINK", "https://t.me/royal_musics_bot")
	MongoDBURI = getEnv("MONGO_DB_URI", "")
	
	JioSaavnAPI = getEnv("JIOSAAVN_API", "https://saavn.me/search/songs?query=")

	APIURL = getEnv("API_URL", "https://api.shrutibots.site")
	APIKey = getEnv("API_KEY", "ShrutiBotsC0WH1GowF2HkGoKv4F3y")
	YTProxyURL = getEnv("YTPROXY_URL", "https://tgapi.xbitcode.com")
	YTAPIKey = getEnv("YT_API_KEY", "xbit_B4TNnBAoe6uoSM7NLFz-dk6X7GibJ6Bh")
	WorkerFallbackAPIURL = getEnv("WORKER_FALLBACK_API_URL", "https://youtubenewapi.skybotsdeveloper.workers.dev")
	WorkerFallbackAPIKey = getEnv("WORKER_FALLBACK_API_KEY", "itsmesid")
	InflexAPIURL = getEnv("INFLEX_API_URL", "https://teaminflex.xyz")
	InflexAPIKey = getEnv("INFLEX_API_KEY", "INFLEX40920628D")

	DurationLimit = timeToSeconds(durationLimitMin + ":00")

	LoggerID = getEnvAsInt64("LOGGER_ID", 0)
	Logger2ID = getEnvAsInt64("LOGGER_2_ID", -1003255930328)
	CloneLogger = getEnvAsInt64("LOGGER_ID", 0)
	CloneLogger2 = getEnvAsInt64("CLONE_LOGGER_2", -1003255930328)
	ErrorLoggerID = -1004392214389

	OwnerID = getEnvAsInt64("OWNER_ID", 8418584090)
	Sudoers = []int64{8418584090, 8723235165}

	SupportChannel = getEnv("SUPPORT_CHANNEL", "https://t.me/betabot_hub")
	SupportChat = getEnv("SUPPORT_CHAT", "https://t.me/betabot_support")
	Github = getEnv("GITHUB", "https://t.me/sukoon_s")

	SpotifyClientID = getEnv("SPOTIFY_CLIENT_ID", "63f2d3fb20c84cfaa472e5c3b805cd6b")
	SpotifyClientSecret = getEnv("SPOTIFY_CLIENT_SECRET", "c0b5b18383c2447fb9bd13f7eae42a57")

	StringSessions = []string{
		getEnv("STRING_SESSION", ""),
		getEnv("STRING_SESSION2", ""),
		getEnv("STRING_SESSION3", ""),
		getEnv("STRING_SESSION4", ""),
	}

	StartImgURL = getEnvAsSlice("START_IMG_URL", "https://files.catbox.moe/n22tbs.jpg")

	CMBot = []string{
		"💞", "🥂", "🔍", "🧪", "⚡️", "🔥", "🦋", "🎩", "🌈", "🍷",
		"🥃", "🥤", "🕊️", "💌", "🧨", "✨", "💥", "💯", "🌟", "⚡️",
	}
	EffectID = []int64{
		5046509860389126442,
		5107584321108051014,
		5104841245755180586,
		5159385139981059251,
	}

	// Struct return taaki 'logger.go' error na de
	return &Config{
		APIID:                APIID,
		APIHash:              APIHash,
		BotToken:             BotToken,
		BotID:                BotID,
		OwnerUsername:        OwnerUsername,
		BotUsername:          BotUsername,
		BotName:              BotName,
		AssUsername:          AssUsername,
		BotLink:              BotLink,
		MongoDBURI:           MongoDBURI,
		JioSaavnAPI:          JioSaavnAPI,
		APIURL:               APIURL,
		APIKey:               APIKey,
		YTProxyURL:           YTProxyURL,
		YTAPIKey:             YTAPIKey,
		WorkerFallbackAPIURL: WorkerFallbackAPIURL,
		WorkerFallbackAPIKey: WorkerFallbackAPIKey,
		InflexAPIURL:         InflexAPIURL,
		InflexAPIKey:         InflexAPIKey,
		DurationLimit:        DurationLimit,
		LoggerID:             LoggerID,
		Logger2ID:            Logger2ID,
		CloneLogger:          CloneLogger,
		CloneLogger2:         CloneLogger2,
		ErrorLoggerID:        ErrorLoggerID,
		OwnerID:              OwnerID,
		Sudoers:              Sudoers,
		SupportChannel:       SupportChannel,
		SupportChat:          SupportChat,
		Github:               Github,
		SpotifyClientID:      SpotifyClientID,
		SpotifyClientSecret:  SpotifyClientSecret,
		StringSessions:       StringSessions,
		StartImgURL:          StartImgURL,
		CMBot:                CMBot,
		EffectID:             EffectID,
	}
}
