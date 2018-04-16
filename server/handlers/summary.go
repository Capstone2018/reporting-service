package handlers

import (
	"io"
	"net/http"

	"github.com/Capstone2018/reporting-service/server/models/pages"
)

//SummaryHandler handles requests for the page summary APIandler
func (ctx *Context) SummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	pageURL := r.FormValue("url")
	if len(pageURL) == 0 {
		http.Error(w, "missing `url` parameter", http.StatusBadRequest)
		return
	}

	// fetch the html
	body, err := fetchHTML(pageURL)
	if err != nil {
		http.Error(w, "error fetching URL: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer body.Close()

	// create a new opengraph struct and process the stream
	og := pages.NewOpenGraph()
	if err := og.ProcessStream(pageURL, body); err != nil && err != io.EOF {
		http.Error(w, "error extracting summary: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, og, http.StatusOK)
}
