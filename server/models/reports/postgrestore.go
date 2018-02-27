package reports

import (
	"database/sql"
)

//PostgreStore implements Store for a Postgres database
type PostgreStore struct {
	db *sql.DB
}

//NewPostgreStore constructs a PostgreStore
func NewPostgreStore(db *sql.DB) *PostgreStore {
	if db == nil {
		panic("nil pointer passed to NewMySQLStore")
	}

	return &PostgreStore{
		db: db,
	}
}

// Insert inserts a new report into the database
func (s *PostgreStore) Insert(report *Report) (*Report, error) {
	return nil, nil
}
