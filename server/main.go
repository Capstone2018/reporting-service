package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Capstone2018/reporting-service/server/databases/postgres"
	"github.com/Capstone2018/reporting-service/server/handlers"
	"github.com/Capstone2018/reporting-service/server/models/pages"
	"github.com/Capstone2018/reporting-service/server/models/reports"
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
	// redisAddr := getenv("REDISADDR", "localhost:6379")
	// sessionsSigKey := getenv("SESSIONKEY", "")

	psqlHost := getenv("POSTGRES_HOST", "localhost")
	psqlPort := getenv("POSTGRES_PORT", "5432")
	psqlUser := getenv("POSTGRES_USER", "admin")
	psqlPassword := getenv("POSTGRES_PASSWORD", "")
	psqlDatabase := getenv("POSTGRES_DB", "reports")
	// connect to the sql DB and create a new store
	cfg := postgres.Config{
		Host:     psqlHost,
		Port:     psqlPort,
		User:     psqlUser,
		Password: psqlPassword,
		Database: psqlDatabase,
	}
	db, err := postgres.New(cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// remember to close the database TODO: abstract this method to check error
	defer db.Close()
	reportStore := reports.NewPostgreStore(db)
	pageStore := pages.NewPostgreStore(db)

	// connect to redis, default timeout for sessions to be 1 hour
	// TODO: (maybe use it later for other caching?)
	// client := redis.NewClient(&redis.Options{Addr: redisAddr})
	// redisStore := sessions.NewRedisStore(client, time.Hour)

	// create a handler context
	hctx := handlers.NewHandlerContext(reportStore, pageStore)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc(apiRoot+"reports", hctx.ReportsHandler)
	apiMux.HandleFunc(apiRoot+"pages", hctx.PagesHandler)
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
