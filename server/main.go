package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Capstone2018/reporting-service/server/handlers"
	"github.com/Capstone2018/reporting-service/server/models/reports"
	"github.com/Capstone2018/reporting-service/server/sessions"
	"github.com/go-redis/redis"
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
	// connect to the mysql database
	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		db.Close()
		log.Fatalf("error opening mysql database: %v", err)
	}
	// connection retry logic
	for i := 1; i < maxConnRetries; i++ {
		err = db.Ping()
		// return if we don't find an error
		if err == nil {
			return db, nil
		}
		log.Printf("error connecting to DB server at %s: %s", mysqlConfig.FormatDSN(), err)
		log.Printf("will attempt another connection in %d seconds", i*2)
		time.Sleep(time.Duration(i*2) * time.Second)
	}
	// only close the connection if we hit an error and reached maxConnRetries
	db.Close()
	return nil, err
}

func main() {
	addr := getenv("ADDR", ":443")
	tlsKey := getenv("TLSKEY", "")
	tlsCert := getenv("TLSCERT", "")
	redisAddr := getenv("REDISADDR", "localhost:6379")
	sessionsSigKey := getenv("SESSIONKEY", "")
	// connect to the sql DB and create a new store
	db, _ := connectToMySQL()
	// remember to close the database
	defer db.Close()
	mysqlStore := reports.NewMySQLStore(db)

	// connect to redis, default timeout for sessions to be 1 hour
	// TODO: determine how long a user should be signed in for
	client := redis.NewClient(&redis.Options{Addr: redisAddr})
	redisStore := sessions.NewRedisStore(client, time.Hour)

	// create a handler context
	hctx := handlers.NewHandlerContext(mysqlStore, redisStore, sessionsSigKey)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc(apiRoot+"reports", hctx.Authenticated(hctx.ReportsHandler))
	//apiMux.HandleFunc(apiRoot+"reports/", hctx.Authenticated(hctx.ReportIDHandler))
	serverMux := http.NewServeMux()
	serverMux.Handle(apiRoot, handlers.Adapt(apiMux,
		handlers.CORS(),
	))

	log.Printf("loading TLS certificate from %s\n", tlsCert)
	log.Printf("loading TLS private key from %s\n", tlsKey)
	log.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCert, tlsKey, serverMux))
}
