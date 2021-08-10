package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverHelper writes an error message and stack trace to the errorLog
// then sends a 500 response to the client.
func (app *application) serverError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// clientError helper sends specific status code and description to the client.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}


// notFound handler sends a 404 response to the client.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}