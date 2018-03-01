package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Capstone2018/reporting-service/server/models/pages"
)

// PagesHandler handles the /pages resource
func (ctx *Context) PagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method must be GET", http.StatusBadRequest)
		return
	}
	pageURL := r.FormValue("url")
	if len(pageURL) == 0 {
		http.Error(w, "missing `url` parameter", http.StatusBadRequest)
		return
	}

	// validate the url (use ParseRequestURI first so we can make sure that it's not relative path)
	u, err := url.ParseRequestURI(pageURL)
	if err != nil {
		http.Error(w, "invalid url: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Parse the url so we can keep the query fragment
	u, err = url.Parse(pageURL)
	if err != nil {
		http.Error(w, "invalid url: "+err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: figure out if we should quit if the url fetching broke
	// fetch the page so we can parse it
	body, err := fetchHTML(u.String())
	if err != nil {
		http.Error(w, "error fetching URL: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer body.Close()

	// get the opengraph
	og := pages.NewOpenGraph()
	if err := og.ProcessStream(u.String(), body); err != nil {
		// don't quit, we still want to store the page, but just insert a null
		http.Error(w, "error fetching Opengraph: "+err.Error(), http.StatusInternalServerError)
	}

	// archive the url -- this is super slow and should be on it's own goroutine
	a := pages.NewArchive()
	if err := a.Archive(u.String()); err != nil {
		http.Error(w, "error archiving url: "+err.Error(), http.StatusInternalServerError)
	}

	// TODO: remove this, object replacement
	og = &pages.OpenGraph{
		CreatedAt:        time.Now(),
		Title:            "test",
		Description:      "blah blah blah",
		LocalesAlternate: []string{"french", "english", "latin"},
		Images:           []*pages.Image{&pages.Image{URL: "http://google.com", SecureURL: "https://google.com", Type: "uh"}},
	}

	// enter a row in the database for the page TODO: decide if we should separate this into
	// seperate interface methods????
	np := &pages.Page{
		CreatedAt: time.Now(),
		URL:       u,
		URLString: u.String(),
		OpenGraph: og,
		WaybackID: a.WaybackID,
	}

	p, err := ctx.PageStore.Insert(np)
	if err != nil {
		http.Error(w, "error inserting page: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// concurrently? write to message queue (for python workers)

	// write the page back to the user
	respond(w, p, http.StatusCreated)
}

var client = &http.Client{
	Timeout: time.Second * 30,
}

//fetchHTML does an HTTP GET for the pageURL
//and returns an error if the status code is >= 400
//or the content type doesn't start with `text/html`.
//If all goes well, it returns the response body.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	resp, err := client.Get(pageURL)
	if err != nil {
		return nil, err
	}

	//check status code
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("response status code was %d", resp.StatusCode)
	}

	//check content-type
	contentType := resp.Header.Get(headerContentType)
	if !strings.HasPrefix(contentType, contentTypeTextHTML) {
		return nil, fmt.Errorf("requested URL is not an HTML page: content type was %s", contentType)
	}

	return resp.Body, nil
}
