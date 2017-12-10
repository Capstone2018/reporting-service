package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func receive(r *http.Request, target interface{}) error {
	//ensure Content-Type is JSON
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, contentTypeJSON) {
		return fmt.Errorf("`%s` is not a supported Content-Type; must be `%s`", contentType, contentTypeJSONUTF8)
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("error decoding JSON in request body: %v", err)
	}

	return nil
}
