package main

import (
	"fmt"
	"net/http"
)



func secureHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w,r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w,r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Create a deferred function (which will always run in the event of a panic)
		defer func(){
			// Use recover to check if there has been a panic
			if err := recover(); err != nil {
				// Set header for Connection: close
				w.Header().Set("Connection", "close")
				// Return a 500 Internal Server Error
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		
		next.ServeHTTP(w,r)
	})
}