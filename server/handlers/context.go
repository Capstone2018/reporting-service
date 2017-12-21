package handlers

import "github.com/Capstone2018/reporting-service/models/reports"

// Context holds global context values for the handlers
type Context struct {
	ReportsStore reports.Store
}

// NewHandlerContext returns a new handler context for globals
func NewHandlerContext(reportsStore reports.Store) *Context {
	return &Context{
		ReportsStore: reportsStore,
	}
}
