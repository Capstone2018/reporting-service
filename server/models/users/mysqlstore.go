package users

import "database/sql"

const sqlInsertUser = `insert into users(username, email) values (?, ?)`
const sqlSelectID = `select id, username, email
from users where users.id=?`

const sqlSelectEmail = `select id, username, email
from users where users.email=?`

const sqlSelectUserName = `select id, username, email
from users where users.username=?`

type userRow struct {
	id       int64
	username string
	email    string
}

//MySQLStore implements Store for a MySQL database
type MySQLStore struct {
	db *sql.DB
}

// NewMySQLStore constructs a MySQLStore
func NewMySQLStore(db *sql.DB) *MySQLStore {
	if db == nil {
		panic("nil pointer passed to NewMySQLStore")
	}

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

// GetByUserName returns a user with a given username
func (s *MySQLStore) GetByUserName(username string) (*User, error) {
	return nil, nil
}

func scanUsers(rows *sql.Rows) ([]*User, error) {
	return nil, nil
}
