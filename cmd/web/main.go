package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ahojo/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql" // import mysql driver
)

type Config struct {
	Addr      string
	StaticDir string
}

// Application struct to inject the dependencies of our application
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

func main() {

	cfg := new(Config)
	// Set the flags for our server. flag.Parse() is needed otherwise it will always use the defaults.
	flag.StringVar(&cfg.Addr, "addr", ":4000", "Http Network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Http Network address")
	dsn := flag.String("dsn", "root:password@/snippetbox?parseTime=true", "Mysql connection info") // parseTime=true is needed to use time.Time. Converts MYSQL datetime to time.Time
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

	// Database connection

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{

		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &mysql.SnippetModel{DB: db}, // pass the database connection to our snippet model
	}
	// End Database setup

	infoLog.Printf("Starting server on %s", cfg.Addr)
	// Start our server on port 4000, pass in our mux.
	// err := http.ListenAndServe(cfg.Addr, mux)
	// errorLog.Fatal(err) // calls os.exit(1)
	srv := &http.Server{
		Addr:     cfg.Addr,
		Handler:  app.routes(cfg), // cfg is already a pointer.,
		ErrorLog: errorLog,
	}
	// ListenAndServe is blocking, so we need to start it in a goroutine
	err = srv.ListenAndServe()
	errorLog.Fatal(err) // calls os.exit(1)
}

// openDB wraps the mysql.Open function and returns a database connection
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
