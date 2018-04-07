package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	// import init for postgres
	_ "github.com/lib/pq"
)

const maxConnRetries = 10

// Config holds the configuration used for instantiating a new Roach.
type Config struct {
	// Address that locates our postgres instance
	Host string
	// Port to connect to
	Port string
	// User that has access to the database
	User string
	// Password so that the user can login
	Password string
	// Database to connect to (must have been created priorly)
	Database string
}

// Open opens a new db connection to a postgres instance
// Abstracted without attempting to ping for faster tests without retry
func Open(cfg Config) (db *sqlx.DB, err error) {
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" ||
		cfg.Password == "" || cfg.Database == "" {
		err = errors.Errorf(
			"All fields must be set (%s)",
			spew.Sdump(cfg))
		return
	}
	// The first argument corresponds to the driver name that the driver
	// (in this case, `lib/pq`) used to register itself in `database/sql`.
	// The next argument specifies the parameters to be used in the connection.
	// Details about this string can be seen at https://godoc.org/github.com/lib/pq
	conf := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port)

	db, err = sqlx.Open("postgres", conf)
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't open connection to postgre database (%s)",
			spew.Sdump(cfg))
		return
	}

	return
}

// Connect connects to a Postgres instance
func Connect(cfg Config) (db *sqlx.DB, err error) {
	// Open the connection
	d, err := Open(cfg)
	// if we got an error don't try and ping (nil pointer?)
	if err != nil {
		return
	}
	// retry logic to wait for the database to initialize
	for i := 1; i < maxConnRetries; i++ {
		err = d.Ping()
		// return if we don't find an error
		if err == nil {
			db = d
			return
		}
		log.Printf("error connecting to DB server at %s: %s", spew.Sdump(cfg), err)
		log.Printf("will attempt another connection in %d seconds", i*2)
		time.Sleep(time.Duration(i*2) * time.Second)
	}

	// if we got here we've reached the max retries -- return an error
	err = errors.Wrapf(err,
		"Couldn't ping postgre database (%s)",
		spew.Sdump(cfg))
	return
}

// initialize a new database called test
func initTestDB(db *sqlx.DB) (err error) {
	// crete a test database
	if _, err = db.Exec(`create database test with template reporting`); err != nil {
		err = errors.Wrap(err,
			"Couldn't create test database from template reporting database")
		return
	}
	return
}

// TestConnect returns a postgres connection and creates a test db to test against
func TestConnect(cfg Config) (db *sqlx.DB, err error) {
	// Open the connection
	d, err := Open(cfg)
	// don't attempt to ping
	if err != nil {
		return
	}
	// attempt a ping
	err = d.Ping()
	// initialize a new test database with tables from 'reporting db
	if err == nil {
		initTestDB(d)
		db = d
		return
	}
	// otherwise return an error
	err = errors.Wrapf(err,
		"Couldn't ping postgre database (%s)",
		spew.Sdump(cfg))
	return
}

// TestClose tears down a test database connection
func TestClose() (err error) {
	return
}

// // Close performs the release of any resources that
// // `sql/database` DB pool created. This is usually meant
// // to be used in the exitting of a program or `panic`ing.
// func (p *Postgres) Close() (err error) {
// 	if p.Db == nil {
// 		return
// 	}

// 	if err = p.Db.Close(); err != nil {
// 		err = errors.Wrapf(err,
// 			"Errored closing database connection",
// 			spew.Sdump(p.cfg))
// 	}

// 	return
// }
