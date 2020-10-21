package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"text/template"
)

func (s *server) handleAPI(rand *rand.Rand) http.HandlerFunc {
	thing := prepareThing(rand) //one-time, read-only
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, thing)
	}
}

func prepareThing(rand *rand.Rand) int64 {
	return rand.Int63()
}

func (s *server) handleGreet(format string) http.HandlerFunc {
	// types just for this handler
	type request struct {
		Name string
	}
	type response struct {
		Greeting string `json:"greeting"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			respond(w, r, nil, http.StatusBadRequest)
			return
		}
		greeting := fmt.Sprintf(format, req.Name)
		resp := response{Greeting: greeting}
		respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleTemplate() http.HandlerFunc {
	// lazy init just for this handler and not for the whole app
	// this is both thread safe and speeds app startup
	var (
		init   sync.Once
		tpl    *template.Template
		tplerr error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, tplerr = template.ParseFiles("some", "files")
		})
		if tplerr != nil {
			http.Error(w, tplerr.Error(), http.StatusInternalServerError)
			return
		}
		respond(w, r, tpl, http.StatusOK)
	}
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, "here is the index!", http.StatusOK)
	}
}

// a helper utility to respond
func respond(w http.ResponseWriter, r *http.Request,
	data interface{}, status int) {
    w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// a helper for decoding stuff
func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// later check content type here. initially, just JSON
	return json.NewDecoder(r.Body).Decode(v)
}
