package reports

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// NewReport represents a new report being posted to the service
type NewReport struct {
	Description     string `json:"description"`
	ReportType      string `json:"reportType"`
	ContentCategory string `json:"contentCategory"`
	//ContentSummary string `json:"contentSummary"` //TODO: decide if this is parsed in front or back end
	Title        string `json:"title"`
	AuthorString string `json:"authorString"`
	URL          string `json:"url"`
	UserID       int64  `json:"userID"`
}

// Report represents a fully validated report
type Report struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"userID"`
	Description string     `json:"description"`
	ReportType  string     `json:"reportType"`
	ReportURL   *ReportURL `json:"reportURL"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// ReportURL represents the data assocated with the URL/site of a report
type ReportURL struct {
	ID              int64    `json:"id"`
	URL             *url.URL `json:"url"`
	ArchiveURL      string   `json:"archiveURL"`
	Title           string   `json:"title"`
	AuthorString    string   `json:"authorString"`
	ContentSummary  string   `json:"contentSummary"`
	ContentCategory string   `json:"contentCategory"`
}

// Validate checks that a new report is valid
func (nr *NewReport) Validate() (*url.URL, error) {
	// if nr.UserID == 0 {
	// 	return nil, fmt.Errorf("no creator ID provided")
	// } // TODO: authentication decisions
	if len(nr.Description) == 0 {
		return nil, fmt.Errorf("no description provided")
	}
	if len(nr.ReportType) == 0 {
		return nil, fmt.Errorf("no report type provided")
	}
	if len(nr.ContentCategory) == 0 {
		return nil, fmt.Errorf("no content category provided")
	}
	u, err := url.ParseRequestURI(nr.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid report url: %v", err)
	}
	return u, nil
}

// ToReport turns a new report into a full report
func (nr *NewReport) ToReport() (*Report, error) {
	u, err := nr.Validate()
	if err != nil {
		return nil, err
	}
	// archive.org the submitted URL
	waybackID, err := archive(u.RawPath)
	if err != nil {
		// TODO: decide to not actually fail if you can't archive, you probably still want to save to db...?
		return nil, err
	}
	// TODO: figure out where content summary is generated from

	// create the report struct so we can store it as a report
	reportURL := &ReportURL{
		URL:             u,
		ArchiveURL:      waybackID,
		Title:           nr.Title,
		AuthorString:    nr.AuthorString,
		ContentCategory: nr.ContentCategory,
	}
	report := &Report{
		Description: nr.Description,
		ReportType:  nr.ReportType,
		CreatedAt:   time.Now(),
		UserID:      nr.UserID,
		ReportURL:   reportURL,
	}
	return report, nil
}

type archiveResponse struct {
	Domain    string `json:"domain"`
	ID        string `json:"id"`
	WaybackID string `json:"waybackID"`
}

// triggers the wayback machine to archive a url
func archive(pageURL string) (string, error) {
	jsonReq := fmt.Sprintf(`{"url":"%s"}`, pageURL)
	resp, err := http.Post("https://pragma.archivelab.org", "application/json", strings.NewReader(jsonReq))
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

func receive(r *http.Response, target interface{}) error {
	//ensure Content-Type is JSON
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return fmt.Errorf("`%s` is not a supported Content-Type; must be `%s`", contentType, "application/json")
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("error decoding JSON in request body: %v", err)
	}

	return nil
}
