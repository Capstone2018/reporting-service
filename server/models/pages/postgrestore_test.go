package pages

import (
	"testing"
)

func TestPageInsert(t *testing.T) {
	// // initalize a database connection
	// cfg := postgres.Config{
	// 	Host:     "devpsql",
	// 	Port:     "5432",
	// 	User:     "admin",
	// 	Password: "supersecret",
	// 	Database: "reporting",
	// }
	// db, err := postgres.Open(cfg)
	// if err != nil {
	// 	t.Errorf("error establishing postgres connection")
	// }
	// defer db.Close()
	// // create a pages store
	// s := pages.NewPostgreStore(db)

}

func TestMain(m *testing.M) {
	//
	// models.TestDBManager.Enter()
	// // os.Exit() does not respect defer statements
	// ret := m.Run()
	// models.TestDBManager.Exit()
	// os.Exit(ret)
}
