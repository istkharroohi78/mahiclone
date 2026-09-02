package utils

import "ANJALI/config"

// IsSudoer checks if the user ID exists in the Sudoers list
func IsSudoer(userID int64) bool {
	for _, id := range config.Sudoers {
		if id == userID {
			return true
		}
	}
	return false
}
