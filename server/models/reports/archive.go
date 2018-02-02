package reports

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const waybackURL = "https://pragma.archivelab.org"

// Archive triggers the wayback machine to archive a url
func Archive(pageURL string) (string, error) {
	jsonReq := fmt.Sprintf(`{"url":"%s"}`, pageURL)
	resp, err := http.Post(waybackURL, "application/json", strings.NewReader(jsonReq))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body := &archiveResponse{}
	if err = receive(resp, body); err != nil {
		return "", err
	}
	return body.WaybackID, nil
}

type archiveResponse struct {
	ID           int    `json:"id"`
	AnnotationID int    `json:"annotation_id"`
	Protocol     string `json:"protocol"`
	Domain       string `json:"domain"`
	Path         string `json:"path"`
	WaybackID    string `json:"wayback_id"`
}

func receive(r *http.Response, target interface{}) error {
	//ensure Content-Type is JSON
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return fmt.Errorf("archive url returned `%s`, which is not a supported Content-Type; must be `%s`", contentType, "application/json")
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("error decoding JSON in archive url response body: %v", err)
	}

	return nil
}
