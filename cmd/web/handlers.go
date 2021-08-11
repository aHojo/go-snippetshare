package main

import (
	"fmt"

	// "html/template"
	"net/http"
	"strconv"

	"github.com/ahojo/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	/* CODE BELOW IS HANDLED BY PAT NOW
	// Check if the currenct request URL path exactly matches "/"
	// If it doesn't, give a not found error and return from this handler
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }
	*/

	snippets, err := app.snippets.GetRecent()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := templateData{
		Snippets: snippets,
	}

	app.render(w, r, "home.page.tmpl", &data)
	/*
		OLD WAY OF DOING THINGS
		// TODO REMOVE
	*/
	// // Initialize a string slice for the paths of all tempates
	// // home MUST be first
	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }

	// // Parsefiles function reads the template file into a template set
	// ts, err := template.ParseFiles(files...)

	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// // Write the template content as the response body.
	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	/*
		END OF OLD WAY
	*/

	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%v\n", snippet)
	// }
}

// show a snippet
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	// Get the id from the query params - ?id=<somenumber>
	// .Query().Get() returns "" if it doesn't exist
	// Atoi will return an error if it can't convert the value.
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// Pat doesn't strip the : from the param
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

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

	// Create the type that will hold the data that we will pass into the template
	snippetData := templateData{
		Snippet: snippet,
	}
	app.render(w, r, "show.page.tmpl", &snippetData)
	/* OLD WAY OF DOING THINGS
	// TODO REMOVE
	// Render the template
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// execute the template
	err = ts.Execute(w, snippetData)
	if err != nil {
		app.serverError(w, err)
	}
	// fmt.Fprintf(w, "%v", snippet)
	*/
}

// Handler that will return a form to create a new snippet
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "create.page.tmpl", nil)

	
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// No longer need this because of PAT
	// if r.Method != "POST" {
	// 	w.Header().Set("Allow", "POST")
	// 	// w.WriteHeader(http.StatusMethodNotAllowed)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	// Parse the form data is the POST body
	// puts it in a map called r.PostForm
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// use the r.PostForm map to get the values of the form fields
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	
	// Pass the data to the SnippetModel.Insert function. Get the id back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect to the page for the snippet
	// http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
	// We are now using semantic url
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
