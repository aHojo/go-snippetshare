package main

import "net/http"

func (app *application) routes(cfg *Config) *http.ServeMux{
	// Use the http.NewServeMux() to initialize a new servemux, then
	// register the home function as the handler for the "/" path
	mux := http.NewServeMux() // this is  the default, but still define it for security.

	mux.HandleFunc("/", app.home)                        // subtree path, has an ending /
	mux.HandleFunc("/snippet", app.showSnippet)          // fixed path, url must match this exactly.
	mux.HandleFunc("/snippet/create", app.createSnippet) // fixed path, url must match this exactly.

	// Create a fileserver to serve static content from
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))

	// use the mux.Handle() to register the file serveras the handler
	// all url paths start with /static/. Strip the /static prefix before
	// the request reaches the file server
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
