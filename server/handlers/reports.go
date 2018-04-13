package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Capstone2018/reporting-service/server/models/pages"
	"github.com/Capstone2018/reporting-service/server/models/reports"
)

// ReportsHandler handles the /reports resource
func (ctx *Context) ReportsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: add in querying ability? or is this on its own service
	if r.Method != "POST" {
		http.Error(w, "method must be POST", http.StatusBadRequest)
		return
	}

	switch r.Method {
	// post new report to the db
	case "POST":
		// create an empty report struct and decode a json response to it
		nr := &reports.NewReport{}
		if err := receive(r, nr); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// convert the new report to a report
		// this validates the urls sent and discards invalid urls - not throwing error in that case
		report, err := nr.ToReport()
		if err != nil {
			http.Error(w, fmt.Sprintf("error converting new report: %v", err), http.StatusBadRequest)
			return
		}

		// insert a new report, get the returned report id
		report, err = ctx.ReportsStore.Insert(report)
		if err != nil {
			http.Error(w, fmt.Sprintf("error inserting new report: %v", err), http.StatusInternalServerError)
			return
		}

		// handle each url -- opengraph, archive and page insertion
		for _, url := range report.URLs {
			//TODO: decide level of error reporting??? do we want to report to the user that their url wasn't saved properly?
			go ctx.handlePage(url, report.ID)
		}

		respond(w, report, http.StatusCreated)
	}
}

// handlePage parses the metadata of a url's "page" and inserts it into the database
func (ctx *Context) handlePage(u *url.URL, reportID int64) {
	// archive the url
	a := pages.NewArchive()
	if err := a.Archive(u.String()); err != nil {
		log.Printf("Error archiving url: %v", u.String())
	}
	// TODO: figure out if we should quit if the url fetching broke
	// fetch the page so we can parse it
	body, err := fetchHTML(u.String())
	if err != nil {
		log.Printf("error fetching url: %v", u.String())
		return
	}
	defer body.Close()

	// fetch the opengraph
	// get the opengraph
	og := pages.NewOpenGraph()
	if err := og.ProcessStream(u.String(), body); err != nil {
		// don't quit, we still want to store the page, but just insert a null opengraph
		log.Printf("error processing opengraph: %v, err: %v", u.String(), err.Error())
	}

	// insert the page into the database
	np := &pages.Page{
		CreatedAt: time.Now(),
		URL:       u,
		URLString: u.String(),
		OpenGraph: og,
		WaybackID: a.WaybackID,
		ReportID:  reportID,
		Query:     url.QueryEscape(u.RawQuery),
		Fragment:  url.PathEscape(u.Fragment),
	}

	_, err = ctx.PageStore.Insert(np)
	if err != nil {
		log.Printf("error storing new page from: %v, %v", np.URLString, err)
	}
}
