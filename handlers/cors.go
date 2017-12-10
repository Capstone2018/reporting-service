package handlers

import (
	"net/http"
)

//CORS is a middleware adapter that provides Cross Origin Resource Sharing support
func CORS() Adapter {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//add all of the CORS headers
			w.Header().Add(headerAccessControlAllowOrigin, "*")
			w.Header().Add(headerAccessControlAllowMethods, "GET, PUT, POST, PATCH, DELETE, LINK, UNLINK")
			w.Header().Add(headerAccessControlAllowHeaders, "Content-Type, Authorization")
			w.Header().Add(headerAccessControlExposeHeaders, "Authorization")

			//if the request method is OPTIONS (pre-flight CORS request),
			//simply responds with http.StatusOK
			//else, call the handler's ServeHTTP() method
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
			} else {
				handler.ServeHTTP(w, r)
			}
		})
	}
}
