package lib

// IsUsernameValid checks if the username a user has chosen is valid.
func IsUsernameValid(username string) bool {
	return 2 < len(username) && len(username) < 20
}

// IsPasswordValid checks if the password a user has chosen is valid.
func IsPasswordValid(password string) bool {
	return 8 <= len(password) && len(password) <= 32
}
