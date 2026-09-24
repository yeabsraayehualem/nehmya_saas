package utils

import (
	"time"

	"github.com/gorilla/sessions"
)

var CookieStore *sessions.CookieStore
var SESSION_VALID = 20
func InitSessions() {
	CookieStore = sessions.NewCookieStore([]byte("XG1HLd4fjGSzVR4S1Lz5jU5rcUXkwCXAsxH0PDeYkPQy4j1w7qAnTR5nMAUFxu0cqVUDcxWcRjNS1nP1dcxp6HwHy0a6jDKvaJLS"))
	CookieStore.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true in production with HTTPS
		MaxAge:   int(time.Duration(SESSION_VALID) * time.Minute / time.Second),
	}
}
