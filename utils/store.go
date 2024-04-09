package utils

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("wowowow"))

func GetSession(h *http.Request) *sessions.Session {
	session, err := Store.Get(h, "garagesale-sesh")
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
