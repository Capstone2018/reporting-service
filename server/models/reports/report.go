package reports

import (
	"fmt"
	"log"
	"net/url"
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
	URLs            []string `json:"urls"`
	ReportTypes     []string `json:"report_types"`
}

// Report represents a fully validated report
type Report struct {
	ID              int64      `json:"id"`
	Email           string     `json:"email"`
	UserDescription string     `json:"user_description"`
	URLs            []*url.URL `json:"urls"`
	ReportTypes     []string   `json:"report_types"`
	CreatedAt       time.Time  `json:"created_at"`
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

	// create a report with the current time, etc
	report := &Report{
		UserDescription: nr.UserDescription,
		Email:           nr.Email,
		ReportTypes:     nr.ReportTypes,
		CreatedAt:       time.Now(),
	}
	urls := make(map[string]bool)
	// validate the URLs
	for _, pageURL := range nr.URLs {
		// check if url is dupe
		if urls[pageURL] {
			log.Printf("user submitted duplicate url: %v", pageURL)
			continue
		}
		urls[pageURL] = true
		// validate the url (use ParseRequestURI first so we can make sure that it's not relative path)
		_, err := url.ParseRequestURI(pageURL)
		if err != nil {
			log.Printf("invalid url: %v sent from user", pageURL)
			continue
		}
		u, err := url.Parse(pageURL)
		if err != nil {
			log.Printf("error parsing URL: %v", pageURL)
			continue
		}
		report.URLs = append(report.URLs, u)

	}
	return report, nil
}
