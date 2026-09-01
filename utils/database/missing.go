package database

// GetActiveChats returns a list of active streaming chats
func GetActiveChats() []string {
	var chats []string
	for k := range ActiveChats {
		chats = append(chats, k)
	}
	return chats
}

func GetAllClones() []map[string]interface{} { return nil }
func GetCloneBotOwner(botID string) int64 { return 0 }
func GetServedUsersClone(botID string) []int64 { return nil }
