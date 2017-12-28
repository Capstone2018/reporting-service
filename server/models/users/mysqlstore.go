package users

import "database/sql"

type userRow struct {
	id    int64
	email string
}

//MySQLStore implements Store for a MySQL database
type MySQLStore struct {
	db *sql.DB
}

// NewMySQLStore constructs a MySQLStore
func NewMySQLStore(db *sql.DB) *MySQLStore {
	return &MySQLStore{
		db: db,
	}
}

// Insert inserts a new user to the database
func (s *MySQLStore) Insert(newUser *NewUser) (*User, error) {
	return nil, nil
}

// GetByID returns a user with a given id
func (s *MySQLStore) GetByID(id int64) (*User, error) {
	return nil, nil
}

// GetByEmail returns a user with a given email
func (s *MySQLStore) GetByEmail(email string) (*User, error) {
	return nil, nil
}
