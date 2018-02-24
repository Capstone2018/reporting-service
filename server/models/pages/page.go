package pages

import "time"

// Page represents a page on the internet
type Page struct {
	ID              int64     `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	URLString       string    `json:"url_string"`
	WaybackID       string    `json:"wayback_id"`
	URLID           int64     `json:"url_id"`
	OgID            int64     `json:"og_id"`
	ReportID        int64     `json:"report_id"`
	QueryFragmentID int64     `json:"query_fragment_id"`
}
