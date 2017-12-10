package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Capstone2018/reporting-service/handlers"
	"github.com/Capstone2018/reporting-service/models/reports"
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

func main() {
	addr := getenv("ADDR", ":443")
	tlsKey := getenv("TLSKEY", "")
	tlsCert := getenv("TLSCERT", "")
	mysqlAddr := getenv("MYSQL_ADDR", "localhost")

	//construct the connection string
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.Addr = mysqlAddr
	mysqlConfig.DBName = getenv("MYSQL_DATABASE", "")
	mysqlConfig.User = "root"
	mysqlConfig.Passwd = getenv("MYSQL_ROOT_PASSWORD", "")
	//tell the MySQL driver to parse DATETIME
	//column values into go time.Time values
	mysqlConfig.ParseTime = true

	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		log.Fatalf("error opening mysql database: %v", err)
	}
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
