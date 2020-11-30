package lib

import "github.com/gorilla/sessions"

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// GetStore returns a cookie storage, in which we can check and add cookies.
func GetStore() *sessions.CookieStore {
	return store
}
