package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	//- create a new SessionID
	sid, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	//- save the sessionState to the store
	if err := store.Save(sid, sessionState); err != nil {
		return InvalidSessionID, err
	}
	//- add a header to the ResponseWriter that looks like this:
	//    "Authorization: Bearer <sessionID>"
	//  where "<sessionID>" is replaced with the newly-created SessionID
	w.Header().Add(headerAuthorization, schemeBearer+sid.String())

	//return the new SessionID and nil
	return sid, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//get the value of the Authorization header,
	//or the "auth" query string parameter if no Authorization header is present
	auth := r.Header.Get(headerAuthorization)
	if len(auth) == 0 {
		auth = r.FormValue(paramAuthorization)
	}
	//if it's zero-length, return InvalidSessionID and ErrNoSessionID
	if len(auth) == 0 {
		return InvalidSessionID, ErrNoSessionID
	}
	//if it doesn't start with "Bearer ",
	//return InvalidSessionID and ErrInvalidScheme
	if len(auth) == 0 || !strings.HasPrefix(auth, schemeBearer) {
		return InvalidSessionID, ErrInvalidScheme
	}

	//trim off the "Bearer " prefix and validate the remaining id
	//if you get an error return InvalidSessionID and the error
	sid, err := ValidateID(strings.TrimPrefix(auth, schemeBearer), signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	//return the validated SessionID and nil
	return sid, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	//get the SessionID from the request
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return sid, err
	}
	//get the data associated with that SessionID from the store.
	if err := store.Get(sid, sessionState); err != nil {
		return sid, err
	}
	return sid, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	//get the SessionID from the request
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return sid, err
	}
	//delete the data associated with it in the store.
	if err := store.Delete(sid); err != nil {
		return sid, err
	}
	// return the seesionID to the user
	return sid, nil
}
