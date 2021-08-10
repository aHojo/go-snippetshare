package main

import (
	"fmt"
	// "html/template"
	"net/http"
	"strconv"

	"github.com/ahojo/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Check if the currenct request URL path exactly matches "/"
	// If it doesn't, give a not found error and return from this handler
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// TODO  WE WILL USE THIS LATER

	/* 	// Initialize a string slice for the paths of all tempates
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
	*/

	snippets, err := app.snippets.GetRecent()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%v\n", snippet)
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

	// Get the snippet from the database
	// If it doesn't exist, return a 404 not found error
	snippet, err := app.snippets.Get(id)

	if err == models.ErrNoRecords {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Fprintf(w, "%v", snippet)

}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		// w.WriteHeader(http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Create some dummy data to insert into the database
	var title string = "For"
	var content string = `for i <= 3 {\nfmt.Println(i)\ni = i + 1\n}`
	expires := "7"

	// Pass the data to the SnippetModel.Insert function. Get the id back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect to the page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
