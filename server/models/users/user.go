package users

import (
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
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
	_, err := mail.ParseAddress(nu.Email) // This doesn't reallllyyy do real validation..
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

	//ensure UserName has non-zero length
	if len(nu.UserName) == 0 {
		return fmt.Errorf("You must supply a user name")
	}

	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	// validate the new user
	if err := nu.Validate(); err != nil {
		return nil, err
	}

	user := &User{
		Email:    nu.Email,
		UserName: nu.UserName,
	}

	return user, user.SetPassword(nu.Password)
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	//bcrypt generate a new hash of the password
	phash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.PassHash = phash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//bcrypt compare the supplied password with the stored PassHash
	return bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
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
