package utils

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func CreateStore() {
	store = sessions.NewCookieStore([]byte("PotatOBananaSausages"))

	store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
}

func GetSession(h *http.Request) *sessions.Session {
	session, err := store.Get(h, "garagesale")
	if err != nil {
		return nil
	}
	return session
}

func GetCartSize(h *http.Request) int {
	session := GetSession(h)

	// Retrieve cart ID from the session if it exists
	cartSize, ok := session.Values["cartSize"].(int)
	if !ok {
		cartSize = 0
	}

	return cartSize
}
