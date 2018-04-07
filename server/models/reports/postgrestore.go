package reports

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const insertIntoReports = `insert into reports(user_id, user_description, created_at) values($1,$2,$3) returning id`

const insertIntoTypes = `with s as (
		select id
		from report_types where type = $1
	)
	insert into report_types_reports(report_id, report_type_id)
	select $2, id from s`

//PostgreStore implements Store for a Postgres database
type PostgreStore struct {
	db *sqlx.DB
}

//NewPostgreStore constructs a PostgreStore
func NewPostgreStore(db *sqlx.DB) *PostgreStore {
	if db == nil {
		panic("nil pointer passed to reports.NewPostgreStore")
	}

	return &PostgreStore{
		db: db,
	}
}

// Insert inserts a new report into the database
func (s *PostgreStore) Insert(report *Report) (*Report, error) {
	log.Println("begining report insert transaction")
	// begin a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error begining transaction: %v", err)
	}

	// insert into reports TODO: write user logic
	var rID int64
	if err := tx.QueryRow(insertIntoReports,
		nil, report.UserDescription, report.CreatedAt).Scan(&rID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error inserting report: %v", err)
	}
	// set the report id
	report.ID = rID

	// for each reportType, insert into intermediate table
	for _, rtype := range report.ReportTypes {
		if _, err := tx.Exec(insertIntoTypes, rtype, rID); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error inserting type for report: %v", err)
		}
	}

	//now commit the transaction so that all those inserts are atomic
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error committing insert transaction: %v", err)
	}

	return report, nil
}
