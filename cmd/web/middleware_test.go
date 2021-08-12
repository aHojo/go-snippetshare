package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {

	// Setup
	rr := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	framOptions := rs.Header.Get("X-Frame-Options")
	if framOptions != "deny" {
		t.Errorf("Expected X-Frame-Options to be DENY, got %s", framOptions)
	}

	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("Expected X-XSS-Protection to be 1; mode=block, got %s", xssProtection)
	}

	if rs.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be 200, got %d", rs.StatusCode)
	}

	defer rs.Body.Close()

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("Expected body to be OK, got %s", string(body))
	}
}
