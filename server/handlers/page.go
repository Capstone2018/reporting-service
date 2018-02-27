package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Capstone2018/reporting-service/server/models/pages"
)

// PageHandler handles the /page resource
func (ctx *Context) PageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method must be GET", http.StatusBadRequest)
		return
	}
	pageURL := r.FormValue("url")
	if len(pageURL) == 0 {
		http.Error(w, "missing `url` parameter", http.StatusBadRequest)
		return
	}

	// fetch the page so we can parse it
	body, err := fetchHTML(pageURL)
	if err != nil {
		http.Error(w, "error fetching URL: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer body.Close()

	// enter a row in the database for the page
	np := &pages.Page{
		CreatedAt: time.Now(),
		URLString: pageURL,
	}
	p, err := ctx.PageStore.Insert(np)
	fmt.Println(p.ID)
	// write to message queue (for python workers)

	// get the opengraph
	og := pages.NewOpenGraph()
	if err := og.ProcessStream(pageURL, body); err != nil {

	}

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
