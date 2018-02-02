package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

//idLength is the length of the ID portion
const idLength = 32

//signedLength is the full length of the signed session ID
//(ID portion plus signature)
const signedLength = idLength + sha256.Size

//SessionID represents a valid, digitally-signed session ID.
//This is a base64 URL encoded string created from a byte slice
//where the first `idLength` bytes are crytographically random
//bytes representing the unique session ID, and the remaining bytes
//are an HMAC hash of those ID bytes (i.e., a digital signature).
//The byte slice layout is like so:
//+-----------------------------------------------------+
//|...32 crypto random bytes...|HMAC hash of those bytes|
//+-----------------------------------------------------+
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("Invalid Session ID")

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
func NewSessionID(signingKey string) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("singning key must not be zero length")
	}

	//- create a byte slice where the first `idLength` of bytes
	//  are cryptographically random bytes for the new session ID,
	//  and the remaining bytes are an HMAC hash of those ID bytes,
	//  using the provided `signingKey` as the HMAC key.
	buf := make([]byte, signedLength)
	if _, err := rand.Read(buf[0:idLength]); err != nil {
		return InvalidSessionID, err
	}

	// generate hmac from id bytes in the remaining buffer
	mac := genMac(buf[0:idLength], signingKey)
	copy(buf[idLength:], mac)

	// encode that byte slice using base64 URL Encoding and return
	//  the result as a SessionID type
	return SessionID(base64.URLEncoding.EncodeToString(buf)), nil
}

//ValidateID validates the string in the `id` parameter
//using the `signingKey` as the HMAC signing key
//and returns an error if invalid, or a SessionID if valid
func ValidateID(id string, signingKey string) (SessionID, error) {
	//validate the `id` parameter using the provided `signingKey`.
	// base64 decode the id to a byte slice
	buf, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return InvalidSessionID, err
	}
	//if the byte slice length is < signedLength
	//it must be invalid, so return InvalidSessionID
	//and ErrInvalidID
	if len(buf) < signedLength {
		return InvalidSessionID, ErrInvalidID
	}
	// HMAC hash the id from the byte slice
	mac := genMac(buf[0:idLength], signingKey)
	// compare the hmac to the one stored in the remaining bytes
	if !hmac.Equal(mac, buf[idLength:]) {
		return InvalidSessionID, ErrInvalidID
	}
	// session is valid so return it
	return SessionID(id), nil

}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	return string(sid)
}

//genMac generates a MAC for a given id and signing key
func genMac(id []byte, signingKey string) []byte {
	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write(id)
	return h.Sum(nil)
}
