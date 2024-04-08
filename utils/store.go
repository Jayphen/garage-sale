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
