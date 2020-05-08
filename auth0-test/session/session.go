package session

import (
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

const DEFAULT = "default"

func SetupStore() {
	key := os.Getenv("SESSION_KEY")
	Store = sessions.NewCookieStore([]byte(key))
}
