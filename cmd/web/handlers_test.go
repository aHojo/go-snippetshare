package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()

	r, errr := http.NewRequest("GET", "/ping", nil)
	if errr != nil {
		t.Fatal(errr)
	}

	ping(rr, r)

	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("StatusCode: %d, expected %d", rs.StatusCode, http.StatusOK)
	}

	defer rs.Body.Close()

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "pong" {
		t.Errorf("Body: %s, expected %s", string(body), "pong")
	}

}
