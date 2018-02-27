package pages

import (
	"database/sql"
	"fmt"
)

//const insertIntoHostnames = `insert into hostnames(host) values($1)`
const insertIntoURLs = `insert into urls(host_id, path) values($1, $2)`
const insertIntoQueryFragment = `insert into query_fragment(query, fragment) values($1,$2)`
const insertIntoOG = `insert into opengraph(created_at, url, title, type, description, determiner, locale, locales_alternate, images, audios, videos, profile, article, book, blob) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
const insertIntoPage = `insert into pages(created_at, url_id, og_id, report_id, query_fragment_id, wayback, url_string) values($1,$2,$3,$4,$5,$6,$7)`

//PostgreStore implements Store for a Postgres database
type PostgreStore struct {
	db *sql.DB
}

//NewPostgreStore constructs a PostgreStore
func NewPostgreStore(db *sql.DB) *PostgreStore {
	if db == nil {
		panic("nil pointer passed to NewMySQLStore")
	}

	return &PostgreStore{
		db: db,
	}
}

// Insert inserts a page into the database
func (s *PostgreStore) Insert(page *Page) (*Page, error) {
	// begin a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error begining transaction: %v", err)
	}

	var res sql.Result
	// insert into opengraph
	if res, err = tx.Exec(insertIntoOG, page.OpenGraph.CreatedAt, page.OpenGraph.URL, page.OpenGraph.Title, page.OpenGraph.Type, page.OpenGraph.Description, page.OpenGraph.Determiner, page.OpenGraph.LocalesAlternate, page.OpenGraph.Images, page.OpenGraph.Audios, page.OpenGraph.Videos, page.OpenGraph.Profile, page.OpenGraph.Article, page.OpenGraph.Book); err != nil {
		// TODO: handle error, handle json -- needs to be with sqlx probably
	}
	fmt.Println(res)
	return nil, nil
}

// Update updates a page in the database
func (s *PostgreStore) Update(id int64, page *Page) (*Page, error) {
	return nil, nil
}
