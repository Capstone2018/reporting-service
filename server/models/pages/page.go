package pages

import (
	"net/url"
	"time"

	"github.com/Capstone2018/reporting-service/server/models/reports"
)

// Page represents a page on the internet
type Page struct {
	ID        int64           `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	URLString string          `json:"url_string"`
	WaybackID string          `json:"wayback_id,omitempty"`
	URL       *url.URL        `json:"url,omitempty"`
	OpenGraph *OpenGraph      `json:"og,omitempty"`
	Report    *reports.Report `json:"report,omitempty"`
	Query     string          `json:"query,omitempty"`
	Fragment  string          `json:"fragment,omitempty"`
}
