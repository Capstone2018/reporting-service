package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

//respond encodes `val` into JSON and writes it to the response,
//setting the response status code to the value of `statusCode`.
//if `statusCode` is zero, http.StatusOK will be used.
func respond(w http.ResponseWriter, val interface{}, statusCode int) {
	w.Header().Add(headerContentType, contentTypeJSONUTF8)
	if statusCode > 0 {
		w.WriteHeader(statusCode)
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(val); err != nil {
		log.Printf("error encoding response value to JSON: %v", err)
	}
}
