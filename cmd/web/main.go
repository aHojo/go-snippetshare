package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ahojo/snippetbox/pkg/models/mysql"
	"github.com/golangcollege/sessions"

	_ "github.com/go-sql-driver/mysql" // import mysql driver
)


// Defining our new type for context
type contextKey string 
var contextKeyUser = contextKey("user")

type Config struct {
	Addr      string
	StaticDir string
}

// Application struct to inject the dependencies of our application
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	session       *sessions.Session
	users         *mysql.UserModel
}

func main() {

	cfg := new(Config)
	// Set the flags for our server. flag.Parse() is needed otherwise it will always use the defaults.
	flag.StringVar(&cfg.Addr, "addr", ":4000", "Http Network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Http Network address")
	dsn := flag.String("dsn", "root:password@/snippetbox?parseTime=true", "Mysql connection info") // parseTime=true is needed to use time.Time. Converts MYSQL datetime to time.Time
	// Needs to be 32 bytes long. Used to encrypt and authenticate session cookies.
	secret := flag.String("secret", "z6MAh3pPbnEHbf*+3Gd8qGWKTzbpa@ge", "Secret")
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

	// Create the template Cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	// Database connection

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Session Management
	// Use the sessions.New() function to initialize a new session manager.
	// Pass in the secret key as the param, sessions will expire after 12 hours.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := &application{

		infoLog:       infoLog,
		errorLog:      errorLog,
		snippets:      &mysql.SnippetModel{DB: db}, // pass the database connection to our snippet model
		users:         &mysql.UserModel{DB: db},
		templateCache: templateCache,
		session:       session,
	}
	// End Database setup

	// Create this tls.Config struct to hold non default settings
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true, //  field controls whether the HTTPS connection should use Go???s favored cipher suites or the user???s favored cipher suites.
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		/*
			If we want to use a custom cipher suite, we can do so by adding it to the CipherSuites field.
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},

			Recommended here: https://wiki.mozilla.org/Security/Server_Side_TLS

			Set accepted TLS versions.
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS12,

		*/
	}

	// Start our server on port 4000, pass in our mux.
	// err := http.ListenAndServe(cfg.Addr, mux)
	// errorLog.Fatal(err) // calls os.exit(1)
	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      app.routes(cfg), // cfg is already a pointer.,
		ErrorLog:     errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s", cfg.Addr)
	// ListenAndServe is blocking, so we need to start it in a goroutine
	// err = srv.ListenAndServe()
	// For TLS
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
