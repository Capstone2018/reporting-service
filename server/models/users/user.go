package users

import (
	"fmt"
	"net/mail"
)

var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	PassHash []byte `json:"-"` //stored, but not encoded to clients
	UserName string `json:"userName"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"` // TODO: decide if we should use usernames
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	//validate the new user
	//- Email field must be a valid email address
	// TODO: decide if we should do email confirmations, would need to get Snopes email creds
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return fmt.Errorf("Email provided is not valid: %s", err.Error())
	}
	//- Password must be at least 6 characters
	if len(nu.Password) < 6 {
		return fmt.Errorf("Password must be at least 6 characters long")
	}
	//- Password and PasswordConf must match
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("Password and password confirmation must match")
	}

	// TODO: decide if usernames should be implemented

	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	//TODO: set the PhotoURL field of the new User to
	//the Gravatar PhotoURL for the user's email address.
	//see https://en.gravatar.com/site/implement/hash/
	//and https://en.gravatar.com/site/implement/images/

	//TODO: also set the ID field of the new User
	//to a new bson ObjectId
	//http://godoc.org/labix.org/v2/mgo/bson

	//TODO: also call .SetPassword() to set the PassHash
	//field of the User to a hash of the NewUser.Password
	return nil, nil
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	//TODO: use the bcrypt package to generate a new hash of the password
	//https://godoc.org/golang.org/x/crypto/bcrypt
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//TODO: use the bcrypt package to compare the supplied
	//password with the stored PassHash
	//https://godoc.org/golang.org/x/crypto/bcrypt
	return nil
}

//Validate validates the user credentials
func (c *Credentials) Validate() error {
	if len(c.Email) == 0 {
		return fmt.Errorf("an email must be provided")
	}
	if len(c.Password) == 0 {
		return fmt.Errorf("a password must be provided")
	}

	return nil
}
