package reports

import (
	"fmt"
	"net/url"
	"time"
)

type NewReport struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

type Report struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Website     *Website  `json:"website"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Website struct {
	ID  int64    `json:"id"`
	URL *url.URL `json:"url"`
}

// Validate checks that a new report is valid
func (nr *NewReport) Validate() (*url.URL, error) {
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
		Website: &Website{
			URL: u,
		},
	}
	return report, nil
}
