package reports

import "database/sql"

// import (
// 	"database/sql"
// 	"fmt"
// )

// const sqlCallInsertReport = `call insert_report(?,?,?,?,?,?,?,?,?,?,?,?)`

// // const sqlInsertReport = `insert into reports(ru_id, meta_id, report_type, description, created_at) values ((select id from report_urls where path = ?), (select id from url_metadata where query = ? and fragment = ?), ?, ?, ?)`

// // const sqlInsertReportURL = `insert into report_urls(host_id, path, archive_url, title, author_string, content_summary, content_category) values((select id from hostnames where host = ?), ?, ?, ?, ?, ?, ?)`

// // const sqlInsertURLMeta = `insert into url_metadata(query, fragment) values(?, ?)`

// // const sqlInsertHostnames = `insert into hostnames(host) values(?)`

// // const sqlSelectID = `select reports.id, description, createdAt, websiteID, userID, url, host
// // from reports inner join websites on (websites.id=reports.websiteID) where reports.id=?`

// // const sqlSelectAll = `select reports.id, description, createdAt, websiteID, userID, url, host
// // from reports inner join websites on (websites.id=reports.websiteID)
// // order by id, createdAt`

// // const sqlSelectURL = `select reports.id, description, createdAt, websiteID, userID, url, host
// // from reports inner join websites on (websites.id=reports.websiteID) where url=?
// // order by id, createdAt`

// // const sqlSelectHost = `select reports.id, description, createdAt, websiteID, url, host
// // from reports inner join websites on (websites.id=reports.websiteID) where host=?
// // order by id, createdAt`

// // type reportRow struct {
// // 	id          int64
// // 	description string
// // 	createdAt   time.Time
// // 	websiteID   int64
// // 	userID      int64
// // 	url         string
// // 	host        string
// // }

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
func (s *MySQLStore) Insert(report *Report) (*Report, error) {
	return nil, nil
}

// 	//since we need to insert into both the `reports` and `websites`
// 	//tables, and since we want those inserts to be atomic (all or nothing)
// 	//we need to start a database transaction
// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		return nil, fmt.Errorf("error begining transaction: %v", err)
// 	}
// 	// Parameters in order
// 	// host, path, archive_url, query, fragment, report_type, description, created_at, title, author_string, content_summary, content_category
// 	res, err := tx.Exec(sqlCallInsertReport, report.ReportURL.URL.Host, report.ReportURL.URL.EscapedPath(),
// 		report.ReportURL.ArchiveURL, report.ReportURL.URL.RawQuery, report.ReportURL.URL.Fragment, report.ReportType, report.Description,
// 		report.CreatedAt, report.ReportURL.Title, report.ReportURL.AuthorString, report.ReportURL.ContentSummary, report.ReportURL.ContentCategory)
// 	if err != nil {
// 		tx.Rollback()
// 		return nil, fmt.Errorf("error inserting report: %v", err)
// 	}
// 	reportID, err := res.LastInsertId()
// 	if err != nil {
// 		tx.Rollback()
// 		return nil, fmt.Errorf("error getting report id: %v", err)
// 	}

// 	// // try to insert the report
// 	// res, err := tx.Exec(sqlInsertReport, report.ReportURL.URL.EscapedPath(), report.ReportURL.URL.RawQuery, report.ReportURL.URL.Fragment, report.ReportType, report.Description, report.CreatedAt)
// 	// if err != nil {
// 	// 	//rollback the transaction if there's an error
// 	// 	tx.Rollback()
// 	// 	return nil, fmt.Errorf("error inserting report: %v", err)
// 	// }
// 	// reportID, err := res.LastInsertId()
// 	// if err != nil {
// 	// 	tx.Rollback()
// 	// 	return nil, fmt.Errorf("error getting report id: %v", err)
// 	// }

// 	// // insert the report_url
// 	// // TODO: test the fuck out of this... Not sure the difference between raw query and escaped path
// 	// res, err = tx.Exec(sqlInsertReportURL, report.ReportURL.URL.Host, report.ReportURL.URL.EscapedPath(),
// 	// 	report.ReportURL.ArchiveURL, report.ReportURL.Title, report.ReportURL.AuthorString, report.ReportURL.ContentSummary, report.ReportURL.ContentCategory)
// 	// if err != nil {
// 	// 	tx.Rollback()
// 	// 	return nil, fmt.Errorf("error inserting report url: %v", err)
// 	// }
// 	// // reportURLID, err := res.LastInsertId()
// 	// // if err != nil {
// 	// // 	tx.Rollback()
// 	// // 	return nil, fmt.Errorf("error getting report url id: %v", err)
// 	// // }

// 	// // next insert into url_metadata
// 	// res, err = tx.Exec(sqlInsertURLMeta, report.ReportURL.URL.RawQuery, report.ReportURL.URL.Fragment)
// 	// if err != nil {
// 	// 	tx.Rollback()
// 	// 	return nil, fmt.Errorf("error inserting url metadata: %v", err)
// 	// }
// 	// // metaID, err := res.LastInsertId()
// 	// // if err != nil {
// 	// // 	tx.Rollback()
// 	// // 	return nil, fmt.Errorf("error getting url metadata id: %v", err)
// 	// // }

// 	// // last insert into hostnames
// 	// res, err = tx.Exec(sqlInsertHostnames, report.ReportURL.URL.Host)
// 	// if err != nil {
// 	// 	tx.Rollback()
// 	// 	return nil, fmt.Errorf("error inserting hostname: %v", err)
// 	// }
// 	// // hostID, err := res.LastInsertId()
// 	// // if err != nil {
// 	// // 	tx.Rollback()
// 	// // 	return nil, fmt.Errorf("error getting host id: %v", err)
// 	// // }

// 	// set the id of the report
// 	report.ID = reportID

// 	//now commit the transaction so that all those inserts are atomic
// 	if err := tx.Commit(); err != nil {
// 		//try to rollback if we can't commit
// 		tx.Rollback()
// 		return nil, fmt.Errorf("error committing insert transaction: %v", err)
// 	}

// 	return report, nil

// }

// // // GetByID returns the Report with the given id
// // func (s *MySQLStore) GetByID(id int64) (*Report, error) {
// // 	rows, err := s.db.Query(sqlSelectID, id)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("error selecting request by id: %v", err)
// // 	}

// // 	reports, err := scanReports(rows)
// // 	if err != nil || len(reports) == 0 {
// // 		return nil, err
// // 	}

// // 	return reports[0], nil
// // }

// // // GetAll returns all of the reports in the database
// // func (s *MySQLStore) GetAll() ([]*Report, error) {
// // 	rows, err := s.db.Query(sqlSelectAll)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("error selecting reports: %v", err)
// // 	}

// // 	return scanReports(rows)
// // }

// // // GetByURL returns the list of Reports with a given url
// // func (s *MySQLStore) GetByURL(url string) ([]*Report, error) {
// // 	rows, err := s.db.Query(sqlSelectURL, url)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("error selecting request by url: %v", err)
// // 	}

// // 	return scanReports(rows)
// // }

// // // GetByHost returns the list of Reports with a given host
// // func (s *MySQLStore) GetByHost(host string) ([]*Report, error) {
// // 	rows, err := s.db.Query(sqlSelectHost, host)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("error selecting request by host: %v", err)
// // 	}

// // 	return scanReports(rows)
// // }

// // func scanReports(rows *sql.Rows) ([]*Report, error) {
// // 	defer rows.Close()
// // 	reports := []*Report{}
// // 	row := reportRow{}

// // 	for rows.Next() {
// // 		err := rows.Scan(&row.id, &row.description, &row.createdAt, &row.websiteID, &row.userID, &row.url, &row.host)
// // 		if err != nil {
// // 			return nil, fmt.Errorf("error scanning row: %v", err)
// // 		}

// // 		u, err := url.ParseRequestURI(row.url)
// // 		if err != nil {
// // 			return nil, fmt.Errorf("invalid parsing url: %v", err)
// // 		}
// // 		// create a new report and append it to the reports slice
// // 		report := &Report{
// // 			ID:          row.id,
// // 			Description: row.description,
// // 			CreatedAt:   row.createdAt,
// // 			UserID:      row.userID,
// // 			Website:     &Website{ID: row.websiteID, URL: u},
// // 		}
// // 		reports = append(reports, report)
// // 	}

// // 	if err := rows.Err(); err != nil {
// // 		return nil, fmt.Errorf("error reading rows: %v", err)
// // 	}

// // 	return reports, nil
// // }
