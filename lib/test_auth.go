package lib

import (
	"net/http"
)

// IsAuth takes in a request pointer and returns a boolean telling if a user has a valid token in their cookies.
func IsAuth(r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil {
		return false
	}

	if _, err := ValidateToken(string(cookie.Value)); err == nil {
		return true
	}

	return false
}
