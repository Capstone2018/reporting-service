package reports

import (
	"github.com/jmoiron/sqlx"
)

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
	return nil, nil
}
