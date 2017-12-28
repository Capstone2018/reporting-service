package handlers

import (
	"net/http"

	"github.com/Capstone2018/reporting-service/server/sessions"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *SessionState)

func (ctx *Context) Authenticated(handlerFunc AuthenticatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessState := &SessionState{}
		_, err := sessions.GetState(r, ctx.SessionSigningKey, ctx.SessionStore, sessState)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		handlerFunc(w, r, sessState)
	}
}
