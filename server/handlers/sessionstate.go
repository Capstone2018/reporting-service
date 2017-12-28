package handlers

import (
	"net/http"
	"time"

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
)

//SessionState represents the state of an authenticated user session
type SessionState struct {
	User       *users.User
	Began      time.Time
	ClientAddr string
}

//NewSessionState constructs a new SessionState
func NewSessionState(r *http.Request, user *users.User) *SessionState {
	return &SessionState{
		Began:      time.Now(),
		ClientAddr: r.RemoteAddr,
		User:       user,
	}
}
