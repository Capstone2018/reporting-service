package reports

import (
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
