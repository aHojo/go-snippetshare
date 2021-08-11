package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// Use to return a *http.ServeMux but we are return a http.Handler because of middleware.
func (app *application) routes(cfg *Config) http.Handler {

	// Create a middleware chain containing our "standard middleware" that is used for every request.
	chain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Use the http.NewServeMux() to initialize a new servemux, then
	// register the home function as the handler for the "/" path
	//mux := http.NewServeMux()                            // this is  the default, but still define it for security.
	// Starting to use the GIN framework
	//mux.HandleFunc("/", app.home)                        // subtree path, has an ending /
	//mux.HandleFunc("/snippet", app.showSnippet)          // fixed path, url must match this exactly.
	//mux.HandleFunc("/snippet/create", app.createSnippet) // fixed path, url must match this exactly.
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	// Create a fileserver to serve static content from
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))

	// use the mux.Handle() to register the file serveras the handler
	// all url paths start with /static/. Strip the /static prefix before
	// the request reaches the file server
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	// without middleware
	// return mux

	// If we do not use alice
	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	// With Alice
	return chain.Then(mux)
}
