package handlers

import (
	"net/http"
)

type requestContextKey string

//Adapter represents a function that wraps an http.Handler with
//some sort of middleware functionality
type Adapter func(http.Handler) http.Handler

//Adapt wraps the http.Handler with the provided Adapters
func Adapt(handler http.Handler, adapters ...Adapter) http.Handler {
	for i := len(adapters) - 1; i >= 0; i-- {
		handler = adapters[i](handler)
	}
	return handler
}
