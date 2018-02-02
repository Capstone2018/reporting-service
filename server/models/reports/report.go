package reports

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Capstone2018/reporting-service/server/models/opengraph"
)

// NewReport represents a new report being posted to the service
type NewReport struct {
	UserDescription string `json:"userDescription"`
	ReportType      string `json:"reportType"`
	//OpenGraph       *opengraph.OpenGraph `json:"og"`
	UserID int64 `json:"userID"`
}

// Report represents a fully validated report
type Report struct {
	ID              int64                `json:"id"`
	UserID          int64                `json:"userID"`
	UserDescription string               `json:"userDescription"`
	ReportType      string               `json:"reportType"`
	OpenGraph       *opengraph.OpenGraph `json:"og"`
	CreatedAt       time.Time            `json:"createdAt"`
}

// // ReportURL represents the data assocated with the URL/site of a report
// type ReportURL struct {
// 	ID           int64    `json:"id"`
// 	URL          *url.URL `json:"url"`
// 	ArchiveURL   string   `json:"archiveURL"`
// 	Title        string   `json:"title"`
// 	AuthorString string   `json:"authorString"`
// }

// Validate checks that a new report is valid
func (nr *NewReport) Validate() error {
	// make sure that a userID was passed when creating a NewReport
	if nr.UserID == 0 {
		return fmt.Errorf("no creator ID provided")
	}

	// validate user input
	if len(nr.UserDescription) == 0 {
		return fmt.Errorf("no description provided")
	}
	if len(nr.UserDescription) > 300 {
		return fmt.Errorf("user description len > 300")
	}
	if !(nr.ReportType == "Misleading Title" || nr.ReportType == "False Information" || nr.ReportType == "Other") {
		return fmt.Errorf("report type not allowed value")
	}

	return nil
}

// ToReport turns a new report into a full report
func (nr *NewReport) ToReport() (*Report, error) {
	// validate the new report
	err := nr.Validate()
	if err != nil {
		return nil, err
	}

	report := &Report{
		UserDescription: nr.UserDescription,
		UserID:          nr.UserID,
		ReportType:      nr.ReportType,
		CreatedAt:       time.Now(),
	}

	return report, nil
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
