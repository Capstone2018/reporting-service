package handlers

import (
	"fmt"
	"net/http"

	"github.com/Capstone2018/reporting-service/server/models/reports"
)

// ReportsHandler handles the /reports resource
func (ctx *Context) ReportsHandler(w http.ResponseWriter, r *http.Request, sessState *SessionState) {
	// TODO: add in querying ability
	if r.Method != "POST" {
		http.Error(w, "method must be POST", http.StatusBadRequest)
		return
	}

	switch r.Method {
	// get reports with query string
	case "GET":
		// // TODO: write authentication logic.. (currently you just need to be signed in to access this -- so anyone with an account)
		// // who is allowed to query the reports from the database?
		// host := strings.ToLower(r.FormValue("host"))
		// url := strings.ToLower(r.FormValue("url"))
		// if len(host) == 0 && len(url) == 0 {
		// 	// return all the reports in the database
		// 	reports, err := ctx.ReportsStore.GetAll()
		// 	if err != nil {
		// 		http.Error(w, fmt.Sprintf("error getting all reports: %v", err), http.StatusInternalServerError)
		// 		return
		// 	}
		// 	respond(w, reports, http.StatusOK)
		// } else {
		// 	if len(host) != 0 && len(url) != 0 {
		// 		http.Error(w, "can't provide both url and host query", http.StatusBadRequest)
		// 		return
		// 	}
		// 	// do host database query
		// 	if len(host) != 0 {
		// 		reports, err := ctx.ReportsStore.GetByHost(host)
		// 		if err != nil {
		// 			http.Error(w, fmt.Sprintf("error querying by host name: %v", err), http.StatusInternalServerError)
		// 			return
		// 		}
		// 		respond(w, reports, http.StatusOK)
		// 	} else { // do url database query
		// 		reports, err := ctx.ReportsStore.GetByURL(host)
		// 		if err != nil {
		// 			http.Error(w, fmt.Sprintf("error querying by url: %v", err), http.StatusInternalServerError)
		// 			return
		// 		}
		// 		respond(w, reports, http.StatusOK)
		// 	}
		// }
	// post new report to the db
	case "POST":
		// create an empty report struct and decode a json response to it
		nr := &reports.NewReport{}
		if err := receive(r, nr); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// convert the new report to a report
		report, err := nr.ToReport()
		if err != nil {
			http.Error(w, fmt.Sprintf("error converting new report: %v", err), http.StatusBadRequest)
			return
		}

		// TODO: get the Opengraph properties

		// archive.org the submitted URL
		// waybackID, err := archive(u.String())
		// if err != nil {
		// 	log.Printf("unable to archive URL: %v", err)
		// }

		// write the report to the database
		r, err := ctx.ReportsStore.Insert(report)
		if err != nil {
			http.Error(w, fmt.Sprintf("error saving new report: %v", err), http.StatusInternalServerError)
			return
		}

		// TODO: write to the mq when we decide to build a triage system

		respond(w, r, http.StatusCreated)
	}
}

// // ReportIDHandler handles the /users/<report-id> resource
// func (ctx *Context) ReportIDHandler(w http.ResponseWriter, r *http.Request, sessState *SessionState) {
// 	// TODO: ensure authentication

// 	// TODO: decide if a person can edit their report..
// 	if r.Method != "GET" {
// 		http.Error(w, "method must be GET", http.StatusBadRequest)
// 		return
// 	}
// 	// get the report ID and query the database
// 	_, idString := path.Split(r.URL.Path)
// 	// make sure that the user passed in an integer
// 	reportID, err := strconv.ParseInt(idString, 10, 64)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("id must be an integer: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	// query db and write the report back to the user
// 	report, err := ctx.ReportsStore.GetByID(reportID)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("error querying by id: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	respond(w, report, http.StatusOK)
// }
