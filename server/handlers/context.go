package handlers

import (
	"github.com/Capstone2018/reporting-service/server/models/pages"
	"github.com/Capstone2018/reporting-service/server/models/reports"
	"github.com/Capstone2018/reporting-service/server/sessions"
)

// Context holds global context values for the handlers
type Context struct {
	ReportsStore      reports.Store
	PageStore         pages.Store
	SessionStore      sessions.Store
	SessionSigningKey string
}

// NewHandlerContext returns a new handler context for globals
func NewHandlerContext(reportsStore reports.Store, pageStore pages.Store) *Context {
	return &Context{
		ReportsStore: reportsStore,
		PageStore:    pageStore,
		//SessionStore:      sessionsStore,
		//SessionSigningKey: sessionsSigningKey,
	}
}
