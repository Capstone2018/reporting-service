package pages

import (
	"fmt"
	"log"
	"net/url"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//const insertIntoHostnames = `insert into hostnames(host) values($1)`
const insertIntoURLs = `with s as (
		select id, "host", "path"
		from urls
		where host = $1 and path = $2
	), i as (
		insert into urls ("host", "path")
		select $1, $2
		where not exists (select 1 from s)
		returning id, "host", "path"
	)
	select id
	from i
	union all
	select id
	from s`

//const insertIntoQueryFragment = `insert into query_fragment(query, fragment) values($1,$2) returning id`
const insertIntoQueryFragment = `with s as (
	select id, "query", "fragment"
	from urls
	where query = $1 and fragment = $2
	), i as (
		insert into urls ("query", "fragment")
		select $1, $2
		where not exists (select 1 from s)
		returning id, "query", "fragment"
	)
	select id
	from i
	union all
	select id
	from s`
const insertIntoOG = `insert into opengraph(
	created_at, url, title, type, description, determiner, locale, locales_alternate, 
	images, audios, videos, profile, article, book, blob) 
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) returning id`
const insertIntoPage = `insert into pages(
	created_at, url_id, og_id, report_id, query_fragment_id, wayback, url_string) 
	values($1,$2,$3,$4,$5,$6,$7) returning id`

//PostgreStore implements Store for a Postgres database
type PostgreStore struct {
	db *sqlx.DB
}

//NewPostgreStore constructs a PostgreStore
func NewPostgreStore(db *sqlx.DB) *PostgreStore {
	if db == nil {
		panic("nil pointer passed to NewMySQLStore")
	}

	return &PostgreStore{
		db: db,
	}
}

// Insert inserts a page into the database
func (s *PostgreStore) Insert(page *Page) (*Page, error) {
	log.Println("begining page insert transaction")
	// begin a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error begining transaction: %v", err)
	}
	// marshall the opengraph
	// _, err := json.Marshal(page.OpenGraph)
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, fmt.Errorf("error marshalling opengraph: %v", err)
	// }
	//var res sqlx.Result
	// insert into opengraph
	var ogID int64
	if err := tx.QueryRow(insertIntoOG,
		page.OpenGraph.CreatedAt, page.OpenGraph.URL,
		page.OpenGraph.Title, page.OpenGraph.Type,
		page.OpenGraph.Description, page.OpenGraph.Determiner,
		page.OpenGraph.Locale, pq.Array(page.OpenGraph.LocalesAlternate),
		page.OpenGraph.Images, page.OpenGraph.Audios,
		page.OpenGraph.Videos, page.OpenGraph.Profile,
		page.OpenGraph.Article, page.OpenGraph.Book,
		page.OpenGraph).Scan(&ogID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error inserting opengraph: %v", err)
	}
	// set the opengraph id
	page.OpenGraph.ID = ogID

	// insert into the URL
	var urlID int64
	if err := tx.QueryRow(insertIntoURLs, page.URL.Hostname(),
		page.URL.EscapedPath()).Scan(&urlID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error inserting url: %v", err)
	}

	// insert into the query_fragment
	var qfID int64
	if err := tx.QueryRow(insertIntoQueryFragment, url.QueryEscape(page.URL.RawQuery),
		url.PathEscape(page.URL.Fragment)).Scan(&qfID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error inserting query_fragment: %v", err)
	}

	//now commit the transaction so that all those inserts are atomic
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error committing insert transaction: %v", err)
	}

	return page, nil
}

// Update updates a page in the database
func (s *PostgreStore) Update(id int64, page *Page) (*Page, error) {
	return nil, nil
}
