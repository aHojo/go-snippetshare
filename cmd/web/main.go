package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr string
	StaticDir string
}

// Application struct to inject the dependencies of our application
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	cfg := new(Config)
	// Set the flags for our server. flag.Parse() is needed otherwise it will always use the defaults.
	flag.StringVar(&cfg.Addr, "addr", ":4000", "Http Network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Http Network address")
	flag.Parse()

	/* // if we want to log to a file, we can use the standard log package
	f, err := os.OpenFile("web.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime) */
	
	// Create a logger for writing informational messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
	
		infoLog: infoLog,
		errorLog: errorLog,
	}

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

	infoLog.Printf("Starting server on %s", cfg.Addr)
	// Start our server on port 4000, pass in our mux.
	// err := http.ListenAndServe(cfg.Addr, mux)
	// errorLog.Fatal(err) // calls os.exit(1)
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
		ErrorLog: errorLog,
	}
	// ListenAndServe is blocking, so we need to start it in a goroutine
	err := srv.ListenAndServe()
	errorLog.Fatal(err) // calls os.exit(1)
}
