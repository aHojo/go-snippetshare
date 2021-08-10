package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Check if the currenct request URL path exactly matches "/"
	// If it doesn't, give a not found error and return from this handler
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Initialize a string slice for the paths of all tempates
	// home MUST be first
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Parsefiles function reads the template file into a template set
	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the template content as the response body.
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

// show a snippet
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	// Get the id from the query params - ?id=<somenumber>
	// .Query().Get() returns "" if it doesn't exist
	// Atoi will return an error if it can't convert the value.

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "We are in show a snippet, displaying snippet: %d", id)

}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		// w.WriteHeader(http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Create a snippet")
}
