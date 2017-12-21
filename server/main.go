package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Capstone2018/reporting-service/server/handlers"
	"github.com/Capstone2018/reporting-service/server/models/reports"
	"github.com/go-sql-driver/mysql"
)

const apiRoot = "/v1/"

func getenv(name string, def string) string {
	val := os.Getenv(name)
	if len(val) == 0 {
		if len(def) > 0 {
			return def
		}
		log.Fatalf("please set the %s environment variable", name)
	}
	return val
}

const maxConnRetries = 10

func connectToMySQL() (*sql.DB, error) {
	mysqlAddr := getenv("MYSQL_ADDR", "localhost")

	//construct the connection string
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.Addr = mysqlAddr
	mysqlConfig.Net = "tcp"

	mysqlConfig.DBName = getenv("MYSQL_DATABASE", "")
	mysqlConfig.User = "root"
	mysqlConfig.Passwd = getenv("MYSQL_ROOT_PASSWORD", "")
	//tell the MySQL driver to parse DATETIME
	//column values into go time.Time values
	mysqlConfig.ParseTime = true

	log.Println(mysqlConfig.FormatDSN())
	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		db.Close()
		log.Fatalf("error opening mysql database: %v", err)
	}
	for i := 1; i < maxConnRetries; i++ {
		err = db.Ping()
		if err == nil {
			return db, nil
		}
		log.Printf("error connecting to DB server at %s: %s", mysqlConfig.FormatDSN(), err)
		log.Printf("will attempt another connection in %d seconds", i*2)
		time.Sleep(time.Duration(i*2) * time.Second)
	}
	db.Close()
	return nil, err
}

func main() {
	addr := getenv("ADDR", ":443")
	tlsKey := getenv("TLSKEY", "")
	tlsCert := getenv("TLSCERT", "")
	db, _ := connectToMySQL()
	// remember to close the database
	defer db.Close()
	mysqlStore := reports.NewMySQLStore(db)
	hctx := handlers.NewHandlerContext(mysqlStore)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc(apiRoot+"reports", hctx.ReportsHandler)

	serverMux := http.NewServeMux()
	serverMux.Handle(apiRoot, handlers.Adapt(apiMux,
		handlers.CORS(),
	))

	log.Printf("loading TLS certificate from %s\n", tlsCert)
	log.Printf("loading TLS private key from %s\n", tlsKey)
	log.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCert, tlsKey, serverMux))
}
