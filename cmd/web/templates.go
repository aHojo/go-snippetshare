package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/ahojo/snippetbox/pkg/forms"
	"github.com/ahojo/snippetbox/pkg/models"
)

// Define a templateData type to act as the holding structure for the data
// Any dynamic data that we want to pass to our HTML templates.
//
type templateData struct {
	AuthenticatedUser *models.User
	CurrentYear       int
	Flash             string
	CSRFToken         string
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	// FormData		url.Values // Same underlying type as r.PostForm
	// FormErrors	map[string]string
	Form *forms.Form // This is replaceing the FormData and FormErrors field.
}

// initialize a template.FuncMap
// This creates a lookup for us
var functions template.FuncMap = template.FuncMap{
	"humanDate": humanDate,
}

// Create a human readable string of the time given by the database.
// THIS CAN ONLY RETURN 1 value.
func humanDate(t time.Time) string {
	if t.IsZero(){
		return ""
	}
	return t.Format("02 Jan 2006 at 15:04")
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	// Initialize a new map to act as the cache for our templates.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths, with the extension '.page.tmpl'
	// This will return a slice of all the page templates.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages and load them into the cache
	for _, page := range pages {
		// Extract the file name from the full file path
		name := filepath.Base(page)

		//Parse the page template file in to a template set
		// Registers our functions to the template set
		// .New() creates an empty template set
		// .Funcs() returns the template.FuncMap
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use th ParseGlob method to add any layout templates
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use th ParseGlob method to add any layout templates
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache
		cache[name] = ts
	}
	// Return the map
	return cache, nil
}
