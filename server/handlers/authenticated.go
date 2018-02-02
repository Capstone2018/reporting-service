package handlers

import (
	"net/http"

	"github.com/Capstone2018/reporting-service/server/sessions"
	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
)

// AuthenticatedHandler defines an interface of a handler that can take a session state as parameter
type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *SessionState)

// Authenticated is a middleware handler that returns a HanderFunc
func (ctx *Context) Authenticated(handlerFunc AuthenticatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessState := &SessionState{}
		_, err := sessions.GetState(r, ctx.SessionSigningKey, ctx.SessionStore, sessState)
		if err != nil {
			// generate a dummy user for users who aren't authenticated
			// only on the POST /reports
			if r.Method == "POST" && r.URL.Path == "/reports" {
				sessState.User = &users.User{
					ID:    1,
					Email: "test@snopes.com",
				}
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

		}

		handlerFunc(w, r, sessState)
	}
}
