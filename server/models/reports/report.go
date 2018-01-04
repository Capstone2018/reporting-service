package reports

import (
	"fmt"
	"net/url"
	"time"
)

// NewReport represents a new report being posted to the service
type NewReport struct {
	Description string `json:"description"`
	URL         string `json:"url"`
	CreatorID   int64  `json:"creatorID"`
}

// Report represents a fully validated report
type Report struct {
	ID          int64     `json:"id"`
	CreatorID   int64     `json:"creatorID"`
	Description string    `json:"description"`
	Website     *Website  `json:"website"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Website represents a website that has been reported on
type Website struct {
	ID  int64    `json:"id"`
	URL *url.URL `json:"url"`
}

// Validate checks that a new report is valid
func (nr *NewReport) Validate() (*url.URL, error) {
	if nr.CreatorID == 0 {
		return nil, fmt.Errorf("no creator ID provided")
	}
	if len(nr.Description) == 0 {
		return nil, fmt.Errorf("no description provided")
	}
	u, err := url.ParseRequestURI(nr.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %v", err)
	}
	return u, nil
}

// ToReport turns a new report into a full report
func (nr *NewReport) ToReport() (*Report, error) {
	u, err := nr.Validate()
	if err != nil {
		return nil, err
	}

	report := &Report{
		Description: nr.Description,
		CreatedAt:   time.Now(),
		CreatorID:   nr.CreatorID,
		Website: &Website{
			URL: u,
		},
	}
	return report, nil
}
