package reports

import (
	"fmt"
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
