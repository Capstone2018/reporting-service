package reports

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"
)

const sqlInsertReport = `insert into reports(description, createdAt, websiteID) values (?, ?, ?)`

const sqlInsertWebsite = `insert into websites(url, host) values(?, ?)`

const sqlSelectID = `select id, description, createdAt, websiteID, url, host 
from reports inner join websites on (websites.id=reports.websiteID) where id=?`

const sqlSelectURL = `select id, description, createdAt, websiteID, url, host 
from reports inner join websites on (websites.id=reports.websiteID) where url=?
order by id, createdAt`

const sqlSelectHost = `select id, description, createdAt, websiteID, url, host 
from reports inner join websites on (websites.id=reports.websiteID) where host=?
order by id, createdAt`

type reportRow struct {
	id          int64
	description string
	createdAt   time.Time
	websiteID   int64
	url         string
	host        string
}

//MySQLStore implements Store for a MySQL database
type MySQLStore struct {
	db *sql.DB
}

//NewMySQLStore constructs a MySQLStore
func NewMySQLStore(db *sql.DB) *MySQLStore {
	if db == nil {
		panic("nil pointer passed to NewMySQLStore")
	}

	return &MySQLStore{
		db: db,
	}
}

// Insert inserts a new report into the database
func (s *MySQLStore) Insert(nr *NewReport) (*Report, error) {
	report, err := nr.ToReport()
	if err != nil {
		return nil, err
	}

	//since we need to insert into both the `reports` and `websites`
	//tables, and since we want those inserts to be atomic (all or nothing)
	//we need to start a database transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error begining transaction: %v", err)
	}

	// first insert into the website so we can get the id
	res, err := tx.Exec(sqlInsertWebsite, report.Website.URL.EscapedPath(), report.Website.URL.Host)
	if err != nil {
		//rollback the transaction if there's an error
		tx.Rollback()
		return nil, fmt.Errorf("error inserting website: %v", err)
	}
	websiteID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error getting website id: %v", err)
	}

	// then insert the report
	res, err = tx.Exec(sqlInsertReport, report.Description, report.CreatedAt, websiteID)
	if err != nil {
		//rollback the transaction if there's an error
		tx.Rollback()
		return nil, fmt.Errorf("error inserting report: %v", err)
	}
	reportID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error getting report id: %v", err)
	}

	report.ID = reportID

	//now commit the transaction so that all those inserts are atomic
	if err := tx.Commit(); err != nil {
		//try to rollback if we can't commit
		tx.Rollback()
		return nil, fmt.Errorf("error committing insert transaction: %v", err)
	}

	return report, nil

}

// GetByID returns the Report with the given id
func (s *MySQLStore) GetByID(id int64) (*Report, error) {
	rows, err := s.db.Query(sqlSelectID, id)
	if err != nil {
		return nil, fmt.Errorf("error selecting request by id: %v", err)
	}

	reports, err := scanReports(rows)
	if err != nil || len(reports) == 0 {
		return nil, err
	}

	return reports[0], nil
}

// GetByURL returns the list of Reports with a given url
func (s *MySQLStore) GetByURL(url string) ([]*Report, error) {
	rows, err := s.db.Query(sqlSelectURL, url)
	if err != nil {
		return nil, fmt.Errorf("error selecting request by url: %v", err)
	}

	return scanReports(rows)
}

// GetByHost returns the list of Reports with a given host
func (s *MySQLStore) GetByHost(host string) ([]*Report, error) {
	rows, err := s.db.Query(sqlSelectHost, host)
	if err != nil {
		return nil, fmt.Errorf("error selecting request by host: %v", err)
	}

	return scanReports(rows)
}

func scanReports(rows *sql.Rows) ([]*Report, error) {
	defer rows.Close()
	reports := []*Report{}
	row := reportRow{}

	for rows.Next() {
		err := rows.Scan(&row.id, &row.description, &row.createdAt, &row.websiteID, &row.url, &row.host)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		u, err := url.ParseRequestURI(row.url)
		if err != nil {
			return nil, fmt.Errorf("invalid parsing url: %v", err)
		}
		// create a new report and append it to the reports slice
		report := &Report{
			ID:          row.id,
			Description: row.description,
			CreatedAt:   row.createdAt,
			Website:     &Website{ID: row.websiteID, URL: u},
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	return reports, nil
}
