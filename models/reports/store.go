package reports

import "errors"

// ErrReportNotFound is returned when a report can't be found
var ErrReportNotFound = errors.New("error not found")

// Store is a store for reports
type Store interface {
	// Insert converts a NewReport to a Report, inserts
	// it into the db and returns it
	Insert(newReport *NewReport) (*Report, error)

	// GetByID returns the Report with the given id
	GetByID(id int64) (*Report, error)

	// GetByUrl returns the list of Reports with a given url
	GetByUrl(url string) ([]*Report, error)

	// GetByHost returns the list of Reports with a given host
	GetByHost(host string) ([]*Report, error)
}
