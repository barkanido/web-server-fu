package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestGreet(t *testing.T) {
	is := is.New(t)
	db := make(map[string]string)
	srv := newServer(&db)
	r := httptest.NewRequest("POST", "/greet", strings.NewReader("{\"name\":\"ido\"}"))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)
	w.Result().Body.Close()
	is.Equal(w.Result().StatusCode, http.StatusOK)
	is.Equal(w.Body.String(), "{\"greeting\":\"hello ido\"}\n")
}
