package users

import "errors"

//ErrUserNotFound is returned when the user can't be found
var ErrUserNotFound = errors.New("user not found")

//Store is a store for users
type Store interface {
	// Insert inserts a new user to the database
	Insert(newUser *NewUser) (*User, error)

	// GetByID returns a user with a given id
	GetByID(id int64) (*User, error)

	// GetByEmail returns a user with a given email
	GetByEmail(email string) (*User, error)
}
