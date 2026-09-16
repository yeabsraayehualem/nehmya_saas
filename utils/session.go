package utils

import (
	"time"

	"github.com/gorilla/sessions"
)

// NewCookieStore returns a gorilla session cookie store configured with sane
// cookie defaults. The secret should be a long random value; keep it out of
// source control (e.g. load it from an environment variable).
func NewCookieStore(secret []byte, maxAge time.Duration) *sessions.CookieStore {
	store := sessions.NewCookieStore(secret)
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set to true when serving over HTTPS
		MaxAge:   int(maxAge.Seconds()),
	}
	return store
}
