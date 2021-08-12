package main

import (
	"fmt"

	// "html/template"
	"net/http"
	"strconv"

	"github.com/ahojo/snippetbox/pkg/forms"
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

	// Use the PopString method to retrieve the value for the "flash" key.
	// PopString will also delete the key and value from the session data, so 
	// it acts like a one-time fetch. 
	// If there is no matching key, it will return an empty string.
	//flash := app.session.PopString(r, "flash") no longer needed because it's in addDefaultData

	// Create the type that will hold the data that we will pass into the template
	snippetData := templateData{
		Snippet: snippet,
		//Flash: flash, no longer needed because it's in addDefaultData
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

	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

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

	// Create a new form Struct that contains the POST form data
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// If the form isn't valid, redisplay the form passin in the form.Form object
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded
	// We can use the Get() method to retrieve the validated data
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
	}

	// Use the Put() method from golangcollege session package to add a string value, 
	// and the key "flash" to the session data.
	// If there is no session for the user, it will create a new empty session. 
	app.session.Put(r, "flash", "Snippet successfully created!")


	// Redirect to the page for the snippet
	// http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
	// We are now using semantic url
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

/* SIGNUP SECTION */
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request){
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {

	// Parse the form Data
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}	

	// Validate the form contents using the form helper we made earlier
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRegex)
	form.MinLength("password", 10)

	// If there are errors, redisplay the form
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	// Insert the user to the database
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		fmt.Println(err == models.ErrDuplicateEmail)
		form.Errors.Add("email", "Email already exists")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Add a confirmation flash message
	app.session.Put(r, "flash", "You have successfully signed up!")
	// redirect to the login page
	http.Redirect(w,r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "login.page.tmpl", &templateData{ Form: forms.New(nil) })

}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check if valid credentials
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Wrong email or password")
		app.render(w,r, "login.page.tmpl", &templateData{ Form: form })
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Add the id to the current session
	app.session.Put(r, "id", id)

	// redirect to the create snippet page
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the id from the current session
	app.session.Remove(r, "id")
	// add a flash message to the session to confirm they've been logged out
	app.session.Put(r, "flash", "You have successfully logged out!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}


/*
WE USED THIS BEFORE CREATING THE PARSE FORM STRUCT

	// use the r.PostForm map to get the values of the form fields
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	// Create a map to hold any of our errors
	errors := make(map[string]string)

	// Check if the title is not blank and not more than 100 chars
	if strings.TrimSpace(title) == "" {
		errors["title"] = "Title cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "Title cannot be more than 100 characters"
	}

	// Check that the Content is not blank
	if strings.TrimSpace(content) == "" {
		errors["content"] = "Content cannot be blank"
	}

	// Check that the expires is not blank and matches one of the permitted values
	// Permitted Values: "1", "7", "365"
	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "Expires cannot be blank"
	} else if expires != "1" && expires != "7" && expires != "365" {
		errors["expires"] = "Expires must be 1, 7, or 365"
	}

	// If there are any errors, send them back to the client
	// Sends the errors, and previously submited form data
	if len(errors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
	}

	// Pass the data to the SnippetModel.Insert function. Get the id back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
*/
