package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ahojo/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})

}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: true,
	})
	
	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check if a userID value is present in the session
		// If it is not present, call the next handler in the chain
		exists := app.session.Exists(r, "id")
		if !exists {
			next.ServeHTTP(w,r)
			return
		}

		// Fetch the details of the current user from the DB.
		// If no matching records are found, remove the invalid userID from
		// their session and call the next handler 
		user, err := app.users.Get(app.session.GetInt(r, "id"))
		if err == models.ErrNoRecords {
			app.session.Remove(r,"id")
			next.ServeHTTP(w,r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		// If we get here, the request is valid, 
		// Authenticated user
		// Create a new copy of the request wit hthe user information added
		// call the next handler in the chain
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})

}
func (app *application) recoverPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Create a deferred function (which will always run in the event of a panic)
		defer func() {
			// Use recover to check if there has been a panic
			if err := recover(); err != nil {
				// Set header for Connection: close
				w.Header().Set("Connection", "close")
				// Return a 500 Internal Server Error
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
