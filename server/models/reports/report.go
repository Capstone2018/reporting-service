package reports

import (
	"fmt"
	"time"
)

type page struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

// NewReport represents a new report being posted to the service
type NewReport struct {
	Email           string   `json:"email"`
	UserDescription string   `json:"user_description"`
	Pages           []page   `json:"pages"`
	ReportTypes     []string `json:"report_types"`
}

// Report represents a fully validated report
type Report struct {
	ID              int64     `json:"id"`
	Email           string    `json:"email"`
	UserDescription string    `json:"user_description"`
	Pages           []page    `json:"pages"`
	ReportTypes     []string  `json:"report_types"`
	CreatedAt       time.Time `json:"created_at"`
}

// Validate checks that a new report is valid
func (nr *NewReport) Validate() error {

	// validate user input
	if len(nr.UserDescription) == 0 {
		return fmt.Errorf("no description provided")
	}
	if len(nr.UserDescription) > 300 {
		return fmt.Errorf("user description len > 300")
	}
	// check that the user didn't provide too many report types
	if len(nr.ReportTypes) > 5 {
		return fmt.Errorf("user provided more than 5 report types")
	}

	// TODO: validate all the report types from a list of valid reports

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
		Pages:           nr.Pages,
		Email:           nr.Email,
		ReportTypes:     nr.ReportTypes,
		CreatedAt:       time.Now(),
	}

	return report, nil
}
