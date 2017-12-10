package handlers

import (
	"fmt"
	"net/http"

	"github.com/Capstone2018/reporting-service/models/reports"
)

// ReportsHandler handles the /reports resource
func (ctx *Context) ReportsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// get all reports
	case "GET":

	case "POST":
		nr := &reports.NewReport{}
		if err := receive(r, nr); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		r, err := ctx.ReportsStore.Insert(nr)
		if err != nil {
			http.Error(w, fmt.Sprintf("error saving new report: %v", err), http.StatusInternalServerError)
			return
		}

		// TODO: write to the mq when we decide to build a triage system

		respond(w, r, http.StatusCreated)
	}
}
