package lib

import "unicode"

// IsUsernameValid checks if the username a user has chosen is valid.
func IsUsernameValid(username string) bool {
	return 2 < len(username) && len(username) < 20
}

// IsPasswordValid checks if the password a user has chosen is valid.
func IsPasswordValid(password string) bool {
	return 8 <= len(password) && len(password) <= 32
}

func ValidatePasswordStrength(password string) bool {
	var (
		upper, lower, numbers bool
		totalLength           uint8
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upper = true
			totalLength++
		case unicode.IsLower(char):
			upper = true
			totalLength++
		case unicode.IsNumber(char):
			numbers = true
			totalLength++
		default:
			return false
		}
	}

	if !upper || !lower || !numbers || totalLength < 8 {
		return false
	}

	return true
}
