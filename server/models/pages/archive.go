package pages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const waybackURL = "https://pragma.archivelab.org"

var client = &http.Client{
	Timeout: time.Second * 30,
}

// Archive represents a wayback archive
type Archive struct {
	ID           int    `json:"id"`
	AnnotationID int    `json:"annotation_id"`
	Protocol     string `json:"protocol"`
	Domain       string `json:"domain"`
	Path         string `json:"path"`
	WaybackID    string `json:"wayback_id"`
}

// NewArchive returns a new instance of an Archive struct
func NewArchive() *Archive {
	return &Archive{}
}

// Archive triggers the wayback machine to archive a url
func (a *Archive) Archive(pageURL string) error {
	jsonReq := fmt.Sprintf(`{"url":"%s"}`, pageURL)
	resp, err := client.Post(waybackURL, "application/json", strings.NewReader(jsonReq))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = receive(resp, a)
	return err
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
